package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
)

// distinctSources have deliberately different lengths so that each parse
// requests a differently sized allocation from the WASM heap.
var distinctSources = []string{
	`class Foo; def bar; 1; end; end`,
	`module A
  module B
    class C
      def initialize(x, y)
        @x = x
        @y = y
      end
    end
  end
end`,
	`puts 1`,
	`x = [1,2,3].map { |i| i * 2 }.select { |i| i > 2 }
y = x.reduce(0) { |a, b| a + b }
puts y`,
	`def f(a, b = 2, *rest, k:, **opts, &blk)
  yield a
end`,
	`# frozen_string_literal: true
require "set"
S = Set.new([1,2,3])`,
	``,
	`begin
  raise ArgumentError, "x"
rescue => e
  retry
ensure
  puts "done"
end`,
	strings.Repeat("a = 1\n", 500),
	`case v
in [1, *rest] then rest
in {k: Integer => n} then n
else 0
end`,
}

// TestParserReuseAcrossDistinctSources is a regression test for the WASM
// out-of-bounds trap reported in issue #4. Reusing a single parser across
// sources of differing sizes used to trap roughly half the time because the
// options blob was read out of memory that still held bytes from the previous
// parse.
func TestParserReuseAcrossDistinctSources(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	for round := 0; round < 30; round++ {
		for i, src := range distinctSources {
			if _, err := p.Parse(ctx, []byte(src)); err != nil {
				t.Fatalf("round %d, source %d: %v", round, i, err)
			}
		}
	}
}

// TestParserReuseWithFilePath covers the same loop with a non-empty filepath,
// which takes the other branch of pm_options_read.
func TestParserReuseWithFilePath(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx, WithFilePath("test.rb"))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	for round := 0; round < 30; round++ {
		for i, src := range distinctSources {
			if _, err := p.Parse(ctx, []byte(src)); err != nil {
				t.Fatalf("round %d, source %d: %v", round, i, err)
			}
		}
	}
}

// TestParserReuseMatchesFreshParser guards against the weaker failure mode:
// no trap, but a tree built from stale options.
func TestParserReuseMatchesFreshParser(t *testing.T) {
	ctx := context.Background()

	reused, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer reused.Close(ctx)

	// Warm the reused parser so it is not on its first allocation.
	for _, src := range distinctSources {
		if _, err := reused.Parse(ctx, []byte(src)); err != nil {
			t.Fatalf("warmup: %v", err)
		}
	}

	for i, src := range distinctSources {
		fresh, err := NewParser(ctx)
		if err != nil {
			t.Fatalf("failed to create fresh parser: %v", err)
		}

		want, err := fresh.Parse(ctx, []byte(src))
		fresh.Close(ctx)
		if err != nil {
			t.Fatalf("source %d: fresh parser failed: %v", i, err)
		}

		got, err := reused.Parse(ctx, []byte(src))
		if err != nil {
			t.Fatalf("source %d: reused parser failed: %v", i, err)
		}

		wantJSON, err := json.Marshal(want.Value)
		if err != nil {
			t.Fatalf("source %d: %v", i, err)
		}

		gotJSON, err := json.Marshal(got.Value)
		if err != nil {
			t.Fatalf("source %d: %v", i, err)
		}

		if string(wantJSON) != string(gotJSON) {
			t.Errorf("source %d: reused parser produced a different AST than a fresh parser", i)
		}

		if len(want.Errors) != len(got.Errors) {
			t.Errorf("source %d: error count %d (fresh) != %d (reused)", i, len(want.Errors), len(got.Errors))
		}
	}
}

// TestParseEmptySource exercises the zero-length allocation path directly.
func TestParseEmptySource(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	for i := 0; i < 20; i++ {
		if _, err := p.Parse(ctx, []byte("")); err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
	}
}

// TestParseReportsSyntaxErrors makes sure the cleanup changes did not swallow
// diagnostics coming back from Prism.
func TestParseReportsSyntaxErrors(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	result, err := p.Parse(ctx, []byte("def broken("))
	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}

	if len(result.Errors) == 0 {
		t.Fatal("expected syntax errors, got none")
	}
}

// TestParserConcurrentUse verifies the mutex actually serializes Parse, since
// the README documents a shared parser as safe to call from many goroutines.
func TestParserConcurrentUse(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	var wg sync.WaitGroup
	errs := make(chan error, 16*len(distinctSources))

	for g := 0; g < 16; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, src := range distinctSources {
				if _, err := p.Parse(ctx, []byte(src)); err != nil {
					errs <- fmt.Errorf("source %d: %w", i, err)
				}
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}
