package parser

import (
	"context"
	"testing"

	"github.com/tetratelabs/wazero"
)

// A parser built with a compilation cache parses like one without; the second
// parser sharing the cache instantiates from it.
func TestWithRuntimeConfig(t *testing.T) {
	ctx := context.Background()
	cache := wazero.NewCompilationCache()
	defer cache.Close(ctx)
	config := wazero.NewRuntimeConfig().WithCompilationCache(cache)

	for i := 0; i < 2; i++ {
		p, err := NewParser(ctx, WithPoolSize(1), WithRuntimeConfig(config))
		if err != nil {
			t.Fatalf("NewParser: %v", err)
		}
		result, err := p.Parse(ctx, []byte("1 + 1"))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if result == nil || result.Value == nil {
			t.Fatal("no program node")
		}
		if err := p.Close(ctx); err != nil {
			t.Fatalf("Close: %v", err)
		}
	}
}
