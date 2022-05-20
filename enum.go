package enum

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

// Enum represents a named Enum that is associaterd with a value. Enum values
// are auto-generated starting form 0 and monotonically increasing in
// declaration order.
type Enum[T constraints.Integer] interface {
	Name() string
	Value() T

	// Makes sure no external package can implement this interface.
	unimplementable()
}

// We need to use any here because each set will have a different type. This is
// ok though as we will always know the exact type stored and will always
// expose it as the actual type.
var setByTypeName = make(map[string]any)

// getTypeName returns the unique name of the associated type T.
func getTypeName[T any]() string {
	var tInstance T

	tType := reflect.TypeOf(tInstance)

	return tType.PkgPath() + "." + tType.Name()
}

// New returns a new Enum associated with the given name and type T. 
func New[T constraints.Integer](name string) Enum[T] {
	typeName := getTypeName[T]()

	var s *internalSet[T]

	if as, ok := setByTypeName[typeName]; !ok {
		s = newInternalSet[T]()
		setByTypeName[typeName] = s
	} else {
		s = as.(*internalSet[T])
	}

	return s.Add(name)
}

// internalEnum is the only concrete implementation of the Enum interface. It
// holds the name and typed value.
type internalEnum[T comparable] struct {
	name  string
	value T
}

// Name returns the name associated with this Enum instance.
func (e internalEnum[T]) Name() string {
	return e.name
}

// Value returns the value associated with this Enum instance.
func (e internalEnum[T]) Value() T {
	return e.value
}

// unimplementanle makes sure we implement the Enum interface.
func (e internalEnum[T]) unimplementable() {}
