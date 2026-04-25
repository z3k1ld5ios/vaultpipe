package env

// ChangeType describes the kind of change detected between two env maps.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change represents a single key-level difference between two env maps.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Diff computes the ordered list of changes between two env maps.
// Keys present in next but not in prev are Added.
// Keys present in prev but not in next are Removed.
// Keys present in both but with different values are Updated.
func Diff(prev, next map[string]string) []Change {
	var changes []Change

	for k, nv := range next {
		if ov, ok := prev[k]; !ok {
			changes = append(changes, Change{Key: k, Type: ChangeAdded, NewVal: nv})
		} else if ov != nv {
			changes = append(changes, Change{Key: k, Type: ChangeUpdated, OldVal: ov, NewVal: nv})
		}
	}

	for k, ov := range prev {
		if _, ok := next[k]; !ok {
			changes = append(changes, Change{Key: k, Type: ChangeRemoved, OldVal: ov})
		}
	}

	sortChanges(changes)
	return changes
}

// HasChanges returns true when Diff finds at least one difference.
func HasChanges(prev, next map[string]string) bool {
	return len(Diff(prev, next)) > 0
}

// FilterByType returns only the changes that match the given ChangeType.
func FilterByType(changes []Change, ct ChangeType) []Change {
	out := make([]Change, 0, len(changes))
	for _, c := range changes {
		if c.Type == ct {
			out = append(out, c)
		}
	}
	return out
}

func sortChanges(changes []Change) {
	// insertion sort — maps are small in practice
	for i := 1; i < len(changes); i++ {
		for j := i; j > 0 && changes[j].Key < changes[j-1].Key; j-- {
			changes[j], changes[j-1] = changes[j-1], changes[j]
		}
	}
}
