Go+ Specification
=====

## Comments

See [Comments](mini-spec.md#comments).


## Literals

See [Literals](mini-spec.md#literals).

### String literals

#### C style string literals

TODO

```go
c"Hello, world!\n"
```

#### Python string literals

TODO

```go
py"Hello, world!\n"
```

## Types

### Boolean types

See [Boolean types](mini-spec.md#boolean-types).

### Numeric types

See [Numeric types](mini-spec.md#numeric-types).

#### String types

See [String types](mini-spec.md#string-types).

##### C style string types

```go
import "c"

*c.Char  // alias for *int8
```

##### Python string types

```go
import "py"

*py.Object  // TODO: *py.String?
```

### Array types

See [Array types](mini-spec.md#array-types).

An array type T may not have an element of type T, or of a type containing T as a component, directly or indirectly, if those containing types are only array or struct types.

```go
// invalid array types
type (
	T1 [10]T1                 // element type of T1 is T1
	T2 [10]struct{ f T2 }     // T2 contains T2 as component of a struct
	T3 [10]T4                 // T3 contains T3 as component of a struct in T4
	T4 struct{ f T3 }         // T4 contains T4 as component of array T3 in a struct
)

// valid array types
type (
	T5 [10]*T5                // T5 contains T5 as component of a pointer
	T6 [10]func() T6          // T6 contains T6 as component of a function type
	T7 [10]struct{ f []T7 }   // T7 contains T7 as component of a slice in a struct
)
```

### Slice types

See [Slice types](mini-spec.md#slice-types).

### Struct types

A _struct_ is a sequence of named elements, called fields, each of which has a name and a type. Field names may be specified explicitly (IdentifierList) or implicitly (EmbeddedField). Within a struct, non-[blank]() field names must be [unique]().

```go
// An empty struct.
struct {}

// A struct with 6 fields.
struct {
	x, y int
	u float32
	_ float32  // padding
	A *[]int
	F func()
}
```

A field declared with a type but no explicit field name is called an _embedded field_. An embedded field must be specified as a type name T or as a pointer to a non-interface type name *T, and T itself may not be a pointer type. The unqualified type name acts as the field name.

```go
// A struct with four embedded fields of types T1, *T2, P.T3 and *P.T4
struct {
	T1        // field name is T1
	*T2       // field name is T2
	P.T3      // field name is T3
	*P.T4     // field name is T4
	x, y int  // field names are x and y
}
```

The following declaration is illegal because field names must be unique in a struct type:

```go
struct {
	T     // conflicts with embedded field *T and *P.T
	*T    // conflicts with embedded field T and *P.T
	*P.T  // conflicts with embedded field T and *T
}
```

A field or [method]() f of an embedded field in a struct x is called _promoted_ if x.f is a legal [selector]() that denotes that field or method f.

Promoted fields act like ordinary fields of a struct except that they cannot be used as field names in [composite literals]() of the struct.

Given a struct type S and a [named type]() T, promoted methods are included in the method set of the struct as follows:

* If S contains an embedded field T, the [method sets]() of S and *S both include promoted methods with receiver T. The method set of *S also includes promoted methods with receiver *T.
* If S contains an embedded field *T, the method sets of S and *S both include promoted methods with receiver T or *T.

A field declaration may be followed by an optional string literal tag, which becomes an attribute for all the fields in the corresponding field declaration. An empty tag string is equivalent to an absent tag. The tags are made visible through a [reflection interface]() and take part in [type identity]() for structs but are otherwise ignored.

```go
struct {
	x, y float64 ""  // an empty tag string is like an absent tag
	name string  "any string is permitted as a tag"
	_    [4]byte "ceci n'est pas un champ de structure"
}

// A struct corresponding to a TimeStamp protocol buffer.
// The tag strings define the protocol buffer field numbers;
// they follow the convention outlined by the reflect package.
struct {
	microsec  uint64 `protobuf:"1"`
	serverIP6 uint64 `protobuf:"2"`
}
```

A struct type T may not contain a field of type T, or of a type containing T as a component, directly or indirectly, if those containing types are only array or struct types.

```go
// invalid struct types
type (
	T1 struct{ T1 }            // T1 contains a field of T1
	T2 struct{ f [10]T2 }      // T2 contains T2 as component of an array
	T3 struct{ T4 }            // T3 contains T3 as component of an array in struct T4
	T4 struct{ f [10]T3 }      // T4 contains T4 as component of struct T3 in an array
)

// valid struct types
type (
	T5 struct{ f *T5 }         // T5 contains T5 as component of a pointer
	T6 struct{ f func() T6 }   // T6 contains T6 as component of a function type
	T7 struct{ f [10][]T7 }    // T7 contains T7 as component of a slice in an array
)
```

### Pointer types

See [Pointer types](mini-spec.md#pointer-types).
