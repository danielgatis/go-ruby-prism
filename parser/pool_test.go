package parser

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestPoolParsesConcurrently checks that a shared parser runs parses on
// several instances at once rather than serializing them.
func TestPoolParsesConcurrently(t *testing.T) {
	ctx := context.Background()

	const workers = 8

	p, err := NewParser(ctx, WithPoolSize(workers))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	// A source big enough that a parse is not over instantly, so overlapping
	// work is observable.
	source := []byte(strings.Repeat("def m(a, b) = a + b\n", 2000))

	var inFlight, peak int64

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				// Count only while a parse is actually running, so this
				// measures parses in flight rather than goroutines
				// queued up waiting for an instance.
				if _, err := p.parseTracked(ctx, source, &inFlight, &peak); err != nil {
					t.Errorf("parse: %v", err)
				}
			}
		}()
	}

	wg.Wait()

	if peak < 2 {
		t.Errorf("peak concurrent parses = %d, want at least 2; the pool is serializing", peak)
	}
}

// TestPoolSizeOneSerializes checks that a single-instance pool never runs two
// parses at the same time.
func TestPoolSizeOneSerializes(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx, WithPoolSize(1))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	source := []byte(strings.Repeat("x = 1\n", 500))

	var inFlight, peak int64

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				// Count only while a parse is actually running, so this
				// measures parses in flight rather than goroutines
				// queued up waiting for an instance.
				if _, err := p.parseTracked(ctx, source, &inFlight, &peak); err != nil {
					t.Errorf("parse: %v", err)
				}
			}
		}()
	}

	wg.Wait()

	if peak != 1 {
		t.Errorf("peak concurrent parses = %d, want 1", peak)
	}
}

// TestPoolDoesNotExceedItsSize checks that the pool never creates more
// instances than it was configured for.
func TestPoolDoesNotExceedItsSize(t *testing.T) {
	ctx := context.Background()

	const size = 3

	p, err := NewParser(ctx, WithPoolSize(size))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if _, err := p.Parse(ctx, []byte("puts 1")); err != nil {
				t.Errorf("parse: %v", err)
			}
		}()
	}

	wg.Wait()

	p.mutex.Lock()
	created := p.created
	p.mutex.Unlock()

	if created > size {
		t.Errorf("created %d instances, want at most %d", created, size)
	}
}

// TestParseAfterCloseFails checks that a closed parser reports an error
// instead of panicking or blocking.
func TestParseAfterCloseFails(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	if _, err := p.Parse(ctx, []byte("puts 1")); err != nil {
		t.Fatalf("parse before close: %v", err)
	}

	if err := p.Close(ctx); err != nil {
		t.Fatalf("close: %v", err)
	}

	if _, err := p.Parse(ctx, []byte("puts 1")); err == nil {
		t.Error("Parse after Close returned no error")
	}
}

// TestCloseIsIdempotent checks that closing twice is safe.
func TestCloseIsIdempotent(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	if err := p.Close(ctx); err != nil {
		t.Fatalf("first close: %v", err)
	}

	if err := p.Close(ctx); err != nil {
		t.Errorf("second close: %v", err)
	}
}

// TestCloseWaitsForInFlightParses checks that Close does not tear down an
// instance while a parse is still using it.
func TestCloseWaitsForInFlightParses(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx, WithPoolSize(4))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	source := []byte(strings.Repeat("def m(a, b) = a + b\n", 2000))

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < 3; j++ {
				if _, err := p.Parse(ctx, source); err != nil {
					// A parse racing Close may legitimately find the
					// parser closed; a crash or a hang would not be.
					return
				}
			}
		}()
	}

	// Let the parses get going, then close underneath them.
	time.Sleep(20 * time.Millisecond)

	if err := p.Close(ctx); err != nil {
		t.Fatalf("close: %v", err)
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("Close left parses blocked")
	}
}

// TestParseRespectsContextCancellation checks that a parse waiting for a busy
// pool gives up when its context is cancelled.
func TestParseRespectsContextCancellation(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx, WithPoolSize(1))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	// Hold the only instance for the duration of the test.
	held, err := p.acquire(ctx)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	defer p.release(held)

	cancelCtx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()

	if _, err := p.Parse(cancelCtx, []byte("puts 1")); err == nil {
		t.Error("Parse returned no error after its context was cancelled")
	}

	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("Parse took %v to notice cancellation", elapsed)
	}
}

// TestInvalidPoolSizeIsRejected checks that a nonsensical pool size fails
// loudly at construction.
func TestInvalidPoolSizeIsRejected(t *testing.T) {
	ctx := context.Background()

	if _, err := NewParser(ctx, WithPoolSize(0)); err == nil {
		t.Error("NewParser accepted a pool size of 0")
	}

	if _, err := NewParser(ctx, WithPoolSize(-1)); err == nil {
		t.Error("NewParser accepted a negative pool size")
	}
}

// parseTracked runs a parse and records how many parses were executing at the
// same time. It acquires an instance itself so the counter only covers the
// window where the instance is checked out, which is what the pool actually
// serializes.
func (p *Parser) parseTracked(ctx context.Context, source []byte, inFlight, peak *int64) (*ParseResult, error) {
	instance, err := p.acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer p.release(instance)

	current := atomic.AddInt64(inFlight, 1)

	for {
		observed := atomic.LoadInt64(peak)
		if current <= observed || atomic.CompareAndSwapInt64(peak, observed, current) {
			break
		}
	}

	defer atomic.AddInt64(inFlight, -1)

	return p.parseWith(ctx, instance, source)
}
