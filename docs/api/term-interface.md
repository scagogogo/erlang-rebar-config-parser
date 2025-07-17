# Term Interface

The `Term` interface is the foundation of the type system in the Erlang Rebar Config Parser. All Erlang data types implement this interface, providing consistent methods for string representation and comparison.

## Interface Definition

```go
type Term interface {
    String() string
    Compare(other Term) bool
}
```

## Methods

### String()

Returns a string representation of the term in Erlang syntax.

#### Returns

- `string`: Erlang syntax representation of the term

#### Examples

```go
atom := parser.Atom{Value: "debug_info", IsQuoted: false}
fmt.Println(atom.String()) // Output: debug_info

quotedAtom := parser.Atom{Value: "complex-name", IsQuoted: true}
fmt.Println(quotedAtom.String()) // Output: 'complex-name'

str := parser.String{Value: "hello world"}
fmt.Println(str.String()) // Output: "hello world"

integer := parser.Integer{Value: 42}
fmt.Println(integer.String()) // Output: 42

float := parser.Float{Value: 3.14}
fmt.Println(float.String()) // Output: 3.14
```

### Compare(other Term)

Compares this term with another term for equality.

#### Parameters

- `other` (Term): The term to compare with

#### Returns

- `bool`: `true` if the terms are equal, `false` otherwise

#### Comparison Rules

1. **Type matching**: Terms must be of the same type to be equal
2. **Value comparison**: The actual values must be identical
3. **Special cases**:
   - For `Atom`: Only `Value` is compared, `IsQuoted` is ignored
   - For nested structures (`Tuple`, `List`): All elements must be equal
   - For numbers: Exact value match required (no type coercion)

#### Examples

```go
// Atom comparison (ignores IsQuoted)
atom1 := parser.Atom{Value: "test", IsQuoted: false}
atom2 := parser.Atom{Value: "test", IsQuoted: true}
fmt.Println(atom1.Compare(atom2)) // Output: true

// String comparison
str1 := parser.String{Value: "hello"}
str2 := parser.String{Value: "hello"}
str3 := parser.String{Value: "world"}
fmt.Println(str1.Compare(str2)) // Output: true
fmt.Println(str1.Compare(str3)) // Output: false

// Number comparison
int1 := parser.Integer{Value: 42}
int2 := parser.Integer{Value: 42}
float1 := parser.Float{Value: 42.0}
fmt.Println(int1.Compare(int2))   // Output: true
fmt.Println(int1.Compare(float1)) // Output: false (different types)
```

## Working with Terms

### Type Assertions

Use Go's type assertion to work with specific term types:

```go
func processTerm(term parser.Term) {
    switch t := term.(type) {
    case parser.Atom:
        fmt.Printf("Atom: %s\n", t.Value)
        if t.IsQuoted {
            fmt.Println("  (quoted)")
        }
    case parser.String:
        fmt.Printf("String: %s\n", t.Value)
    case parser.Integer:
        fmt.Printf("Integer: %d\n", t.Value)
    case parser.Float:
        fmt.Printf("Float: %f\n", t.Value)
    case parser.Tuple:
        fmt.Printf("Tuple with %d elements:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    case parser.List:
        fmt.Printf("List with %d elements:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    default:
        fmt.Printf("Unknown term type: %T\n", t)
    }
}
```

### Safe Type Checking

```go
func isAtom(term parser.Term) bool {
    _, ok := term.(parser.Atom)
    return ok
}

func getAtomValue(term parser.Term) (string, bool) {
    if atom, ok := term.(parser.Atom); ok {
        return atom.Value, true
    }
    return "", false
}

func getStringValue(term parser.Term) (string, bool) {
    if str, ok := term.(parser.String); ok {
        return str.Value, true
    }
    return "", false
}
```

### Working with Collections

```go
func processCollection(term parser.Term) {
    switch t := term.(type) {
    case parser.Tuple:
        fmt.Printf("Processing tuple with %d elements\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("Element %d: %s\n", i, elem.String())
        }
    case parser.List:
        fmt.Printf("Processing list with %d elements\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("Element %d: %s\n", i, elem.String())
        }
    default:
        fmt.Println("Not a collection type")
    }
}
```

## Comparison Examples

### Simple Comparisons

```go
// Create terms
atom1 := parser.Atom{Value: "debug_info"}
atom2 := parser.Atom{Value: "debug_info"}
atom3 := parser.Atom{Value: "warnings_as_errors"}

// Compare
fmt.Println(atom1.Compare(atom2)) // true
fmt.Println(atom1.Compare(atom3)) // false

// Cross-type comparison
str := parser.String{Value: "debug_info"}
fmt.Println(atom1.Compare(str)) // false (different types)
```

### Complex Structure Comparisons

```go
// Create identical tuples
tuple1 := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "cowboy"},
        parser.String{Value: "2.9.0"},
    },
}

tuple2 := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "cowboy"},
        parser.String{Value: "2.9.0"},
    },
}

tuple3 := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "cowboy"},
        parser.String{Value: "2.8.0"}, // Different version
    },
}

fmt.Println(tuple1.Compare(tuple2)) // true
fmt.Println(tuple1.Compare(tuple3)) // false
```

### List Comparisons

```go
list1 := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "warnings_as_errors"},
    },
}

list2 := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "warnings_as_errors"},
    },
}

list3 := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        // Missing second element
    },
}

fmt.Println(list1.Compare(list2)) // true
fmt.Println(list1.Compare(list3)) // false (different lengths)
```

## Utility Functions

### Finding Terms in Collections

```go
func findAtomInList(list parser.List, atomValue string) bool {
    for _, elem := range list.Elements {
        if atom, ok := elem.(parser.Atom); ok && atom.Value == atomValue {
            return true
        }
    }
    return false
}

func findTupleByFirstAtom(list parser.List, atomValue string) (parser.Tuple, bool) {
    for _, elem := range list.Elements {
        if tuple, ok := elem.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == atomValue {
                return tuple, true
            }
        }
    }
    return parser.Tuple{}, false
}
```

### Term Validation

```go
func validateDependency(term parser.Term) error {
    tuple, ok := term.(parser.Tuple)
    if !ok {
        return fmt.Errorf("dependency must be a tuple")
    }
    
    if len(tuple.Elements) < 2 {
        return fmt.Errorf("dependency tuple must have at least 2 elements")
    }
    
    if _, ok := tuple.Elements[0].(parser.Atom); !ok {
        return fmt.Errorf("dependency name must be an atom")
    }
    
    switch tuple.Elements[1].(type) {
    case parser.String, parser.Tuple:
        // Valid version specification
        return nil
    default:
        return fmt.Errorf("dependency version must be a string or tuple")
    }
}
```
