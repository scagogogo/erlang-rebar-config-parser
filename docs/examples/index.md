# Examples

This section provides practical examples of using the Erlang Rebar Config Parser in real-world scenarios. Each example includes complete, runnable code with explanations.

## Overview

The examples are organized by complexity and use case:

- **[Basic Parsing](./basic-parsing)** - Simple parsing and term access
- **[Configuration Access](./config-access)** - Using helper methods to access common configurations
- **[Pretty Printing](./pretty-printing)** - Formatting and displaying configurations
- **[Term Comparison](./comparison)** - Comparing configurations and terms
- **[Complex Analysis](./complex-analysis)** - Advanced parsing and analysis scenarios

## Quick Examples

### Parse and Display

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    config, err := parser.ParseFile("rebar.config")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Configuration has %d terms\n", len(config.Terms))
    fmt.Println(config.Format(2))
}
```

### Extract Dependencies

```go
func extractDependencies(configPath string) ([]string, error) {
    config, err := parser.ParseFile(configPath)
    if err != nil {
        return nil, err
    }
    
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
    
    return deps, nil
}
```

### Configuration Validation

```go
func validateConfig(config *parser.RebarConfig) []string {
    var warnings []string
    
    // Check for required sections
    if _, ok := config.GetDeps(); !ok {
        warnings = append(warnings, "No dependencies defined")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        warnings = append(warnings, "No Erlang options defined")
    }
    
    if _, ok := config.GetAppName(); !ok {
        warnings = append(warnings, "No application name defined")
    }
    
    return warnings
}
```

## Common Patterns

### Error Handling

```go
func parseWithErrorHandling(path string) (*parser.RebarConfig, error) {
    config, err := parser.ParseFile(path)
    if err != nil {
        // Provide context for different error types
        if strings.Contains(err.Error(), "no such file") {
            return nil, fmt.Errorf("configuration file '%s' not found", path)
        } else if strings.Contains(err.Error(), "syntax error") {
            return nil, fmt.Errorf("invalid syntax in '%s': %w", path, err)
        }
        return nil, fmt.Errorf("failed to parse '%s': %w", path, err)
    }
    return config, nil
}
```

### Safe Type Checking

```go
func safelyAccessTuple(term parser.Term, minElements int) (parser.Tuple, bool) {
    if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= minElements {
        return tuple, true
    }
    return parser.Tuple{}, false
}

func getAtomValue(term parser.Term) (string, bool) {
    if atom, ok := term.(parser.Atom); ok {
        return atom.Value, true
    }
    return "", false
}
```

### Working with Collections

```go
func processErlangList(list parser.List, processor func(parser.Term)) {
    for _, element := range list.Elements {
        processor(element)
    }
}

func findInList(list parser.List, predicate func(parser.Term) bool) (parser.Term, bool) {
    for _, element := range list.Elements {
        if predicate(element) {
            return element, true
        }
    }
    return nil, false
}
```

## Real-World Use Cases

### 1. Dependency Analyzer

Analyze project dependencies and their versions:

```go
type Dependency struct {
    Name    string
    Version string
    Source  string
}

func analyzeDependencies(configPath string) ([]Dependency, error) {
    config, err := parser.ParseFile(configPath)
    if err != nil {
        return nil, err
    }
    
    var deps []Dependency
    if depsTerms, ok := config.GetDeps(); ok && len(depsTerms) > 0 {
        if depsList, ok := depsTerms[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if d := parseDependency(dep); d != nil {
                    deps = append(deps, *d)
                }
            }
        }
    }
    
    return deps, nil
}
```

### 2. Configuration Merger

Merge multiple rebar.config files:

```go
func mergeConfigs(paths ...string) (*parser.RebarConfig, error) {
    var allTerms []parser.Term
    
    for _, path := range paths {
        config, err := parser.ParseFile(path)
        if err != nil {
            return nil, fmt.Errorf("failed to parse %s: %w", path, err)
        }
        allTerms = append(allTerms, config.Terms...)
    }
    
    return &parser.RebarConfig{
        Terms: allTerms,
        Raw:   "", // Combined raw content would be complex
    }, nil
}
```

### 3. Configuration Generator

Generate rebar.config programmatically:

```go
func generateConfig(appName string, deps []Dependency) string {
    config := &parser.RebarConfig{
        Terms: []parser.Term{
            parser.Tuple{
                Elements: []parser.Term{
                    parser.Atom{Value: "app_name"},
                    parser.Atom{Value: appName},
                },
            },
            generateDepsConfig(deps),
            generateErlOptsConfig(),
        },
    }
    
    return config.Format(2)
}
```

## Testing Examples

### Unit Testing with the Parser

```go
func TestConfigParsing(t *testing.T) {
    configContent := `
    {erl_opts, [debug_info]}.
    {deps, [{cowboy, "2.9.0"}]}.
    `
    
    config, err := parser.Parse(configContent)
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    // Test erl_opts
    opts, ok := config.GetErlOpts()
    if !ok {
        t.Error("erl_opts not found")
    }
    
    // Test deps
    deps, ok := config.GetDeps()
    if !ok {
        t.Error("deps not found")
    }
    
    // Validate structure
    if len(config.Terms) != 2 {
        t.Errorf("Expected 2 terms, got %d", len(config.Terms))
    }
}
```

### Benchmark Testing

```go
func BenchmarkParsing(b *testing.B) {
    configContent := generateLargeConfig() // Your test data
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := parser.Parse(configContent)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Next Steps

Explore the detailed examples:

1. **[Basic Parsing](./basic-parsing)** - Start with simple parsing examples
2. **[Configuration Access](./config-access)** - Learn to access specific configurations
3. **[Pretty Printing](./pretty-printing)** - Format configurations beautifully
4. **[Term Comparison](./comparison)** - Compare and validate configurations
5. **[Complex Analysis](./complex-analysis)** - Advanced use cases and patterns

Each example page includes:
- Complete, runnable code
- Step-by-step explanations
- Common pitfalls and how to avoid them
- Performance considerations
- Testing strategies
