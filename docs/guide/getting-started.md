# Getting Started

Welcome to the Erlang Rebar Config Parser! This guide will help you get up and running quickly with parsing and manipulating Erlang rebar configuration files in Go.

## What is this library?

The Erlang Rebar Config Parser is a Go library that allows you to:

- Parse Erlang rebar.config files into structured Go objects
- Access configuration elements through convenient helper methods
- Format and pretty-print configurations
- Compare Erlang terms for equality
- Handle all common Erlang data types (atoms, strings, numbers, tuples, lists)

## Prerequisites

- Go 1.18 or later
- Basic understanding of Erlang syntax (helpful but not required)
- Familiarity with Go programming

## Installation

Add the library to your Go project:

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

## Your First Program

Let's start with a simple example that parses a rebar.config file:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Sample rebar.config content
    configContent := `
    {erl_opts, [debug_info, warnings_as_errors]}.
    {deps, [
        {cowboy, "2.9.0"},
        {jsx, "3.1.0"}
    ]}.
    `
    
    // Parse the configuration
    config, err := parser.Parse(configContent)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Printf("Successfully parsed %d configuration terms\n", len(config.Terms))
    
    // Access dependencies
    deps, ok := config.GetDeps()
    if ok {
        fmt.Println("Dependencies found!")
    }
    
    // Pretty print the configuration
    fmt.Println("\nFormatted configuration:")
    fmt.Println(config.Format(2))
}
```

## Understanding the Output

When you run this program, you'll see:

```
Successfully parsed 2 configuration terms
Dependencies found!

Formatted configuration:
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
  {cowboy, "2.9.0"},
  {jsx, "3.1.0"}
]}.
```

## Key Concepts

### 1. Terms

Everything in Erlang is a "term". The library represents these as Go types that implement the `Term` interface:

- **Atoms**: `debug_info`, `'quoted-atom'`
- **Strings**: `"hello world"`
- **Numbers**: `123`, `3.14`
- **Tuples**: `{key, value}`
- **Lists**: `[item1, item2]`

### 2. Configuration Structure

A rebar.config file consists of top-level terms, typically tuples where the first element is an atom identifying the configuration section:

```erlang
{erl_opts, [debug_info]}.        % Erlang compilation options
{deps, [{cowboy, "2.9.0"}]}.     % Dependencies
{profiles, [{test, [...]}]}.     % Build profiles
```

### 3. Helper Methods

The library provides convenient methods to access common configuration sections:

```go
// Instead of manually parsing tuples
term, ok := config.GetTerm("deps")

// Use helper methods
deps, ok := config.GetDeps()
erlOpts, ok := config.GetErlOpts()
appName, ok := config.GetAppName()
```

## Common Patterns

### Parsing from Different Sources

```go
// From file
config, err := parser.ParseFile("rebar.config")

// From string
config, err := parser.Parse(configString)

// From any io.Reader
config, err := parser.ParseReader(reader)
```

### Safe Type Checking

```go
// Always use type assertions safely
if atom, ok := term.(parser.Atom); ok {
    fmt.Printf("Atom value: %s\n", atom.Value)
}

// Check for specific structures
if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
    // Safe to access tuple.Elements[0] and tuple.Elements[1]
}
```

### Error Handling

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    // Handle different error types
    if strings.Contains(err.Error(), "no such file") {
        log.Fatal("Configuration file not found")
    } else if strings.Contains(err.Error(), "syntax error") {
        log.Fatalf("Invalid syntax: %v", err)
    } else {
        log.Fatalf("Parse error: %v", err)
    }
}
```

## Next Steps

Now that you understand the basics, explore these topics:

1. **[Installation](./installation)** - Detailed installation instructions
2. **[Basic Usage](./basic-usage)** - Common usage patterns and examples
3. **[Advanced Usage](./advanced-usage)** - Complex scenarios and best practices
4. **[API Reference](../api/)** - Complete API documentation
5. **[Examples](../examples/)** - Real-world examples and use cases

## Quick Reference

### Essential Functions

```go
// Parsing
config, err := parser.ParseFile("rebar.config")
config, err := parser.Parse(configString)
config, err := parser.ParseReader(reader)

// Accessing configuration
deps, ok := config.GetDeps()
erlOpts, ok := config.GetErlOpts()
appName, ok := config.GetAppName()

// Formatting
formatted := config.Format(2) // 2-space indentation
```

### Essential Types

```go
// Check term types
switch t := term.(type) {
case parser.Atom:
    // t.Value, t.IsQuoted
case parser.String:
    // t.Value
case parser.Integer:
    // t.Value
case parser.Tuple:
    // t.Elements
case parser.List:
    // t.Elements
}
```

## Getting Help

- **Documentation**: Browse the complete [API Reference](../api/)
- **Examples**: Check out [practical examples](../examples/)
- **Issues**: Report bugs on [GitHub](https://github.com/scagogogo/erlang-rebar-config-parser/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/scagogogo/erlang-rebar-config-parser/discussions)
