# API Reference

The Erlang Rebar Config Parser provides a comprehensive API for parsing and manipulating Erlang rebar configuration files. This section documents all public functions, types, and methods available in the library.

## Package Overview

The main package `github.com/scagogogo/erlang-rebar-config-parser/pkg/parser` provides functionality to:

- **Parse rebar.config files** from various sources (files, strings, io.Reader)
- **Access configuration elements** through convenient helper methods
- **Format and pretty-print** configurations with customizable indentation
- **Compare Erlang terms** for equality with type-aware comparison
- **Handle all common Erlang data types** (atoms, strings, numbers, tuples, lists)
- **Process escape sequences** in strings and quoted atoms
- **Validate configuration structure** and provide detailed error messages

## Quick Navigation

- **[Core Functions](./core-functions)** - Main parsing functions (`ParseFile`, `Parse`, `ParseReader`)
- **[Types](./types)** - Data type definitions (`RebarConfig`, `Term`, `Atom`, `String`, etc.)
- **[RebarConfig Methods](./rebar-config)** - Configuration access methods (`GetDeps`, `GetErlOpts`, `Format`, etc.)
- **[Term Interface](./term-interface)** - Term types and operations (`String()`, `Compare()`)

## Complete API Overview

### Core Parsing Functions

| Function | Description | Input | Output |
|----------|-------------|-------|--------|
| `ParseFile(path string)` | Parse rebar.config from file | File path | `*RebarConfig`, `error` |
| `Parse(input string)` | Parse rebar.config from string | Config string | `*RebarConfig`, `error` |
| `ParseReader(r io.Reader)` | Parse rebar.config from reader | io.Reader | `*RebarConfig`, `error` |
| `NewParser(input string)` | Create new parser instance | Input string | `*Parser` |

### RebarConfig Methods

| Method | Description | Returns |
|--------|-------------|---------|
| `GetTerm(name string)` | Get specific term by name | `Term`, `bool` |
| `GetTupleElements(name string)` | Get tuple elements (excluding name) | `[]Term`, `bool` |
| `GetDeps()` | Get dependencies configuration | `[]Term`, `bool` |
| `GetErlOpts()` | Get Erlang compilation options | `[]Term`, `bool` |
| `GetAppName()` | Get application name | `string`, `bool` |
| `GetPlugins()` | Get plugins configuration | `[]Term`, `bool` |
| `GetRelxConfig()` | Get relx (release) configuration | `[]Term`, `bool` |
| `GetProfilesConfig()` | Get profiles configuration | `[]Term`, `bool` |
| `Format(indent int)` | Format config with indentation | `string` |

### Term Types

| Type | Description | Fields/Methods |
|------|-------------|----------------|
| `Term` | Base interface for all Erlang terms | `String()`, `Compare(Term)` |
| `Atom` | Erlang atom (symbol) | `Value string`, `IsQuoted bool` |
| `String` | Erlang string (double-quoted) | `Value string` |
| `Integer` | Erlang integer | `Value int64` |
| `Float` | Erlang float | `Value float64` |
| `Tuple` | Erlang tuple `{a, b, c}` | `Elements []Term` |
| `List` | Erlang list `[a, b, c]` | `Elements []Term` |

### Utility Functions

| Function | Description | Purpose |
|----------|-------------|---------|
| `processEscapes(s string)` | Process escape sequences | Handle `\"`, `\\`, `\n`, `\r`, `\t` |
| `isDigit(ch byte)` | Check if character is digit | Character classification |
| `isAtomStart(ch byte)` | Check if valid atom start | Atom validation |
| `isAtomChar(ch byte)` | Check if valid atom character | Atom validation |
| `isSimpleTerm(term Term)` | Check if term is simple | Formatting decisions |
| `allSimpleTerms(terms []Term)` | Check if all terms are simple | Formatting decisions |

## Basic Usage Pattern

