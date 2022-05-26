package enum

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

// Enum represents a named Enum that is associaterd with an ID. Enum IDs
// are auto-generated starting from 0 and monotonically increasing in
// declaration order. The zero value of an Enum is not valid.
type Enum[T constraints.Integer] struct {
	*internalEnum[T]
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
func New[T constraints.Integer](name string) *Enum[T] {
	if name == "" {
		panic("enum name cannot be empty")
	}

	typeName := getTypeName[T]()

	var s *internalSet[T]

	if as, ok := setByTypeName[typeName]; !ok {
		s = newInternalSet[T]()
		setByTypeName[typeName] = s
	} else {
		s = as.(*internalSet[T])
	}

	return &Enum[T]{s.Add(name)}
}

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

// Valid returns true if the Enum is valid or false otherwise. Default Enum
// instances are invalid. Use New to create a valid one.
func (e *internalEnum[T]) Valid() bool {
	return e != nil
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

	e, ok := s.nameEnumMap[name]
	if !ok {
		return defaultID, fmt.Errorf("name %s could not be found in enum set for type %s", name, typeName)
	}

	return e.id, nil
}

func (e *Enum[T]) createOrUpdateInternalEnum(id T, name string) {
	if e.internalEnum == nil {
		e.internalEnum = &internalEnum[T]{name, id}
	} else {
		// This is just to avoid garbage collection pressure. If we already have
		// an internalEnum, we can just update it.
		e.name = name
		e.id = T(id)
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Enum[T]) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("source should be a string, got %s", data)
	}

	id, err := getIDForName[T](name)
	if err != nil {
		return err
	}

	e.createOrUpdateInternalEnum(id, name)

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e internalEnum[T]) MarshalText() ([]byte, error) {
	return []byte(e.Name()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *Enum[T]) UnmarshalText(text []byte) error {
	name := string(text)

	id, err := getIDForName[T](name)
	if err != nil {
		return err
	}

	e.createOrUpdateInternalEnum(id, name)

	return nil
}

// Value implements the driver.Valuer interface.
func (e internalEnum[T]) Value() (driver.Value, error) {
	return e.Name(), nil
}

// Scan implements the sql.Scanner interface.
func (e *Enum[T]) Scan(value any) error {
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

	e.createOrUpdateInternalEnum(id, name)

	return nil
}

// String implements the fmt.Stringer interface.
func (e internalEnum[T]) String() string {
	return e.name
}
