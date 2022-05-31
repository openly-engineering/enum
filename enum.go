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

// EnumsForType returns all enums associated with the given type T.
func EnumsForType[T constraints.Integer]() []*Enum[T] {
	s := setByTypeName[getTypeName[T]()]

	nameEnumMap := s.(*internalSet[T]).nameEnumMap

	enums := make([]*Enum[T], 0, len(nameEnumMap))
	for _, e := range nameEnumMap {
		enums = append(enums, &Enum[T]{e})
	}

	return enums
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

	return Enum[T]{s.Add(name)}
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
func (e *Enum[T]) UnmarshalJSON(data []byte) error {
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
func (e internalEnum[T]) MarshalText() ([]byte, error) {
	return []byte(e.Name()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *Enum[T]) UnmarshalText(text []byte) error {
	name := string(text)

	var err error
	e.internalEnum, err = getInternalEnumForName[T](name)
	if err != nil {
		return err
	}

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

	var err error
	e.internalEnum, err = getInternalEnumForName[T](name)
	if err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (e internalEnum[T]) String() string {
	return e.name
}
