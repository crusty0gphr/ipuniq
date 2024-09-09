package ipuniq

import (
	"sync"
)

// Set using sync.Map for concurrency
type Set struct {
	elements sync.Map
}

// NewSet initializes a new Set
func NewSet() *Set {
	return &Set{}
}

// Add inserts an element into the set
func (s *Set) Add(element uint32) {
	s.elements.Store(element, struct{}{})
}

// Count returns the number of unique elements in the set
func (s *Set) Count() int {
	count := 0
	s.elements.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
