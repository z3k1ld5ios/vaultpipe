package env

// ChainStep is a function that transforms an environment map.
// It receives the current map and returns a new map or an error.
type ChainStep func(map[string]string) (map[string]string, error)

// Chain applies a sequence of ChainSteps in order, threading the output
// of each step into the next. If any step returns an error the chain
// halts immediately and that error is returned.
//
// Chain is intentionally thin — it does not copy the initial map itself;
// individual steps are responsible for immutability if required.
type Chain struct {
	steps []ChainStep
}

// NewChain constructs a Chain from the provided steps.
func NewChain(steps ...ChainStep) *Chain {
	return &Chain{steps: steps}
}

// Apply runs every step in sequence starting from env.
// The returned map is the result of the final step.
func (c *Chain) Apply(env map[string]string) (map[string]string, error) {
	current := env
	for _, step := range c.steps {
		result, err := step(current)
		if err != nil {
			return nil, err
		}
		current = result
	}
	return current, nil
}

// Append returns a new Chain with the given steps appended.
func (c *Chain) Append(steps ...ChainStep) *Chain {
	next := make([]ChainStep, len(c.steps)+len(steps))
	copy(next, c.steps)
	copy(next[len(c.steps):], steps)
	return &Chain{steps: next}
}

// Len returns the number of steps in the chain.
func (c *Chain) Len() int { return len(c.steps) }
