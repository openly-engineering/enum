package enum

import (
	"fmt"
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

// internalSet collects all enums associated with a specific type T.
type internalSet[T constraints.Integer] struct {
	nextValue int64 // Atomically updated.
	values    map[string]int64
}

// newInternalSet returns a new empty set.
func newInternalSet[T constraints.Integer]() *internalSet[T] {
	return &internalSet[T]{
		0,
		make(map[string]int64),
	}
}

// Add adds a new enum with the given name to the set. The enum value is
// auto-generated based on the instantiation order of enums. This panics if
// an attempt is made to add an enum with a name that already exists in the
// set.
func (s *internalSet[T]) Add(name string) Enum[T] {
	if _, ok := s.values[name]; ok {
		panic("duplicate name in set")
	}

	value := atomic.AddInt64(&s.nextValue, 1)
	value--

	// TODO(bga): Check for integer overflow as T might be any integer type
	// 			  (say, int8).

	s.values[name] = value

	return internalEnum[T]{name: name, value: T(value)}
}

// GetByName returns the Enum associated with the given name and type T.
func (s internalSet[T]) GetByName(name string) (Enum[T], error) {
	value, ok := s.values[name]
	if !ok {
		return nil, fmt.Errorf("name %s could not be found in set", name)
	}

	return internalEnum[T]{name: name, value: T(value)}, nil
}

// GetByValue returns the Enum associated with the given value and type T.
func (s internalSet[T]) GetByValue(value int64) (Enum[T], error) {
	for name, v := range s.values {
		if v == value {
			return internalEnum[T]{name: name, value: T(value)}, nil
		}
	}

	return nil, fmt.Errorf("value %d could not be found in set", value)
}
