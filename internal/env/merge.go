package env

// MergeStrategy controls how conflicting keys are resolved when merging
// two environment maps together.
type MergeStrategy int

const (
	// StrategySecretWins causes secret values to override base env values.
	StrategySecretWins MergeStrategy = iota
	// StrategyBaseWins causes existing base env values to be preserved.
	StrategyBaseWins
	// StrategyError causes an error to be returned on any key conflict.
	StrategyError
)

// Merger combines a base environment map with a secrets map according to
// a configurable conflict resolution strategy.
type Merger struct {
	strategy MergeStrategy
}

// NewMerger returns a Merger configured with the given strategy.
func NewMerger(strategy MergeStrategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines base and secrets into a single map. The returned map is
// always a new allocation and neither input is mutated.
func (m *Merger) Merge(base, secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(base)+len(secrets))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range secrets {
		if existing, conflict := out[k]; conflict {
			switch m.strategy {
			case StrategyBaseWins:
				// keep existing value, skip secret
				_ = existing
				continue
			case StrategyError:
				return nil, &MergeConflictError{Key: k}
			default: // StrategySecretWins
				out[k] = v
			}
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// MergeConflictError is returned when StrategyError is in use and a key
// appears in both the base and secrets maps.
type MergeConflictError struct {
	Key string
}

func (e *MergeConflictError) Error() string {
	return "env/merge: conflicting key: " + e.Key
}
