package signal_test

import (
	"context"
	"errors"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/yourusername/vaultpipe/internal/signal"
)

func newTestHandler(t *testing.T) *signal.Handler {
	t.Helper()
	log, _ := zap.NewDevelopment()
	return signal.NewHandler(log)
}

func TestWait_ContextCancelled_RunsCallbacks(t *testing.T) {
	h := newTestHandler(t)

	var called int32
	h.Register(func(ctx context.Context) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		h.Wait(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after context cancellation")
	}

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("expected 1 callback call, got %d", called)
	}
}

func TestWait_MultipleCallbacks_AllInvoked(t *testing.T) {
	h := newTestHandler(t)

	var count int32
	for i := 0; i < 3; i++ {
		h.Register(func(ctx context.Context) error {
			atomic.AddInt32(&count, 1)
			return nil
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		h.Wait(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return")
	}

	if atomic.LoadInt32(&count) != 3 {
		t.Errorf("expected 3 callback calls, got %d", count)
	}
}

func TestWait_CallbackError_ContinuesOthers(t *testing.T) {
	h := newTestHandler(t)

	var second int32
	h.Register(func(ctx context.Context) error {
		return errors.New("shutdown failure")
	})
	h.Register(func(ctx context.Context) error {
		atomic.AddInt32(&second, 1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		h.Wait(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return")
	}

	if atomic.LoadInt32(&second) != 1 {
		t.Error("second callback was not called despite first returning error")
	}
}

func TestWait_SIGTERMReceived_RunsCallbacks(t *testing.T) {
	h := newTestHandler(t)

	var called int32
	h.Register(func(ctx context.Context) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	ctx := context.Background()
	done := make(chan struct{})
	go func() {
		h.Wait(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after SIGTERM")
	}

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("expected callback to be called, got %d", called)
	}
}
