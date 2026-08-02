package parser

import (
	"context"
	"fmt"
	"sync"

	"github.com/danielgatis/go-ruby-prism/wasm"
)

type Parser struct {
	mutex               sync.Mutex
	runtime             *wasm.Runtime
	filepath            []byte
	line                int
	encoding            []byte
	frozenStringLiteral bool
	commandLine         []CommandLine
	version             SyntaxVersion
	encodingLocked      bool
	mainScript          bool
	partialScript       bool
	scopes              [][][]byte
	logger              Logger
}

func NewParser(ctx context.Context, options ...ParserOption) (*Parser, error) {
	runtime, err := wasm.NewRuntime(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasm runtime: %w", err)
	}

	parser := &Parser{
		runtime: runtime,
		logger:  NewNullLogger(),
	}

	for _, opt := range options {
		opt(parser)
	}

	return parser, nil
}

func (p *Parser) Close(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if err := p.runtime.Close(ctx); err != nil {
		return fmt.Errorf("failed to close the wasm runtime: %w", err)
	}

	return nil
}

func (p *Parser) Parse(ctx context.Context, source []byte) (result *ParseResult, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	result = nil
	err = nil

	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	// Always allocate at least one byte. calloc(1, 0) may return a zero-sized
	// allocation that Prism still dereferences, and it makes the pointer
	// indistinguishable from a failed allocation.
	sourcePtr, err := p.runtime.Calloc(ctx, 1, uint64(len(source))+1)
	p.logger.Debug("sourcePtr: %v", sourcePtr)

	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for source: %w", err)
	}

	defer func() {
		if freeErr := p.runtime.Free(ctx, sourcePtr); freeErr != nil && err == nil {
			result = nil
			err = fmt.Errorf("failed to free memory for source ptr: %w", freeErr)
		}
	}()

	if !p.runtime.MemoryWrite(sourcePtr, source) {
		return nil, fmt.Errorf("failed to write the source into memory")
	}

	p.logger.Debug("source: %v", source)
	p.logger.Debug("filepath: %v", p.filepath)
	p.logger.Debug("line: %v", p.line)
	p.logger.Debug("encoding: %v", p.encoding)
	p.logger.Debug("frozenStringLiteral: %v", p.frozenStringLiteral)
	p.logger.Debug("commandLine: %v", p.commandLine)
	p.logger.Debug("version: %v", p.version)
	p.logger.Debug("encodingLocked: %v", p.encodingLocked)
	p.logger.Debug("mainScript: %v", p.mainScript)
	p.logger.Debug("partialScript: %v", p.partialScript)
	p.logger.Debug("scopes: %v", p.scopes)

	serializedOptions, err := serializeParserOptions(
		[]byte(p.filepath),
		p.line,
		[]byte(p.encoding),
		p.frozenStringLiteral,
		p.commandLine,
		p.version,
		p.encodingLocked,
		p.mainScript,
		p.partialScript,
		p.scopes,
	)

	p.logger.Debug("serializedOptions: %v", serializedOptions)

	if err != nil {
		return nil, fmt.Errorf("failed to serialize the parser options: %w", err)
	}

	optPtr, err := p.runtime.Calloc(ctx, 1, uint64(len(serializedOptions))+1)
	p.logger.Debug("optPtr: %v", optPtr)

	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for options: %w", err)
	}

	defer func() {
		if freeErr := p.runtime.Free(ctx, optPtr); freeErr != nil && err == nil {
			result = nil
			err = fmt.Errorf("failed to free memory for option ptr: %w", freeErr)
		}
	}()

	if !p.runtime.MemoryWrite(optPtr, serializedOptions) {
		return nil, fmt.Errorf("failed to write the options into memory")
	}

	// call the serialize parse function
	bufferSizeOf, err := p.runtime.BufferSizeOf(ctx)
	p.logger.Debug("bufferSizeOf: %v", bufferSizeOf)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer size: %w", err)
	}

	bufferPtr, err := p.runtime.Calloc(ctx, bufferSizeOf, 1)
	p.logger.Debug("bufferPtr: %v", bufferPtr)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer ptr: %w", err)
	}

	if err := p.runtime.BufferInit(ctx, bufferPtr); err != nil {
		// pm_buffer_init failed, so there is no internal value to clean up.
		// Release the struct allocation directly.
		if freeErr := p.runtime.Free(ctx, bufferPtr); freeErr != nil {
			p.logger.Debug("failed to free buffer ptr after failed init: %v", freeErr)
		}

		return nil, fmt.Errorf("failed to init the buffer: %w", err)
	}

	// pm_buffer_free releases both the internal value and the pm_buffer_t
	// struct itself, so it must not be paired with an extra free(bufferPtr).
	defer func() {
		if freeErr := p.runtime.BufferFree(ctx, bufferPtr); freeErr != nil && err == nil {
			result = nil
			err = fmt.Errorf("failed to free memory for buffer ptr: %w", freeErr)
		}
	}()

	if _, err := p.runtime.SerializeParse(ctx, bufferPtr, sourcePtr, uint64(len(source)), optPtr); err != nil {
		return nil, fmt.Errorf("failed to call the parse function: %w", err)
	}

	// read result from memory
	bufferValue, err := p.runtime.BufferValue(ctx, bufferPtr)
	p.logger.Debug("bufferValue: %v", bufferValue)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer value: %w", err)
	}

	bufferLen, err := p.runtime.BufferLength(ctx, bufferPtr)
	p.logger.Debug("bufferLen: %v", bufferLen)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer length: %w", err)
	}

	serializedBytes, ok := p.runtime.MemoryRead(bufferValue, bufferLen)
	p.logger.Debug("serializedBytes: %v", serializedBytes)

	if !ok {
		return nil, fmt.Errorf("failed to read the buffer content from memory")
	}

	result, err = Deserialize(source, serializedBytes)
	p.logger.Debug("result: %v", result)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialize the result: %w", err)
	}

	return result, nil
}

type ParserOption func(*Parser)

func WithFilePath(filepath string) ParserOption {
	return func(p *Parser) {
		p.filepath = []byte(filepath)
	}
}

func WithLine(line int) ParserOption {
	return func(p *Parser) {
		p.line = line
	}
}

func WithEncoding(encoding string) ParserOption {
	return func(p *Parser) {
		p.encoding = []byte(encoding)
	}
}

func WithFrozenStringLiteral(frozenStringLiteral bool) ParserOption {
	return func(p *Parser) {
		p.frozenStringLiteral = frozenStringLiteral
	}
}

func WithCommandLine(commandLine []CommandLine) ParserOption {
	return func(p *Parser) {
		p.commandLine = commandLine
	}
}

func WithVersion(version SyntaxVersion) ParserOption {
	return func(p *Parser) {
		p.version = version
	}
}

func WithEncodingLocked(encodingLocked bool) ParserOption {
	return func(p *Parser) {
		p.encodingLocked = encodingLocked
	}
}

func WithMainScript(mainScript bool) ParserOption {
	return func(p *Parser) {
		p.mainScript = mainScript
	}
}

func WithPartialScript(partialScript bool) ParserOption {
	return func(p *Parser) {
		p.partialScript = partialScript
	}
}

func WithScopes(scopes [][][]byte) ParserOption {
	return func(p *Parser) {
		p.scopes = scopes
	}
}

func WithLogger(logger Logger) ParserOption {
	return func(p *Parser) {
		p.logger = logger
	}
}
