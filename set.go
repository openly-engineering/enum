package enum

import (
	"fmt"
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

// internalSet collects all enums associated with a specific type T.
type internalSet[T constraints.Integer] struct {
	nextID int64 // Atomically updated.
	ids    map[string]int64
}

// newInternalSet returns a new empty set.
func newInternalSet[T constraints.Integer]() *internalSet[T] {
	return &internalSet[T]{
		0,
		make(map[string]int64),
	}
}

// Add adds a new enum with the given name to the set. The enum ID is
// auto-generated based on the instantiation order of enums. This panics if
// an attempt is made to add an enum with a name that already exists in the
// set.
func (s *internalSet[T]) Add(name string) Enum[T] {
	if _, ok := s.ids[name]; ok {
		panic("duplicate name in set")
	}

	id := atomic.AddInt64(&s.nextID, 1)
	id--

	// TODO(bga): Check for integer overflow as T might be any integer type
	// 			  (say, int8).

	s.ids[name] = id

	return &internalEnum[T]{name: name, id: T(id)}
}

// GetByName returns the Enum associated with the given name and type T.
func (s *internalSet[T]) GetByName(name string) (Enum[T], error) {
	id, ok := s.ids[name]
	if !ok {
		return nil, fmt.Errorf("name %s could not be found in set", name)
	}

	return &internalEnum[T]{name: name, id: T(id)}, nil
}

// GetByID returns the Enum associated with the given ID and type T.
func (s *internalSet[T]) GetByID(id T) (Enum[T], error) {
	for name, v := range s.ids {
		if v == int64(id) {
			return &internalEnum[T]{name: name, id: T(id)}, nil
		}
	}

	return nil, fmt.Errorf("id %d could not be found in set", id)
}
