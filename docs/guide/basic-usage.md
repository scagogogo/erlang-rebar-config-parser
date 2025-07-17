# Basic Usage

This guide covers the fundamental usage patterns of the Erlang Rebar Config Parser. After reading this, you'll understand how to parse configurations, access common elements, and handle different data types.

## Parsing Configurations

### From Files

The most common use case is parsing an existing rebar.config file:

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
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Printf("Successfully parsed %d configuration terms\n", len(config.Terms))
}
```

### From Strings

When you have configuration content as a string:

```go
configContent := `
{erl_opts, [debug_info, warnings_as_errors]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
`

config, err := parser.Parse(configContent)
if err != nil {
    log.Fatalf("Failed to parse config: %v", err)
}
```

### From Readers

For reading from any io.Reader (files, HTTP responses, etc.):

```go
import (
    "net/http"
    "os"
)

// From file
file, err := os.Open("rebar.config")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

config, err := parser.ParseReader(file)

// From HTTP response
resp, err := http.Get("https://example.com/rebar.config")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

config, err = parser.ParseReader(resp.Body)
```

## Accessing Configuration Elements

### Using Helper Methods

The library provides convenient helper methods for common configuration sections:

```go
config, _ := parser.ParseFile("rebar.config")

// Get dependencies
if deps, ok := config.GetDeps(); ok {
    fmt.Println("Dependencies found!")
    // deps contains the dependency terms
}

// Get Erlang compilation options
if erlOpts, ok := config.GetErlOpts(); ok {
    fmt.Println("Erlang options found!")
    // erlOpts contains the compilation options
}

// Get application name
if appName, ok := config.GetAppName(); ok {
    fmt.Printf("Application name: %s\n", appName)
}

// Get plugins
if plugins, ok := config.GetPlugins(); ok {
    fmt.Println("Plugins found!")
}

// Get profiles
if profiles, ok := config.GetProfilesConfig(); ok {
    fmt.Println("Profiles found!")
}
```

### Manual Term Access

For custom configuration sections, use the generic term access:

```go
// Get any named term
if term, ok := config.GetTerm("custom_config"); ok {
    fmt.Printf("Custom config: %s\n", term.String())
}

// Get tuple elements
if elements, ok := config.GetTupleElements("my_tuple"); ok {
    fmt.Printf("Tuple has %d elements\n", len(elements))
}
```

## Working with Different Term Types

### Type Checking and Conversion

Always use safe type assertions when working with terms:

```go
func processTerm(term parser.Term) {
    switch t := term.(type) {
    case parser.Atom:
        fmt.Printf("Atom: %s", t.Value)
        if t.IsQuoted {
            fmt.Print(" (quoted)")
        }
        fmt.Println()
        
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
```

### Working with Collections

#### Processing Lists

```go
func processList(list parser.List) {
    fmt.Printf("Processing list with %d elements:\n", len(list.Elements))
    
    for i, element := range list.Elements {
        fmt.Printf("  [%d]: %s\n", i, element.String())
    }
}

// Example: Process erl_opts
if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
    if optsList, ok := erlOpts[0].(parser.List); ok {
        processList(optsList)
    }
}
```

#### Processing Tuples

```go
func processTuple(tuple parser.Tuple) {
    fmt.Printf("Processing tuple with %d elements:\n", len(tuple.Elements))
    
    for i, element := range tuple.Elements {
        fmt.Printf("  [%d]: %s\n", i, element.String())
    }
}

// Example: Process a dependency tuple
func processDependency(dep parser.Term) {
    if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
        if name, ok := tuple.Elements[0].(parser.Atom); ok {
            fmt.Printf("Dependency: %s\n", name.Value)
            
            // Version can be string or tuple
            switch version := tuple.Elements[1].(type) {
            case parser.String:
                fmt.Printf("  Version: %s\n", version.Value)
            case parser.Tuple:
                fmt.Printf("  Version spec: %s\n", version.String())
            }
        }
    }
}
```

## Common Patterns

### Extracting Dependencies

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

// Usage
deps := extractDependencies(config)
fmt.Printf("Found dependencies: %v\n", deps)
```

### Checking for Specific Options

```go
func hasDebugInfo(config *parser.RebarConfig) bool {
    if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
        if optsList, ok := erlOpts[0].(parser.List); ok {
            for _, opt := range optsList.Elements {
                if atom, ok := opt.(parser.Atom); ok && atom.Value == "debug_info" {
                    return true
                }
            }
        }
    }
    return false
}

if hasDebugInfo(config) {
    fmt.Println("Debug info is enabled")
}
```

### Finding Specific Profiles

```go
func hasProfile(config *parser.RebarConfig, profileName string) bool {
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            for _, profile := range profilesList.Elements {
                if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == profileName {
                        return true
                    }
                }
            }
        }
    }
    return false
}

if hasProfile(config, "test") {
    fmt.Println("Test profile is configured")
}
```

## Error Handling

### Comprehensive Error Handling

```go
func parseConfigSafely(path string) (*parser.RebarConfig, error) {
    config, err := parser.ParseFile(path)
    if err != nil {
        // Provide context for different error types
        if strings.Contains(err.Error(), "no such file") {
            return nil, fmt.Errorf("configuration file '%s' not found", path)
        } else if strings.Contains(err.Error(), "permission denied") {
            return nil, fmt.Errorf("permission denied reading '%s'", path)
        } else if strings.Contains(err.Error(), "syntax error") {
            return nil, fmt.Errorf("invalid syntax in '%s': %w", path, err)
        }
        return nil, fmt.Errorf("failed to parse '%s': %w", path, err)
    }
    return config, nil
}
```

### Validation After Parsing

```go
func validateConfig(config *parser.RebarConfig) error {
    // Check for required sections
    if _, ok := config.GetDeps(); !ok {
        return fmt.Errorf("missing required 'deps' configuration")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        return fmt.Errorf("missing required 'erl_opts' configuration")
    }
    
    // Validate application name
    if appName, ok := config.GetAppName(); ok {
        if appName == "" {
            return fmt.Errorf("application name cannot be empty")
        }
    } else {
        return fmt.Errorf("missing application name")
    }
    
    return nil
}

// Usage
config, err := parser.ParseFile("rebar.config")
if err != nil {
    log.Fatal(err)
}

if err := validateConfig(config); err != nil {
    log.Fatalf("Configuration validation failed: %v", err)
}
```

## Formatting and Display

### Pretty Printing

```go
// Format with different indentation levels
formatted2 := config.Format(2)  // 2-space indentation
formatted4 := config.Format(4)  // 4-space indentation

fmt.Println("2-space indentation:")
fmt.Println(formatted2)

fmt.Println("\n4-space indentation:")
fmt.Println(formatted4)
```

### Custom Display Functions

```go
func displayConfig(config *parser.RebarConfig) {
    fmt.Printf("Configuration Summary:\n")
    fmt.Printf("  Total terms: %d\n", len(config.Terms))
    
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("  Application: %s\n", appName)
    }
    
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("  Dependencies: %d\n", len(depsList.Elements))
        }
    }
    
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            fmt.Printf("  Profiles: %d\n", len(profilesList.Elements))
        }
    }
    
    fmt.Println("\nFormatted configuration:")
    fmt.Println(config.Format(2))
}
```

## Next Steps

Now that you understand the basics:

1. **[Advanced Usage](./advanced-usage)** - Learn about complex scenarios and best practices
2. **[API Reference](../api/)** - Explore the complete API documentation
3. **[Examples](../examples/)** - See real-world examples and use cases
