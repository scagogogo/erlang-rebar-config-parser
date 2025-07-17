# Erlang Rebar Config Parser

[中文文档](README_zh.md) | [📖 Documentation](https://scagogogo.github.io/erlang-rebar-config-parser/)

[![Go Tests and Examples](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/go.yml/badge.svg)](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/go.yml)
[![Documentation](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/docs.yml/badge.svg)](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/docs.yml)
[![GoDoc](https://godoc.org/github.com/scagogogo/erlang-rebar-config-parser?status.svg)](https://godoc.org/github.com/scagogogo/erlang-rebar-config-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/erlang-rebar-config-parser)](https://goreportcard.com/report/github.com/scagogogo/erlang-rebar-config-parser)
[![Go Version](https://img.shields.io/github/go-mod/go-version/scagogogo/erlang-rebar-config-parser)](https://github.com/scagogogo/erlang-rebar-config-parser/blob/main/go.mod)
[![License](https://img.shields.io/github/license/scagogogo/erlang-rebar-config-parser)](https://github.com/scagogogo/erlang-rebar-config-parser/blob/main/LICENSE)

A Go library for parsing Erlang rebar configuration files. This library allows you to parse `rebar.config` files into structured Go objects, making it easy to programmatically access and manipulate Erlang project configurations.

## 📚 Documentation

- **[Complete Documentation](https://scagogogo.github.io/erlang-rebar-config-parser/)** - Full documentation website
- **[Getting Started Guide](https://scagogogo.github.io/erlang-rebar-config-parser/guide/getting-started)** - Quick start tutorial
- **[API Reference](https://scagogogo.github.io/erlang-rebar-config-parser/api/)** - Complete API documentation
- **[Examples](https://scagogogo.github.io/erlang-rebar-config-parser/examples/)** - Real-world examples

## 🌟 Features

- Parse rebar.config files into structured Go objects
- Support for all common Erlang term types (tuples, lists, atoms, strings, numbers)
- Helper methods to easily access common configuration elements
- Full support for nested data structures
- Handle comments and whitespace correctly
- Pretty-printing with configurable indentation
- Compare functionality to check term equality
- Continuous Integration via GitHub Actions
- Comprehensive documentation with examples in English and Chinese
- 98% test coverage with comprehensive edge case testing

## 📦 Installation

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Parse a rebar.config file
    config, err := parser.ParseFile("path/to/rebar.config")
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Get and print dependencies
    deps, ok := config.GetDeps()
    if ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("Found %d dependencies\n", len(depsList.Elements))
            
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                        fmt.Printf("- Dependency: %s\n", atom.Value)
                    }
                }
            }
        }
    }
    
    // Format and print the config with nice indentation
    fmt.Println("\nFormatted config:")
    fmt.Println(config.Format(2))
}
```

## 📚 Detailed Usage Guide

### Parsing a rebar.config file

There are three ways to parse a rebar.config file:

```go
// 1. Parse from file
config, err := parser.ParseFile("path/to/rebar.config")

// 2. Parse from string
configStr := `
{erl_opts, [debug_info]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
`
config, err = parser.Parse(configStr)

// 3. Parse from io.Reader
file, err := os.Open("path/to/rebar.config")
if err != nil {
    log.Fatalf("Failed to open file: %v", err)
}
defer file.Close()
config, err = parser.ParseReader(file)
```

### Accessing Configuration Elements

The library provides helper methods to access common configuration elements:

```go
// Get application name
if appName, ok := config.GetAppName(); ok {
    fmt.Printf("Application name: %s\n", appName)
}

// Get dependencies
if deps, ok := config.GetDeps(); ok {
    if list, ok := deps[0].(parser.List); ok {
        fmt.Printf("Found %d dependencies\n", len(list.Elements))
    }
}

// Get Erlang compilation options
if erlOpts, ok := config.GetErlOpts(); ok {
    fmt.Printf("Erlang options: %v\n", erlOpts)
}

// Get plugins
if plugins, ok := config.GetPlugins(); ok {
    fmt.Printf("Plugins: %v\n", plugins)
}

// Get relx configuration
if relx, ok := config.GetRelxConfig(); ok {
    fmt.Printf("Relx config: %v\n", relx)
}

// Get profiles configuration
if profiles, ok := config.GetProfilesConfig(); ok {
    fmt.Printf("Profiles: %v\n", profiles)
}

// Get any arbitrary term by name
if term, ok := config.GetTerm("minimum_otp_vsn"); ok {
    fmt.Printf("Minimum OTP version: %s\n", term)
}

// Get the elements of a named tuple
if elements, ok := config.GetTupleElements("shell"); ok {
    fmt.Printf("Shell configuration elements: %v\n", elements)
}
```

### Working with Different Term Types

The library represents Erlang terms as Go structs that implement the `Term` interface:

```go
// Working with Atoms
if atom, ok := term.(parser.Atom); ok {
    fmt.Printf("Atom value: %s, Quoted: %t\n", atom.Value, atom.IsQuoted)
}

// Working with Strings
if str, ok := term.(parser.String); ok {
    fmt.Printf("String value: %s\n", str.Value)
}

// Working with Integers
if integer, ok := term.(parser.Integer); ok {
    fmt.Printf("Integer value: %d\n", integer.Value)
}

// Working with Floats
if float, ok := term.(parser.Float); ok {
    fmt.Printf("Float value: %f\n", float.Value)
}

// Working with Tuples
if tuple, ok := term.(parser.Tuple); ok {
    fmt.Printf("Tuple with %d elements\n", len(tuple.Elements))
    for i, elem := range tuple.Elements {
        fmt.Printf("  Element %d: %s\n", i, elem.String())
    }
}

// Working with Lists
if list, ok := term.(parser.List); ok {
    fmt.Printf("List with %d elements\n", len(list.Elements))
    for i, elem := range list.Elements {
        fmt.Printf("  Element %d: %s\n", i, elem.String())
    }
}
```

### Pretty Printing

The library includes a `Format` method for pretty-printing rebar configuration files:

```go
// Format with 2-space indentation
fmt.Println(config.Format(2))

// Format with 4-space indentation
fmt.Println(config.Format(4))
```

Example output with 2-space indentation:

```erlang
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
  {cowboy, "2.9.0"},
  {jsx, "3.1.0"}
]}.

{relx, [
  {release, {my_app, "0.1.0"}, [my_app, sasl]},
  {dev_mode, true},
  {include_erts, false}
]}.
```

### Comparing Terms

The library provides a `Compare` method to check if two terms are equal:

```go
// Compare atoms
atom1 := parser.Atom{Value: "test", IsQuoted: false}
atom2 := parser.Atom{Value: "test", IsQuoted: true} // Quoting doesn't affect equality
if atom1.Compare(atom2) {
    fmt.Println("Atoms are equal") // This will print
}

// Compare complex structures
tuple1 := parser.Tuple{Elements: []parser.Term{
    parser.Atom{Value: "cowboy"},
    parser.String{Value: "2.9.0"},
}}
tuple2 := parser.Tuple{Elements: []parser.Term{
    parser.Atom{Value: "cowboy"},
    parser.String{Value: "2.9.0"},
}}
if tuple1.Compare(tuple2) {
    fmt.Println("Tuples are equal") // This will print
}
```

## 🔍 Real-World Examples

### Example 1: Analyzing Dependencies in a Project

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    config, err := parser.ParseFile("./rebar.config")
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    deps, ok := config.GetDeps()
    if !ok {
        fmt.Println("No dependencies found")
        return
    }
    
    if depsList, ok := deps[0].(parser.List); ok {
        fmt.Printf("Project has %d dependencies:\n", len(depsList.Elements))
        
        for _, dep := range depsList.Elements {
            if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("- %s: ", atom.Value)
                    
                    // Check dependency version specification
                    switch v := tuple.Elements[1].(type) {
                    case parser.String:
                        fmt.Printf("Version %s\n", v.Value)
                    case parser.Tuple:
                        if len(v.Elements) >= 2 {
                            if sourceAtom, ok := v.Elements[0].(parser.Atom); ok {
                                switch sourceAtom.Value {
                                case "git":
                                    fmt.Printf("Git repo")
                                    if url, ok := v.Elements[1].(parser.String); ok {
                                        fmt.Printf(" at %s", url.Value)
                                    }
                                    
                                    // Check for tag, branch, or commit reference
                                    if len(v.Elements) >= 3 {
                                        if refTuple, ok := v.Elements[2].(parser.Tuple); ok && len(refTuple.Elements) >= 2 {
                                            if refType, ok := refTuple.Elements[0].(parser.Atom); ok {
                                                if refValue, ok := refTuple.Elements[1].(parser.String); ok {
                                                    fmt.Printf(" (%s: %s)", refType.Value, refValue.Value)
                                                }
                                            }
                                        }
                                    }
                                    fmt.Println()
                                case "hex":
                                    fmt.Printf("Hex package")
                                    if ver, ok := v.Elements[1].(parser.String); ok {
                                        fmt.Printf(" version %s", ver.Value)
                                    }
                                    fmt.Println()
                                default:
                                    fmt.Printf("Custom source: %s\n", sourceAtom.Value)
                                }
                            }
                        }
                    default:
                        fmt.Printf("Unknown version format: %T\n", tuple.Elements[1])
                    }
                }
            }
        }
    }
}
```

### Example 2: Modifying and Writing a Config File

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Parse existing config
    config, err := parser.ParseFile("./rebar.config")
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Get deps
    depsTerms, ok := config.GetDeps()
    if !ok || len(depsTerms) == 0 {
        log.Fatal("No deps found in config")
    }
    
    depsList, ok := depsTerms[0].(parser.List)
    if !ok {
        log.Fatal("Deps is not a list")
    }
    
    // Add a new dependency
    newDep := parser.Tuple{Elements: []parser.Term{
        parser.Atom{Value: "new_dependency"},
        parser.String{Value: "1.0.0"},
    }}
    
    // Create a new deps list with the added dependency
    updatedDepsList := parser.List{
        Elements: append(depsList.Elements, newDep),
    }
    
    // Find and update the deps term in the config
    for i, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == "deps" {
                // Create a new tuple with the updated deps list
                config.Terms[i] = parser.Tuple{
                    Elements: []parser.Term{
                        parser.Atom{Value: "deps"},
                        updatedDepsList,
                    },
                }
                break
            }
        }
    }
    
    // Format and write the updated config
    formatted := config.Format(4)
    err = os.WriteFile("./rebar.config.new", []byte(formatted), 0644)
    if err != nil {
        log.Fatalf("Failed to write updated config: %v", err)
    }
    
    fmt.Println("Updated config written to ./rebar.config.new")
}
```

## 📖 API Documentation

### Core Functions

| Function | Description | Example |
|----------|-------------|---------|
| `ParseFile(path string) (*RebarConfig, error)` | Parses a rebar.config file from the given file path | `config, err := parser.ParseFile("./rebar.config")` |
| `Parse(input string) (*RebarConfig, error)` | Parses a rebar.config file from the given string | `config, err := parser.Parse(configStr)` |
| `ParseReader(r io.Reader) (*RebarConfig, error)` | Parses a rebar.config file from the given reader | `config, err := parser.ParseReader(file)` |

### RebarConfig Methods

| Method | Description | Example |
|--------|-------------|---------|
| `Format(indent int) string` | Returns a formatted string representation of the config | `formatted := config.Format(2)` |
| `GetTerm(name string) (Term, bool)` | Retrieves a specific term from the config by name | `term, ok := config.GetTerm("deps")` |
| `GetTupleElements(name string) ([]Term, bool)` | Gets the elements of a named tuple | `elements, ok := config.GetTupleElements("deps")` |
| `GetDeps() ([]Term, bool)` | Retrieves the deps configuration | `deps, ok := config.GetDeps()` |
| `GetErlOpts() ([]Term, bool)` | Retrieves the erl_opts configuration | `opts, ok := config.GetErlOpts()` |
| `GetAppName() (string, bool)` | Retrieves the application name | `name, ok := config.GetAppName()` |
| `GetPlugins() ([]Term, bool)` | Retrieves the plugins configuration | `plugins, ok := config.GetPlugins()` |
| `GetRelxConfig() ([]Term, bool)` | Retrieves the relx configuration | `relx, ok := config.GetRelxConfig()` |
| `GetProfilesConfig() ([]Term, bool)` | Retrieves the profiles configuration | `profiles, ok := config.GetProfilesConfig()` |

### Term Interface

The `Term` interface is implemented by all Erlang term types:

```go
type Term interface {
    String() string           // Returns a string representation
    Compare(other Term) bool  // Compares this term with another term
}
```

### Erlang Term Types

| Type | Description | Fields |
|------|-------------|--------|
| `Atom` | Represents an Erlang atom | `Value string`, `IsQuoted bool` |
| `String` | Represents an Erlang string (double-quoted) | `Value string` |
| `Integer` | Represents an Erlang integer | `Value int64` |
| `Float` | Represents an Erlang float | `Value float64` |
| `Tuple` | Represents an Erlang tuple | `Elements []Term` |
| `List` | Represents an Erlang list | `Elements []Term` |

## 📋 Supported Erlang Term Types

The library supports the following Erlang term types:

| Erlang Type | Example | Go Representation |
|-------------|---------|-------------------|
| Atoms | `atom_name`, `'quoted-atom'` | `Atom{Value: "atom_name", IsQuoted: false}`, `Atom{Value: "quoted-atom", IsQuoted: true}` |
| Strings | `"hello world"` | `String{Value: "hello world"}` |
| Integers | `123`, `-42` | `Integer{Value: 123}`, `Integer{Value: -42}` |
| Floats | `3.14`, `-1.5e-3` | `Float{Value: 3.14}`, `Float{Value: -0.0015}` |
| Tuples | `{key, value}` | `Tuple{Elements: []Term{Atom{Value: "key"}, Atom{Value: "value"}}}` |
| Lists | `[1, 2, 3]` | `List{Elements: []Term{Integer{Value: 1}, Integer{Value: 2}, Integer{Value: 3}}}` |

## 🔄 Continuous Integration

This project uses GitHub Actions for continuous integration. On every push and pull request to the main branch, the following checks are automatically performed:

- Tests are run on multiple Go versions (1.18, 1.19, 1.20, 1.21)
- All example code is executed to ensure it works as expected
- Code coverage is calculated and reported to Codecov
- Code linting is performed using golangci-lint

The workflow configuration can be found in `.github/workflows/go.yml`.

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes
4. Add tests for your changes
5. Run tests: `go test -v ./...`
6. Commit your changes: `git commit -am 'Add some feature'`
7. Push to the branch: `git push origin feature/your-feature-name`
8. Submit a pull request

Please ensure your code follows the existing style and includes appropriate tests.

### Development Guidelines

- All public functions should have proper documentation comments
- Include examples in documentation where appropriate
- Write tests for all new functionality
- Maintain backward compatibility when possible

## 📜 License

This project is licensed under the terms found in the [LICENSE](LICENSE) file. 