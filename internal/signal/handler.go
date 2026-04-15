// Package signal provides graceful shutdown handling for vaultpipe,
// ensuring leases are revoked and resources cleaned up on termination.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// ShutdownFunc is a callback invoked when a termination signal is received.
type ShutdownFunc func(ctx context.Context) error

// Handler listens for OS signals and triggers registered shutdown callbacks.
type Handler struct {
	log       *zap.Logger
	callbacks []ShutdownFunc
	signals   []os.Signal
}

// NewHandler creates a Handler that responds to SIGINT and SIGTERM by default.
func NewHandler(log *zap.Logger) *Handler {
	return &Handler{
		log:     log,
		signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	}
}

// Register adds a shutdown callback to be called on signal receipt.
func (h *Handler) Register(fn ShutdownFunc) {
	h.callbacks = append(h.callbacks, fn)
}

// Wait blocks until a termination signal is received, then runs all
// registered callbacks in order. The provided context is passed to each.
func (h *Handler) Wait(ctx context.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.signals...)
	defer signal.Stop(ch)

	select {
	case sig := <-ch:
		h.log.Info("signal received, shutting down", zap.String("signal", sig.String()))
	case <-ctx.Done():
		h.log.Info("context cancelled, shutting down")
	}

	for _, fn := range h.callbacks {
		if err := fn(ctx); err != nil {
			h.log.Error("shutdown callback error", zap.Error(err))
		}
	}
}
