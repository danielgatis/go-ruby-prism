package parser

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/danielgatis/go-ruby-prism/wasm"
)

// Parser parses Ruby source into a syntax tree.
//
// A Parser owns a pool of WebAssembly instances. Each instance has its own
// linear memory and can only run one parse at a time, so the pool size sets
// how many parses can run concurrently. A Parser is safe for use by multiple
// goroutines.
type Parser struct {
	// pool hands out instances. It is buffered to poolSize and is filled
	// lazily: an empty pool with room left to grow creates a new instance
	// rather than waiting for one to be returned.
	pool     chan *wasm.Runtime
	poolSize int

	// mutex guards created and closed.
	mutex   sync.Mutex
	created int
	closed  bool

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

// NewParser creates a parser. By default the pool is sized to GOMAXPROCS; use
// WithPoolSize to override it.
//
// One instance is created eagerly so that configuration errors surface here,
// and the rest are created on demand.
func NewParser(ctx context.Context, options ...ParserOption) (*Parser, error) {
	parser := &Parser{
		poolSize: runtime.GOMAXPROCS(0),
		logger:   NewNullLogger(),
	}

	for _, opt := range options {
		opt(parser)
	}

	if parser.poolSize < 1 {
		return nil, fmt.Errorf("pool size must be at least 1, got %d", parser.poolSize)
	}

	parser.pool = make(chan *wasm.Runtime, parser.poolSize)

	instance, err := wasm.NewRuntime(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasm runtime: %w", err)
	}

	parser.created = 1
	parser.pool <- instance

	return parser, nil
}

// Close releases every instance in the pool. It is safe to call more than
// once. Parse returns an error after Close returns.
func (p *Parser) Close(ctx context.Context) error {
	p.mutex.Lock()

	if p.closed {
		p.mutex.Unlock()
		return nil
	}

	p.closed = true
	pending := p.created
	p.mutex.Unlock()

	// Every instance is either idle in the pool or checked out by a Parse
	// that will return it, so draining exactly the number created waits out
	// the in-flight parses without holding the lock.
	var firstErr error

	for i := 0; i < pending; i++ {
		instance := <-p.pool

		if err := instance.Close(ctx); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close the wasm runtime: %w", err)
		}
	}

	close(p.pool)

	return firstErr
}

// acquire checks an instance out of the pool, creating one if the pool is
// empty and has not reached its size limit.
func (p *Parser) acquire(ctx context.Context) (*wasm.Runtime, error) {
	// Prefer an idle instance over paying for a new one.
	select {
	case instance, ok := <-p.pool:
		if !ok {
			return nil, fmt.Errorf("parser is closed")
		}

		return instance, nil
	default:
	}

	p.mutex.Lock()

	if p.closed {
		p.mutex.Unlock()
		return nil, fmt.Errorf("parser is closed")
	}

	if p.created < p.poolSize {
		p.created++
		p.mutex.Unlock()

		instance, err := wasm.NewRuntime(ctx)
		if err != nil {
			// Give the slot back so a later parse can retry.
			p.mutex.Lock()
			p.created--
			p.mutex.Unlock()

			return nil, fmt.Errorf("failed to instantiate wasm runtime: %w", err)
		}

		return instance, nil
	}

	p.mutex.Unlock()

	// The pool is at its limit, so wait for an instance to come back.
	select {
	case instance, ok := <-p.pool:
		if !ok {
			return nil, fmt.Errorf("parser is closed")
		}

		return instance, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// release returns an instance to the pool.
func (p *Parser) release(instance *wasm.Runtime) {
	// Close drains exactly the number of instances created, so a checked-out
	// instance must be returned even once the parser is closing.
	defer func() {
		// The pool is closed only after every instance has been drained,
		// so this can race only if release is called twice.
		_ = recover()
	}()

	p.pool <- instance
}

// Parse parses source and returns the resulting tree. It is safe to call from
// multiple goroutines; each call runs on its own WebAssembly instance.
func (p *Parser) Parse(ctx context.Context, source []byte) (*ParseResult, error) {
	instance, err := p.acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer p.release(instance)

	return p.parseWith(ctx, instance, source)
}

// parseWith runs a parse on an already acquired instance.
func (p *Parser) parseWith(ctx context.Context, instance *wasm.Runtime, source []byte) (result *ParseResult, err error) {
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
	sourcePtr, err := instance.Calloc(ctx, 1, uint64(len(source))+1)
	p.logger.Debug("sourcePtr: %v", sourcePtr)

	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for source: %w", err)
	}

	defer func() {
		if freeErr := instance.Free(ctx, sourcePtr); freeErr != nil && err == nil {
			result = nil
			err = fmt.Errorf("failed to free memory for source ptr: %w", freeErr)
		}
	}()

	if !instance.MemoryWrite(sourcePtr, source) {
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

	optPtr, err := instance.Calloc(ctx, 1, uint64(len(serializedOptions))+1)
	p.logger.Debug("optPtr: %v", optPtr)

	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for options: %w", err)
	}

	defer func() {
		if freeErr := instance.Free(ctx, optPtr); freeErr != nil && err == nil {
			result = nil
			err = fmt.Errorf("failed to free memory for option ptr: %w", freeErr)
		}
	}()

	if !instance.MemoryWrite(optPtr, serializedOptions) {
		return nil, fmt.Errorf("failed to write the options into memory")
	}

	// call the serialize parse function
	bufferSizeOf, err := instance.BufferSizeOf(ctx)
	p.logger.Debug("bufferSizeOf: %v", bufferSizeOf)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer size: %w", err)
	}

	bufferPtr, err := instance.Calloc(ctx, bufferSizeOf, 1)
	p.logger.Debug("bufferPtr: %v", bufferPtr)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer ptr: %w", err)
	}

	if err := instance.BufferInit(ctx, bufferPtr); err != nil {
		// pm_buffer_init failed, so there is no internal value to clean up.
		// Release the struct allocation directly.
		if freeErr := instance.Free(ctx, bufferPtr); freeErr != nil {
			p.logger.Debug("failed to free buffer ptr after failed init: %v", freeErr)
		}

		return nil, fmt.Errorf("failed to init the buffer: %w", err)
	}

	// pm_buffer_free releases both the internal value and the pm_buffer_t
	// struct itself, so it must not be paired with an extra free(bufferPtr).
	defer func() {
		if freeErr := instance.BufferFree(ctx, bufferPtr); freeErr != nil && err == nil {
			result = nil
			err = fmt.Errorf("failed to free memory for buffer ptr: %w", freeErr)
		}
	}()

	if _, err := instance.SerializeParse(ctx, bufferPtr, sourcePtr, uint64(len(source)), optPtr); err != nil {
		return nil, fmt.Errorf("failed to call the parse function: %w", err)
	}

	// read result from memory
	bufferValue, err := instance.BufferValue(ctx, bufferPtr)
	p.logger.Debug("bufferValue: %v", bufferValue)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer value: %w", err)
	}

	bufferLen, err := instance.BufferLength(ctx, bufferPtr)
	p.logger.Debug("bufferLen: %v", bufferLen)

	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer length: %w", err)
	}

	serializedBytes, ok := instance.MemoryRead(bufferValue, bufferLen)
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

// WithPoolSize sets how many WebAssembly instances the parser may hold, which
// is the number of parses it can run at once. It defaults to GOMAXPROCS.
//
// Each instance carries its own linear memory, so a larger pool trades memory
// for concurrency. A size of 1 serializes every parse onto a single instance.
func WithPoolSize(size int) ParserOption {
	return func(p *Parser) {
		p.poolSize = size
	}
}
