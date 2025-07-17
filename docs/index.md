---
layout: home

hero:
  name: "Erlang Rebar Config Parser"
  text: "A Go library for parsing Erlang rebar configuration files"
  tagline: "Parse, access, and format Erlang rebar.config files with ease"
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/erlang-rebar-config-parser

features:
  - icon: 🚀
    title: Easy to Use
    details: Simple API for parsing rebar.config files into structured Go objects with comprehensive helper methods.
  
  - icon: 🔧
    title: Full Feature Support
    details: Support for all common Erlang term types including tuples, lists, atoms, strings, numbers, and nested structures.
  
  - icon: 📝
    title: Pretty Printing
    details: Format and pretty-print rebar configuration files with configurable indentation for better readability.
  
  - icon: ⚡
    title: High Performance
    details: Efficient parsing with 98% test coverage and comprehensive error handling for production use.
  
  - icon: 🌍
    title: Multilingual
    details: Complete documentation and examples available in both English and Chinese.
  
  - icon: 🔍
    title: Term Comparison
    details: Built-in comparison functionality to check equality between different Erlang terms and configurations.
---

## Quick Example

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

## Installation

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

## Features

- **Parse rebar.config files** into structured Go objects
- **Support for all common Erlang term types** (tuples, lists, atoms, strings, numbers)
- **Helper methods** to easily access common configuration elements
- **Full support for nested data structures**
- **Handle comments and whitespace** correctly
- **Pretty-printing** with configurable indentation
- **Compare functionality** to check term equality
- **Continuous Integration** via GitHub Actions
- **Comprehensive documentation** with examples in English and Chinese
- **98% test coverage** with comprehensive edge case testing

## Supported Erlang Term Types

| Erlang Type | Example | Go Representation |
|-------------|---------|-------------------|
| Atoms | `atom_name`, `'quoted-atom'` | `Atom{Value: "atom_name", IsQuoted: false}` |
| Strings | `"hello world"` | `String{Value: "hello world"}` |
| Integers | `123`, `-42` | `Integer{Value: 123}` |
| Floats | `3.14`, `-1.5e-3` | `Float{Value: 3.14}` |
| Tuples | `{key, value}` | `Tuple{Elements: []Term{...}}` |
| Lists | `[1, 2, 3]` | `List{Elements: []Term{...}}` |
