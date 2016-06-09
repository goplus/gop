package stringset

// ---------------------------------------------------------------------------

type empty struct{}

// Type is a set of string, implemented via map[string]struct{} for minimal memory consumption.
type Type map[string]empty

// New creates a Type from a list of values.
func New(items ...string) Type {
	ss := Type{}
	ss.Insert(items...)
	return ss
}

// Add adds one item to the set.
func (s Type) Add(item string) {
	s[item] = empty{}
}

// Insert adds items to the set.
func (s Type) Insert(items ...string) {
	for _, item := range items {
		s[item] = empty{}
	}
}

// Delete removes all items from the set.
func (s Type) Delete(items ...string) {
	for _, item := range items {
		delete(s, item)
	}
}

// Has returns true iff item is contained in the set.
func (s Type) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true iff all items are contained in the set.
func (s Type) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// IsSuperset returns true iff s1 is a superset of s2.
func (s Type) IsSuperset(s2 Type) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// Len returns the size of the set.
func (s Type) Len() int {
	return len(s)
}

// ---------------------------------------------------------------------------