```go
import "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"

// 1. Parse configuration from different sources
config, err := parser.ParseFile("rebar.config")
if err != nil {
    log.Fatalf("Parse error: %v", err)
}

// Alternative: parse from string
configStr := `{erl_opts, [debug_info]}.`
config, err = parser.Parse(configStr)

// Alternative: parse from reader
file, _ := os.Open("rebar.config")
defer file.Close()
config, err = parser.ParseReader(file)

// 2. Access configuration elements
if deps, ok := config.GetDeps(); ok {
    fmt.Println("Dependencies found!")
}

if appName, ok := config.GetAppName(); ok {
    fmt.Printf("Application: %s\n", appName)
}

// 3. Work with terms directly
for _, term := range config.Terms {
    switch t := term.(type) {
    case parser.Tuple:
        fmt.Printf("Tuple with %d elements\n", len(t.Elements))
    case parser.List:
        fmt.Printf("List with %d elements\n", len(t.Elements))
    }
}

// 4. Format and display
formatted := config.Format(2) // 2-space indentation
fmt.Println(formatted)
```

## Error Handling

The parser provides detailed error messages with position information:

```go
config, err := parser.Parse(`{invalid syntax`)
if err != nil {
    // Error will include line and column information
    fmt.Printf("Parse error: %v\n", err)
    // Example: "syntax error at line 1, column 15: expected '}'"
}
```

## Advanced Usage

### Custom Term Processing

```go
func processTerm(term parser.Term) {
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
        fmt.Printf("Tuple with %d elements:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    case parser.List:
        fmt.Printf("List with %d elements:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    }
}
```

### Term Comparison

```go
atom1 := parser.Atom{Value: "debug_info", IsQuoted: false}
atom2 := parser.Atom{Value: "debug_info", IsQuoted: true}

// Compare ignores IsQuoted flag
if atom1.Compare(atom2) {
    fmt.Println("Atoms are equal")
}

// Compare different types
str := parser.String{Value: "debug_info"}
if !atom1.Compare(str) {
    fmt.Println("Different types are not equal")
}
```

## Performance Considerations

- **Memory Usage**: Large configurations are parsed into memory. For very large files, consider streaming approaches.
- **Parsing Speed**: The parser is optimized for typical rebar.config files. Complex nested structures may take longer.
- **Formatting**: The `Format()` method creates a new string. For large configurations, this may use significant memory.

## Thread Safety

The parser types are **not thread-safe**. If you need to access parsed configurations from multiple goroutines, use appropriate synchronization mechanisms or create separate parser instances.

## Error Handling

All parsing functions return an error as the second return value. Common error scenarios include:

- **File not found** - When parsing from a file that doesn't exist
- **Syntax errors** - When the Erlang syntax is invalid
- **Unexpected characters** - When encountering unsupported syntax
- **Unterminated strings/atoms** - When quotes are not properly closed
- **Invalid numbers** - When number format is incorrect

Example error handling:

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    if strings.Contains(err.Error(), "no such file") {
        log.Fatal("Configuration file not found")
    } else if strings.Contains(err.Error(), "syntax error") {
        log.Fatalf("Invalid configuration syntax: %v", err)
    } else {
        log.Fatalf("Failed to parse configuration: %v", err)
    }
}
```

## Type Safety

The library uses Go's type system to provide safe access to Erlang terms. Use type assertions to work with specific term types:

```go
// Safe type assertion
if atom, ok := term.(parser.Atom); ok {
    fmt.Println("Atom value:", atom.Value)
}

// Working with nested structures
if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
    if list, ok := tuple.Elements[1].(parser.List); ok {
        fmt.Printf("List has %d elements\n", len(list.Elements))
    }
}
```

## Performance Considerations

- **Memory usage**: The parser loads the entire configuration into memory
- **Parsing speed**: Linear time complexity relative to input size
- **Type assertions**: Minimal overhead for type checking
- **String formatting**: Lazy evaluation - only computed when needed

## Thread Safety

The library is **read-safe** but not **write-safe**:

- Multiple goroutines can safely read from the same `RebarConfig` instance
- Parsing operations are independent and can be performed concurrently
- Do not modify `RebarConfig` or `Term` instances from multiple goroutines

## Compatibility

- **Go version**: Requires Go 1.18 or later
- **Erlang compatibility**: Supports standard Erlang term syntax
- **Rebar versions**: Compatible with rebar3 and rebar2 configuration formats
