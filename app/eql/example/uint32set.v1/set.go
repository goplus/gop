package uint32set

// ---------------------------------------------------------------------------

type empty struct{}

// Type is a set of uint32, implemented via map[uint32]struct{} for minimal memory consumption.
type Type map[uint32]empty

// New creates a Type from a list of values.
func New(items ...uint32) Type {
	ss := Type{}
	ss.Insert(items...)
	return ss
}

// Add adds one item to the set.
func (s Type) Add(item uint32) {
	s[item] = empty{}
}

// Insert adds items to the set.
func (s Type) Insert(items ...uint32) {
	for _, item := range items {
		s[item] = empty{}
	}
}

// Delete removes all items from the set.
func (s Type) Delete(items ...uint32) {
	for _, item := range items {
		delete(s, item)
	}
}

// Has returns true iff item is contained in the set.
func (s Type) Has(item uint32) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true iff all items are contained in the set.
func (s Type) HasAll(items ...uint32) bool {
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
