package env

// TracedStep wraps any PipelineStep and emits a TraceEntry for every key
// whose value changes after the inner step executes.
type TracedStep struct {
	name  string
	inner PipelineStep
	tracer *Tracer
}

// PipelineStep is the shared interface used by Pipeline and Chain.
type PipelineStep interface {
	Apply(env map[string]string) (map[string]string, error)
}

// NewTracedStep wraps inner with tracing under the given step name.
func NewTracedStep(name string, inner PipelineStep, tr *Tracer) *TracedStep {
	return &TracedStep{name: name, inner: inner, tracer: tr}
}

// Apply executes the inner step, then records a TraceEntry for every key
// that was added, removed, or modified.
func (ts *TracedStep) Apply(env map[string]string) (map[string]string, error) {
	result, err := ts.inner.Apply(env)
	if err != nil {
		ts.tracer.Record(ts.name, "", "", "", TraceLevelError, err.Error())
		return nil, err
	}

	// Detect changes: modified or removed keys.
	for k, oldVal := range env {
		newVal, ok := result[k]
		if !ok {
			ts.tracer.Record(ts.name, k, oldVal, "", TraceLevelInfo, "removed")
			continue
		}
		if newVal != oldVal {
			ts.tracer.Record(ts.name, k, oldVal, newVal, TraceLevelInfo, "modified")
		}
	}

	// Detect added keys.
	for k, newVal := range result {
		if _, existed := env[k]; !existed {
			ts.tracer.Record(ts.name, k, "", newVal, TraceLevelInfo, "added")
		}
	}

	return result, nil
}
