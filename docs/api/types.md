# Types

This section documents all the data types used in the Erlang Rebar Config Parser library.

## RebarConfig

```go
type RebarConfig struct {
    Raw   string
    Terms []Term
}
```

Represents a parsed rebar.config file containing all configuration terms.

### Fields

- `Raw` (string): The original raw content of the configuration file
- `Terms` ([]Term): List of all top-level configuration terms

### Example

```go
config, _ := parser.Parse(`{erl_opts, [debug_info]}.`)
fmt.Printf("Raw content: %s\n", config.Raw)
fmt.Printf("Number of terms: %d\n", len(config.Terms))
```

---

## Term Interface

```go
type Term interface {
    String() string
    Compare(other Term) bool
}
```

The base interface implemented by all Erlang term types.

### Methods

- `String()`: Returns a string representation of the term
- `Compare(other Term)`: Compares this term with another term for equality

---

## Atom

```go
type Atom struct {
    Value    string
    IsQuoted bool
}
```

Represents an Erlang atom (symbol).

### Fields

- `Value` (string): The atom's value
- `IsQuoted` (bool): Whether the atom was quoted in the original syntax

### Methods

```go
func (a Atom) String() string
func (a Atom) Compare(other Term) bool
```

### Examples

```go
// Regular atom: debug_info
atom1 := parser.Atom{Value: "debug_info", IsQuoted: false}
fmt.Println(atom1.String()) // Output: debug_info

// Quoted atom: 'complex-name'
atom2 := parser.Atom{Value: "complex-name", IsQuoted: true}
fmt.Println(atom2.String()) // Output: 'complex-name'

// Comparison (ignores IsQuoted)
atom3 := parser.Atom{Value: "test", IsQuoted: false}
atom4 := parser.Atom{Value: "test", IsQuoted: true}
fmt.Println(atom3.Compare(atom4)) // Output: true
```

---

## String

```go
type String struct {
    Value string
}
```

Represents an Erlang string (double-quoted text).

### Fields

- `Value` (string): The string content

### Methods

```go
func (s String) String() string
func (s String) Compare(other Term) bool
```

### Examples

```go
str := parser.String{Value: "hello world"}
fmt.Println(str.String()) // Output: "hello world"

str1 := parser.String{Value: "test"}
str2 := parser.String{Value: "test"}
fmt.Println(str1.Compare(str2)) // Output: true
```

---

## Integer

```go
type Integer struct {
    Value int64
}
```

Represents an Erlang integer.

### Fields

- `Value` (int64): The integer value

### Methods

```go
func (i Integer) String() string
func (i Integer) Compare(other Term) bool
```

### Examples

```go
num := parser.Integer{Value: 42}
fmt.Println(num.String()) // Output: 42

num1 := parser.Integer{Value: 123}
num2 := parser.Integer{Value: 123}
fmt.Println(num1.Compare(num2)) // Output: true
```

---

## Float

```go
type Float struct {
    Value float64
}
```

Represents an Erlang floating-point number.

### Fields

- `Value` (float64): The float value

### Methods

```go
func (f Float) String() string
func (f Float) Compare(other Term) bool
```

### Examples

```go
num := parser.Float{Value: 3.14}
fmt.Println(num.String()) // Output: 3.14

// Scientific notation
num2 := parser.Float{Value: 1.5e-3}
fmt.Println(num2.String()) // Output: 0.0015
```

---

## Tuple

```go
type Tuple struct {
    Elements []Term
}
```

Represents an Erlang tuple `{elem1, elem2, ...}`.

### Fields

- `Elements` ([]Term): List of elements in the tuple

### Methods

```go
func (t Tuple) String() string
func (t Tuple) Compare(other Term) bool
```

### Examples

```go
// Simple tuple: {key, value}
tuple := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "key"},
        parser.String{Value: "value"},
    },
}
fmt.Println(tuple.String()) // Output: {key, "value"}

// Nested tuple: {deps, [{cowboy, "2.9.0"}]}
nestedTuple := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "deps"},
        parser.List{
            Elements: []parser.Term{
                parser.Tuple{
                    Elements: []parser.Term{
                        parser.Atom{Value: "cowboy"},
                        parser.String{Value: "2.9.0"},
                    },
                },
            },
        },
    },
}
```

---

## List

```go
type List struct {
    Elements []Term
}
```

Represents an Erlang list `[elem1, elem2, ...]`.

### Fields

- `Elements` ([]Term): List of elements in the list

### Methods

```go
func (l List) String() string
func (l List) Compare(other Term) bool
```

### Examples

```go
// Simple list: [debug_info, warnings_as_errors]
list := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "warnings_as_errors"},
    },
}
fmt.Println(list.String()) // Output: [debug_info, warnings_as_errors]

// Mixed type list: [atom, "string", 123]
mixedList := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "atom"},
        parser.String{Value: "string"},
        parser.Integer{Value: 123},
    },
}
```

---

## Type Checking and Conversion

### Safe Type Assertions

```go
func processTerms(terms []parser.Term) {
    for _, term := range terms {
        switch t := term.(type) {
        case parser.Atom:
            fmt.Printf("Atom: %s (quoted: %t)\n", t.Value, t.IsQuoted)
        case parser.String:
            fmt.Printf("String: %s\n", t.Value)
        case parser.Integer:
            fmt.Printf("Integer: %d\n", t.Value)
        case parser.Float:
            fmt.Printf("Float: %f\n", t.Value)
        case parser.Tuple:
            fmt.Printf("Tuple with %d elements\n", len(t.Elements))
        case parser.List:
            fmt.Printf("List with %d elements\n", len(t.Elements))
        default:
            fmt.Printf("Unknown term type: %T\n", t)
        }
    }
}
```

### Working with Nested Structures

```go
func extractDependencies(config *parser.RebarConfig) []string {
    var deps []string
    
    if depsTerms, ok := config.GetDeps(); ok && len(depsTerms) > 0 {
        if depsList, ok := depsTerms[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                        deps = append(deps, atom.Value)
                    }
                }
            }
        }
    }
    
    return deps
}
```
