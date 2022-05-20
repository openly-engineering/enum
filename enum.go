package enum

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

// Enum represents a named Enum that is associaterd with an ID. Enum IDs
// are auto-generated starting from 0 and monotonically increasing in
// declaration order.
type Enum[T constraints.Integer] interface {
	Name() string
	ID() T

	// JSON support.
	json.Marshaler
	json.Unmarshaler

	// Text support.
	encoding.TextMarshaler
	encoding.TextUnmarshaler

	// SQL support.
	driver.Valuer
	sql.Scanner

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
// holds the name and typed ID.
type internalEnum[T constraints.Integer] struct {
	name string
	id   T
}

// Name returns the name associated with this Enum instance.
func (e internalEnum[T]) Name() string {
	return e.name
}

// ID returns the numeric ID associated with this Enum instance.
func (e internalEnum[T]) ID() T {
	return e.id
}

// MarshalJSON implements the json.Marshaler interface.
func (e internalEnum[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Name())
}

func getIDForName[T constraints.Integer](name string) (T, error) {
	typeName := getTypeName[T]()

	var defaultID T

	anySet, ok := setByTypeName[typeName]
	if !ok {
		return defaultID, fmt.Errorf("no enum set associated with type %s", typeName)
	}

	s := anySet.(*internalSet[T])

	id, ok := s.ids[name]
	if !ok {
		return defaultID, fmt.Errorf("name %s could not be found in enum set for type %s", name, typeName)
	}

	return T(id), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *internalEnum[T]) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("source should be a string, got %s", data)
	}

	id, err := getIDForName[T](name)
	if err != nil {
		return err
	}

	e.name = name
	e.id = T(id)

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e internalEnum[T]) MarshalText() ([]byte, error) {
	return []byte(e.Name()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *internalEnum[T]) UnmarshalText(text []byte) error {
	name := string(text)

	id, err := getIDForName[T](name)
	if err != nil {
		return err
	}

	e.name = name
	e.id = T(id)

	return nil
}

func (e internalEnum[T]) Value() (driver.Value, error) {
	return e.Name(), nil
}

func (e *internalEnum[T]) Scan(value any) error {
	if value == nil {
		return nil
	}

	name, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a string or byte slice")
		}

		name = string(bytes[:])
	}

	id, err := getIDForName[T](name)
	if err != nil {
		return err
	}

	e.name = name
	e.id = T(id)

	return nil
}

// unimplementanle makes sure we implement the Enum interface.
func (e internalEnum[T]) unimplementable() {}
