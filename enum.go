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
// declaration order. The zero value of an Enum is not valid. It is safe
// to use this type to create other types (type OtherType Enum[MyEnumType]) as
// it does not implement any methods itself and, instead, delegate all
// methods to embedded types.
type Enum[T constraints.Integer] struct {
	// As internalEnumWrapper is not a pointer, it will never be nil so we use
	// it to implement all methods that we need.
	internalEnumWrapper[T]
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

func getOrCreateSetForType[T constraints.Integer]() *internalSet[T] {
	typeName := getTypeName[T]()

	var s *internalSet[T]

	if as, ok := setByTypeName[typeName]; !ok {
		s = newInternalSet[T]()
		setByTypeName[typeName] = s
	} else {
		s = as.(*internalSet[T])
	}

	return s
}

// New returns a new Enum associated with the given name and type T.
func New[T constraints.Integer](name string) Enum[T] {
	if name == "" {
		panic("enum name cannot be empty")
	}

	s := getOrCreateSetForType[T]()

	return Enum[T]{internalEnumWrapper[T]{s.Add(name)}}
}

// EnumsByType returns all enums associated with the given type T.
func EnumsByType[T constraints.Integer]() []Enum[T] {
	s := setByTypeName[getTypeName[T]()]

	nameEnumMap := s.(*internalSet[T]).nameEnumMap

	enums := make([]Enum[T], 0, len(nameEnumMap))
	for _, e := range nameEnumMap {
		enums = append(enums, Enum[T]{internalEnumWrapper[T]{e}})
	}

	return enums
}

// EnumByTypeAndName returns the enum associated with the given type and name.
// If there is no such enum, a non-nil error is returned.
func EnumByTypeAndName[T constraints.Integer](name string) (Enum[T], error) {
	e, err := getInternalEnumForName[T](name)
	if err != nil {
		return Enum[T]{}, err
	}

	return Enum[T]{internalEnumWrapper[T]{e}}, nil
}

// internalEnumWrapper is the type that implements all Enum methods.
type internalEnumWrapper[T constraints.Integer] struct {
	*internalEnum[T]
}

// Name returns the name associated with this Enum instance.
func (e internalEnumWrapper[T]) Name() string {
	if !e.Valid() {
		panic("enum not initialized")
	}

	return e.internalEnum.name
}

// ID returns the numeric ID associated with this Enum instance.
func (e internalEnumWrapper[T]) ID() T {
	if !e.Valid() {
		panic("enum not initialized")
	}

	return e.internalEnum.id
}

// Valid returns true if the Enum is valid or false otherwise. Default Enum
// instances are invalid. Use New to create a valid one (or use the
// unmarshalling methods to initialize one created in place).
func (e *internalEnumWrapper[T]) Valid() bool {
	return e.internalEnum != nil
}

// MarshalJSON implements the json.Marshaler interface.
func (e internalEnumWrapper[T]) MarshalJSON() ([]byte, error) {
	if !e.Valid() {
		return nil, fmt.Errorf("enum not initialized")
	}

	return json.Marshal(e.Name())
}

func getInternalEnumForName[T constraints.Integer](name string) (*internalEnum[T], error) {
	typeName := getTypeName[T]()

	anySet, ok := setByTypeName[typeName]
	if !ok {
		return nil, fmt.Errorf("no enum set associated with type %s", typeName)
	}

	s := anySet.(*internalSet[T])

	var e *internalEnum[T]
	if e = s.Get(name); e == nil {
		return nil, fmt.Errorf("name %s could not be found in enum set for type %s", name, typeName)
	}

	return e, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *internalEnumWrapper[T]) UnmarshalJSON(data []byte) error {
	var name string
	var err error

	if err = json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("source should be a string, got %s", data)
	}

	e.internalEnum, err = getInternalEnumForName[T](name)
	if err != nil {
		return err
	}

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e internalEnumWrapper[T]) MarshalText() ([]byte, error) {
	if !e.Valid() {
		return nil, fmt.Errorf("enum not initialized")
	}

	return []byte(e.Name()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *internalEnumWrapper[T]) UnmarshalText(text []byte) error {
	name := string(text)

	var err error
	e.internalEnum, err = getInternalEnumForName[T](name)
	if err != nil {
		return err
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (e internalEnumWrapper[T]) Value() (driver.Value, error) {
	if !e.Valid() {
		return nil, fmt.Errorf("enum not initialized")
	}

	return e.Name(), nil
}

// Scan implements the sql.Scanner interface.
func (e *internalEnumWrapper[T]) Scan(value any) error {
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

	var err error
	e.internalEnum, err = getInternalEnumForName[T](name)
	if err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (e internalEnumWrapper[T]) String() string {
	if !e.Valid() {
		panic("enum not initialized")
	}

	return e.name
}

// internalEnum is the internal representation of an Enum and is the type that
// stores the Enum-associated data.
type internalEnum[T constraints.Integer] struct {
	name string
	id   T
}
