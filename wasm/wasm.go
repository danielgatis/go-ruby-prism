package wasm

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed prism.wasm
var prismWasm []byte

type Memory struct {
	mem api.Memory
}

func NewMemory(mem api.Memory) *Memory {
	return &Memory{mem: mem}
}

func (m *Memory) Write(ptr uint64, data []byte) bool {
	return m.mem.Write(uint32(ptr), data)
}

func (m *Memory) Read(offset uint64, byteCount uint64) ([]byte, bool) {
	return m.mem.Read(uint32(offset), uint32(byteCount))
}

type ModFunc struct {
	fn api.Function
}

func NewModFunc(mod api.Module, name string) *ModFunc {
	return &ModFunc{fn: mod.ExportedFunction(name)}
}

func (f *ModFunc) Call(ctx context.Context, params ...uint64) (uint64, error) {
	result, err := f.fn.Call(ctx, params...)
	if err != nil {
		return 0, fmt.Errorf("failed to call the wasm func: %w", err)
	}

	if len(result) > 0 {
		return result[0], nil
	} else {
		return 0, nil
	}
}

type Runtime struct {
	runtime             wazero.Runtime
	mem                 *Memory
	modCalloc           *ModFunc
	modFree             *ModFunc
	modPmSerializeParse *ModFunc
	modPmBufferInit     *ModFunc
	modPmBufferSizeof   *ModFunc
	modPmBufferValue    *ModFunc
	modPmBufferLength   *ModFunc
	modPmBufferFree     *ModFunc
}

func NewRuntime(ctx context.Context) (*Runtime, error) {
	runtime := wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)
	mod, err := runtime.Instantiate(ctx, prismWasm)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate prism: %w", err)
	}

	return &Runtime{
		runtime:             runtime,
		mem:                 NewMemory(mod.Memory()),
		modCalloc:           NewModFunc(mod, "calloc"),
		modFree:             NewModFunc(mod, "free"),
		modPmSerializeParse: NewModFunc(mod, "pm_serialize_parse"),
		modPmBufferInit:     NewModFunc(mod, "pm_buffer_init"),
		modPmBufferSizeof:   NewModFunc(mod, "pm_buffer_sizeof"),
		modPmBufferValue:    NewModFunc(mod, "pm_buffer_value"),
		modPmBufferLength:   NewModFunc(mod, "pm_buffer_length"),
		modPmBufferFree:     NewModFunc(mod, "pm_buffer_free"),
	}, nil
}

func (r *Runtime) Close(ctx context.Context) error {
	if err := r.runtime.Close(ctx); err != nil {
		return fmt.Errorf("failed to close the runtime: %w", err)
	}

	return nil
}

func (r *Runtime) Calloc(ctx context.Context, size uint64, count uint64) (uint64, error) {
	return r.modCalloc.Call(ctx, size, count)
}

func (r *Runtime) Free(ctx context.Context, ptr uint64) error {
	_, err := r.modFree.Call(ctx, ptr)
	if err != nil {
		return fmt.Errorf("failed to free the memory: %w", err)
	}
	return nil
}

func (r *Runtime) BufferSizeOf(ctx context.Context) (uint64, error) {
	return r.modPmBufferSizeof.Call(ctx)
}

func (r *Runtime) BufferInit(ctx context.Context, bufferPtr uint64) error {
	_, err := r.modPmBufferInit.Call(ctx, bufferPtr)
	if err != nil {
		return fmt.Errorf("failed to init the buffer: %w", err)
	}

	return nil
}

func (r *Runtime) BufferValue(ctx context.Context, bufferPtr uint64) (uint64, error) {
	return r.modPmBufferValue.Call(ctx, bufferPtr)
}

func (r *Runtime) BufferLength(ctx context.Context, bufferPtr uint64) (uint64, error) {
	return r.modPmBufferLength.Call(ctx, bufferPtr)
}

func (r *Runtime) BufferFree(ctx context.Context, bufferPtr uint64) error {
	_, err := r.modPmBufferFree.Call(ctx, bufferPtr)
	if err != nil {
		return fmt.Errorf("failed to free the buffer: %w", err)
	}

	return nil
}

func (r *Runtime) SerializeParse(ctx context.Context, bufferPtr, sourcePtr, sourceLen, optPtr uint64) (uint64, error) {
	return r.modPmSerializeParse.Call(ctx, bufferPtr, sourcePtr, sourceLen, optPtr)
}

func (r *Runtime) MemoryWrite(ptr uint64, data []byte) bool {
	return r.mem.Write(ptr, data)
}

func (r *Runtime) MemoryRead(offset uint64, byteCount uint64) ([]byte, bool) {
	return r.mem.Read(offset, byteCount)
}
