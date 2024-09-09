package set

import (
	"sync"
)

// BitwiseSet uses a bitset to track the presence of uint32 elements
type BitwiseSet struct {
	mu     sync.Mutex // Protect concurrent access
	bitset []uint32   // The bitset, using 32-bit words
}

// NewBitwiseSet initializes a new BitwiseSet
func NewBitwiseSet(size int) *BitwiseSet {
	// We need enough space to represent each possible bit, so divide the total number of bits by 32
	bitsetSize := (size / 32) + 1
	return &BitwiseSet{
		bitset: make([]uint32, bitsetSize), // Allocate space for the bitset
	}
}

// Add marks the bit for the given element using bitwise operations.
func (s *BitwiseSet) Add(element uint32) {
	// Calculate the index in the bitset and the bit position
	index := element / 32
	bit := uint32(1 << (element % 32)) // Compute the bit to set

	// Lock, set the bit, and unlock
	s.mu.Lock()
	s.bitset[index] |= bit // Set the bit
	s.mu.Unlock()
}

// Contains checks if the bit for the given element is set.
func (s *BitwiseSet) Contains(element uint32) bool {
	index := element / 32
	bit := uint32(1 << (element % 32))

	// Lock, check the bit, and unlock
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.bitset[index]&bit != 0 // Return true if the bit is set
}

// Count returns the number of set bits (distinct elements in the set)
func (s *BitwiseSet) Count() int {
	count := 0

	// Lock the bitset for reading
	s.mu.Lock()
	defer s.mu.Unlock()

	// Count the number of set bits in the entire bitset
	for _, word := range s.bitset {
		count += popCount(word) // Use popCount to count the set bits in each word
	}

	return count
}

// popCount counts the number of set bits (1s) in uint32 word
func popCount(word uint32) int {
	count := 0
	for word != 0 {
		word &= word - 1 // Clear the lowest set bit
		count++
	}
	return count
}
