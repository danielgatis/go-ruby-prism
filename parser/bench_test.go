package parser

import (
	"context"
	"strings"
	"testing"
)

var benchSource = []byte(strings.Repeat("def m(a, b) = a + b\n", 500))

// BenchmarkParseSerial parses on a single instance.
func BenchmarkParseSerial(b *testing.B) {
	ctx := context.Background()

	p, err := NewParser(ctx, WithPoolSize(1))
	if err != nil {
		b.Fatal(err)
	}
	defer p.Close(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := p.Parse(ctx, benchSource); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseParallel parses from many goroutines against a pooled parser,
// which is the case the pool exists for.
func BenchmarkParseParallel(b *testing.B) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		b.Fatal(err)
	}
	defer p.Close(ctx)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := p.Parse(ctx, benchSource); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkParseParallelPoolSizeOne is the same workload against a
// single-instance pool, for comparison with BenchmarkParseParallel.
func BenchmarkParseParallelPoolSizeOne(b *testing.B) {
	ctx := context.Background()

	p, err := NewParser(ctx, WithPoolSize(1))
	if err != nil {
		b.Fatal(err)
	}
	defer p.Close(ctx)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := p.Parse(ctx, benchSource); err != nil {
				b.Fatal(err)
			}
		}
	})
}
