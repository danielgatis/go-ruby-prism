package parser

import (
	"context"

	"github.com/danielgatis/go-ruby-prism/wasm"
	"github.com/rotisserie/eris"
)

type Parser struct {
	runtime *wasm.Runtime
}

func NewParser(ctx context.Context) (*Parser, error) {
	runtime, err := wasm.NewRuntime(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "failed to instantiate wasm runtime")
	}

	return &Parser{
		runtime: runtime,
	}, nil
}

func (p *Parser) Close(ctx context.Context) error {
	if err := p.runtime.Close(ctx); err != nil {
		return eris.Wrap(err, "failed to close the wasm runtime")
	}

	return nil
}

func (p *Parser) Parse(ctx context.Context, source string) (result *ParseResult, err error) {
	result = nil
	err = nil

	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = eris.Errorf("recovered from panic: %v", r)
		}
	}()

	result, err = p.parseWithOptions(ctx, source, newParseOptions())

	if err != nil {
		return nil, eris.Wrap(err, "failed to parse with options")
	}

	return result, nil
}

func (p *Parser) parseWithOptions(ctx context.Context, source string, opts *parseOptions) (*ParseResult, error) {
	// put source into memory
	sourceBytes := []byte(source)

	sourcePtr, err := p.runtime.Calloc(ctx, 1, uint64(len(sourceBytes)))
	if err != nil {
		return nil, eris.Wrap(err, "failed to allocate memory for source")
	}

	if !p.runtime.MemoryWrite(sourcePtr, sourceBytes) {
		return nil, eris.Wrap(err, "failed to write the source into memory")
	}

	// put option into memory
	optBytes, err := opts.bytes()
	if err != nil {
		return nil, eris.Wrap(err, "failed to convert options into bytes")
	}

	optPtr, err := p.runtime.Calloc(ctx, 1, uint64(len(optBytes)))
	if err != nil {
		return nil, eris.Wrap(err, "failed to allocate memory for options")
	}

	if !p.runtime.MemoryWrite(optPtr, optBytes) {
		return nil, eris.Wrap(err, "failed to write the options into memory")
	}

	// call the serialize parse function
	bufferSizeOf, err := p.runtime.BufferSizeOf(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "failed to get the buffer size")
	}

	bufferPtr, err := p.runtime.Calloc(ctx, bufferSizeOf, 1)
	if err != nil {
		return nil, eris.Wrap(err, "failed to get the buffer ptr")
	}

	if err := p.runtime.BufferInit(ctx, bufferPtr); err != nil {
		return nil, eris.Wrap(err, "failed to init the buffer")
	}

	if _, err := p.runtime.SerializeParse(ctx, bufferPtr, sourcePtr, uint64(len(sourceBytes)), optPtr); err != nil {
		return nil, eris.Wrap(err, "failed to call the parse function")
	}

	// read result from memory
	bufferValue, err := p.runtime.BufferValue(ctx, bufferPtr)
	if err != nil {
		return nil, eris.Wrap(err, "failed to get the buffer value")
	}

	bufferLen, err := p.runtime.BufferLength(ctx, bufferPtr)
	if err != nil {
		return nil, eris.Wrap(err, "failed to get the buffer length")
	}

	serializedBytes, ok := p.runtime.MemoryRead(bufferValue, bufferLen)
	if !ok {
		return nil, eris.Wrap(err, "failed to read the buffer content from memory")
	}

	// free memory
	if err := p.runtime.BufferFree(ctx, bufferPtr); err != nil {
		return nil, eris.Wrap(err, "failed to free memory for buffer ptr")
	}

	if err := p.runtime.Free(ctx, sourcePtr); err != nil {
		return nil, eris.Wrap(err, "failed to free memory for source ptr")
	}

	if err := p.runtime.Free(ctx, bufferPtr); err != nil {
		return nil, eris.Wrap(err, "failed to free memory for buffer ptr")
	}

	if err := p.runtime.Free(ctx, optPtr); err != nil {
		return nil, eris.Wrap(err, "failed to free memory for option ptr")
	}

	result, err := deserialize(serializedBytes, sourceBytes)
	if err != nil {
		return nil, eris.Wrap(err, "failed to deserialize the result")
	}

	return result, nil
}

func ErrToStr(err error) string {
	return eris.ToString(err, true)
}

func ErrToJSON(err error) map[string]interface{} {
	return eris.ToJSON(err, true)
}
