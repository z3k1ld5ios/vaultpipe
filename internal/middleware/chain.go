package middleware

import "context"

// Guard is a function that inspects a context before a Vault operation proceeds.
type Guard func(ctx context.Context) error

// Chain holds an ordered list of Guards applied before each Vault call.
type Chain struct {
	guards []Guard
}

// NewChain creates a Chain from the provided guards.
func NewChain(guards ...Guard) *Chain {
	return &Chain{guards: guards}
}

// Run executes all guards in order. Returns the first error encountered.
func (c *Chain) Run(ctx context.Context) error {
	for _, g := range c.guards {
		if err := g(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RateLimitGuard adapts a RateLimiter into a Guard.
func RateLimitGuard(rl *RateLimiter) Guard {
	return func(ctx context.Context) error {
		return rl.Allow(ctx)
	}
}
