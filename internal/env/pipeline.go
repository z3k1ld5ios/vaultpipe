package env

// Pipeline chains multiple map transformation steps together,
// applying each in order and propagating errors immediately.
// It is designed to compose the various env sub-packages
// (defaults, override, coerce, transform, etc.) without
// coupling them directly to one another.

// StepFunc is a single transformation step that accepts an
// environment map and returns a (possibly modified) copy.
type StepFunc func(map[string]string) (map[string]string, error)

// Pipeline holds an ordered sequence of StepFuncs.
type Pipeline struct {
	steps []StepFunc
}

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Append adds one or more steps to the end of the pipeline.
func (p *Pipeline) Append(steps ...StepFunc) *Pipeline {
	p.steps = append(p.steps, steps...)
	return p
}

// Run executes every step in order, passing the output of each
// step as the input to the next. The initial input is src.
// If any step returns an error the pipeline halts and returns
// that error together with whatever map was produced so far.
func (p *Pipeline) Run(src map[string]string) (map[string]string, error) {
	current := copyMap(src)
	for _, step := range p.steps {
		result, err := step(current)
		if err != nil {
			return current, err
		}
		current = result
	}
	return current, nil
}

// Len returns the number of steps registered in the pipeline.
func (p *Pipeline) Len() int { return len(p.steps) }

// copyMap returns a shallow copy of m.
func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
