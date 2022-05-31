package enum

import (
	"fmt"
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

// internalSet collects all enums associated with a specific type T.
type internalSet[T constraints.Integer] struct {
	nameEnumMap map[string]*internalEnum[T]

	nextID      int64 // Atomically updated.
	exhaustedID bool  // Set to true when there are no more IDs available.
}

// newInternalSet returns a new empty set.
func newInternalSet[T constraints.Integer]() *internalSet[T] {
	return &internalSet[T]{
		make(map[string]*internalEnum[T]),
		0,
		false,
	}
}

// Add adds a new enum with the given name to the set. The enum ID is
// auto-generated based on the instantiation order of enums. This panics if
// an attempt is made to add an enum with a name that already exists in the
// set.
func (s *internalSet[T]) Add(name string) *internalEnum[T] {
	if s.exhaustedID {
		// Run out of IDs.
		panic("too many enums in enum set")
	}

	if _, ok := s.nameEnumMap[name]; ok {
		panic("duplicate name in enum set")
	}

	// Reserve one ID for us and update nextID.
	id := atomic.AddInt64(&s.nextID, 1)
	newID := id - 1

	if T(newID) > T(id) {
		// As we always increment by one, it is guaranteed that we will see the
		// moment id wraps around. If Add() is being called by multiple threads,
		// it is possible that some of those threads will not notice the wrap
		// around but this does not matter as some other thread is still
		// guaranteed to hit this panic here.
		//
		// We mark IDs as exhausthed as the one we just generated is valid.
		s.exhaustedID = true
	}

	e := &internalEnum[T]{
		name: name,
		id:   T(newID),
	}

	s.nameEnumMap[name] = e

	return e
}

// Get returns the enum associated with the given name. If no enum with the
// given name exists, this returns nil.
func (s *internalSet[T]) Get(name string) *internalEnum[T] {
	e, ok := s.nameEnumMap[name]
	if !ok {
		return nil
	}

	return e
}

// GetByName returns the Enum associated with the given name and type T.
func (s *internalSet[T]) GetByName(name string) (*internalEnum[T], error) {
	e, ok := s.nameEnumMap[name]
	if !ok {
		return nil, fmt.Errorf("name %s could not be found in set", name)
	}

	return e, nil
}

// GetByID returns the Enum associated with the given ID and type T.
func (s *internalSet[T]) GetByID(id T) (*internalEnum[T], error) {
	for _, e := range s.nameEnumMap {
		if e.id == id {
			return e, nil
		}
	}

	return nil, fmt.Errorf("id %d could not be found in set", id)
}
