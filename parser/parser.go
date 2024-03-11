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
}

func NewParser(ctx context.Context) (*Parser, error) {
	runtime, err := wasm.NewRuntime(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasm runtime: %w", err)
	}

	return &Parser{
		runtime: runtime,
	}, nil
}

func (p *Parser) Close(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if err := p.runtime.Close(ctx); err != nil {
		return fmt.Errorf("failed to close the wasm runtime: %w", err)
	}

	return nil
}

func (p *Parser) Parse(ctx context.Context, source []byte, options ...ParserOption) (result *ParseResult, err error) {
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

	sourcePtr, err := p.runtime.Calloc(ctx, 1, uint64(len(source)))
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for source: %w", err)
	}

	if !p.runtime.MemoryWrite(sourcePtr, source) {
		return nil, fmt.Errorf("failed to write the source into memory: %w", err)
	}

	// put option into memory
	for _, opt := range options {
		opt(p)
	}

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

	if err != nil {
		return nil, fmt.Errorf("failed to serialize the parser options: %w", err)
	}

	optPtr, err := p.runtime.Calloc(ctx, 1, uint64(len(serializedOptions)))
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for options: %w", err)
	}

	if !p.runtime.MemoryWrite(optPtr, serializedOptions) {
		return nil, fmt.Errorf("failed to write the options into memory: %w", err)
	}

	// call the serialize parse function
	bufferSizeOf, err := p.runtime.BufferSizeOf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer size: %w", err)
	}

	bufferPtr, err := p.runtime.Calloc(ctx, bufferSizeOf, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer ptr: %w", err)
	}

	if err := p.runtime.BufferInit(ctx, bufferPtr); err != nil {
		return nil, fmt.Errorf("failed to init the buffer: %w", err)
	}

	if _, err := p.runtime.SerializeParse(ctx, bufferPtr, sourcePtr, uint64(len(source)), optPtr); err != nil {
		return nil, fmt.Errorf("failed to call the parse function: %w", err)
	}

	// read result from memory
	bufferValue, err := p.runtime.BufferValue(ctx, bufferPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer value: %w", err)
	}

	bufferLen, err := p.runtime.BufferLength(ctx, bufferPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer length: %w", err)
	}

	serializedBytes, ok := p.runtime.MemoryRead(bufferValue, bufferLen)
	if !ok {
		return nil, fmt.Errorf("failed to read the buffer content from memory: %w", err)
	}

	result, err = load(serializedBytes, source)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize the result: %w", err)
	}

	// free memory
	if err := p.runtime.BufferFree(ctx, bufferPtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for buffer ptr: %w", err)
	}

	if err := p.runtime.Free(ctx, sourcePtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for source ptr: %w", err)
	}

	if err := p.runtime.Free(ctx, bufferPtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for buffer ptr: %w", err)
	}

	if err := p.runtime.Free(ctx, optPtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for option ptr: %w", err)
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
