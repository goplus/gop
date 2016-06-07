package stringset

// ---------------------------------------------------------------------------

type empty struct{}

// StringSet is a set of string, implemented via map[string]struct{} for minimal memory consumption.
type StringSet map[string]empty

// New creates a StringSet from a list of values.
func New(items ...string) StringSet {
	ss := StringSet{}
	ss.Insert(items...)
	return ss
}

// Add adds one item to the set.
func (s StringSet) Add(item string) {
	s[item] = empty{}
}

// Insert adds items to the set.
func (s StringSet) Insert(items ...string) {
	for _, item := range items {
		s[item] = empty{}
	}
}

// Delete removes all items from the set.
func (s StringSet) Delete(items ...string) {
	for _, item := range items {
		delete(s, item)
	}
}

// Has returns true iff item is contained in the set.
func (s StringSet) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true iff all items are contained in the set.
func (s StringSet) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// IsSuperset returns true iff s1 is a superset of s2.
func (s StringSet) IsSuperset(s2 StringSet) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// Len returns the size of the set.
func (s StringSet) Len() int {
	return len(s)
}

// ---------------------------------------------------------------------------
