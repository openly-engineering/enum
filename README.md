# enum
Easy to use generic named enum library for go.

## Introduction

This library implements support for named enums without requiring code generation this is achieved by the obvious approach of requiring the user to give enum names (intead of values) and values (IDs in the context of this library) are auto-generated per associated type starting from 0 and monotonicaly increasing in declaration order.

## Creating Enums

Here is how enums are usually created with Go:

```
type MyType int

const (
    Unknown MyType = iota  // 0
    One                    // 1
    Two                    // 2
)
```

In the simple example above, enums have no associated names whatsoever and one would have to manually keep track of names if they needed to use those. Libraries like [enumer](https://github.com/alvaroloes/enumer) try to solve that by using code generation to add support for keeping track of enum names (plus other niceties).

And here is how they would be created with this library:

```
type MyType int

var (
    Unknown = enum.New[MyType]("Unknown")  // 0
    One     = enum.New[MyType]("One")      // 1
    Two     = enum.New[MyType]("Two")      // 2
)
```

Here each enum is associated with a specific type (MyType), has a specific name passed as a parameter to the New call and has an associated auto generated ID.

## Using Enums

Below you can find some possibly interesting usage scenarios in the context of enums.

Getting Enum name:
```
unknownName := Unknown.Name()  // "Unknown"
```

Getting Enum ID/value:
```
unknownID := Unknown.ID()  // 0
```

TODO(bga): Finish this.
