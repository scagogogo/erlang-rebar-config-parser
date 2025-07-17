# Pretty Printing

This example demonstrates how to format and pretty-print Erlang rebar configurations using the library's formatting capabilities.

## Overview

The library provides a `Format()` method that can pretty-print configurations with configurable indentation, making them more readable and properly formatted.

## Basic Formatting

### Simple Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Parse a compact configuration
    compactConfig := `{erl_opts,[debug_info,warnings_as_errors]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"}]}.`
    
    config, err := parser.Parse(compactConfig)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Println("Original (compact):")
    fmt.Println(compactConfig)
    
    fmt.Println("\nFormatted with 2-space indentation:")
    fmt.Println(config.Format(2))
    
    fmt.Println("\nFormatted with 4-space indentation:")
    fmt.Println(config.Format(4))
}
```

### Expected Output

```
Original (compact):
{erl_opts,[debug_info,warnings_as_errors]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"}]}.

Formatted with 2-space indentation:
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
  {cowboy, "2.9.0"},
  {jsx, "3.1.0"}
]}.

Formatted with 4-space indentation:
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
```

## Advanced Formatting Examples

### Complex Configuration

```go
func demonstrateComplexFormatting() {
    complexConfig := `{erl_opts,[debug_info,warnings_as_errors,{parse_transform,lager_transform}]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"},{lager,"3.9.2"}]}.{profiles,[{dev,[{deps,[{sync,"0.1.3"}]}]},{test,[{deps,[{proper,"1.3.0"},{meck,"0.9.0"}]}]}]}.{relx,[{release,{my_app,"0.1.0"},[my_app,sasl]},{dev_mode,true},{include_erts,false}]}.`
    
    config, err := parser.Parse(complexConfig)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Println("=== Complex Configuration Formatting ===")
    fmt.Println("\nOriginal (all on one line):")
    fmt.Println(complexConfig)
    
    fmt.Println("\nBeautifully formatted:")
    fmt.Println(config.Format(2))
}
```

### Nested Structures

```go
func demonstrateNestedFormatting() {
    nestedConfig := `{profiles,[{test,[{erl_opts,[debug_info,export_all]},{deps,[{proper,"1.3.0"},{meck,"0.9.0"},{ct_helper,"1.1.0"}]},{ct_opts,[{sys_config,"test/sys.config"}]}]},{prod,[{erl_opts,[warnings_as_errors,no_debug_info]},{relx,[{dev_mode,false},{include_erts,true}]}]}]}.`
    
    config, err := parser.Parse(nestedConfig)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Println("=== Nested Structure Formatting ===")
    fmt.Println("\nFormatted nested profiles:")
    fmt.Println(config.Format(2))
}
```

## Formatting Different Indentation Styles

### Comparing Indentation Levels

```go
func compareIndentationStyles(config *parser.RebarConfig) {
    indentLevels := []int{2, 4, 8}
    
    for _, indent := range indentLevels {
        fmt.Printf("=== %d-space indentation ===\n", indent)
        fmt.Println(config.Format(indent))
        fmt.Println()
    }
}
```

### Custom Formatting Function

```go
func formatWithCustomStyle(config *parser.RebarConfig) {
    // You can create wrapper functions for different styles
    
    // Compact style (minimal spacing)
    fmt.Println("Compact style:")
    compact := config.Format(0) // No indentation
    fmt.Println(compact)
    
    // Standard style (2 spaces)
    fmt.Println("\nStandard style:")
    standard := config.Format(2)
    fmt.Println(standard)
    
    // Wide style (4 spaces)
    fmt.Println("\nWide style:")
    wide := config.Format(4)
    fmt.Println(wide)
    
    // Extra wide style (8 spaces)
    fmt.Println("\nExtra wide style:")
    extraWide := config.Format(8)
    fmt.Println(extraWide)
}
```

## Practical Use Cases

### Configuration File Cleanup

```go
func cleanupConfigFile(inputPath, outputPath string) error {
    // Read and parse the configuration
    config, err := parser.ParseFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to parse %s: %w", inputPath, err)
    }
    
    // Format with standard 2-space indentation
    formatted := config.Format(2)
    
    // Write the formatted configuration
    err = os.WriteFile(outputPath, []byte(formatted), 0644)
    if err != nil {
        return fmt.Errorf("failed to write %s: %w", outputPath, err)
    }
    
    fmt.Printf("Cleaned up %s -> %s\n", inputPath, outputPath)
    return nil
}

// Usage
err := cleanupConfigFile("messy_rebar.config", "clean_rebar.config")
if err != nil {
    log.Fatal(err)
}
```

### Configuration Diff Tool

```go
func compareConfigurations(path1, path2 string) {
    config1, err := parser.ParseFile(path1)
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", path1, err)
    }
    
    config2, err := parser.ParseFile(path2)
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", path2, err)
    }
    
    fmt.Printf("=== %s ===\n", path1)
    fmt.Println(config1.Format(2))
    
    fmt.Printf("\n=== %s ===\n", path2)
    fmt.Println(config2.Format(2))
    
    // You could add actual diff logic here
    if len(config1.Terms) != len(config2.Terms) {
        fmt.Printf("\nDifference: %s has %d terms, %s has %d terms\n", 
            path1, len(config1.Terms), path2, len(config2.Terms))
    }
}
```

### Configuration Template Generator

```go
func generateConfigTemplate(appName string, deps []string) string {
    // Create a basic configuration structure
    var terms []parser.Term
    
    // Add application name
    terms = append(terms, parser.Tuple{
        Elements: []parser.Term{
            parser.Atom{Value: "app_name"},
            parser.Atom{Value: appName},
        },
    })
    
    // Add basic erl_opts
    terms = append(terms, parser.Tuple{
        Elements: []parser.Term{
            parser.Atom{Value: "erl_opts"},
            parser.List{
                Elements: []parser.Term{
                    parser.Atom{Value: "debug_info"},
                    parser.Atom{Value: "warnings_as_errors"},
                },
            },
        },
    })
    
    // Add dependencies
    var depElements []parser.Term
    for _, dep := range deps {
        depElements = append(depElements, parser.Tuple{
            Elements: []parser.Term{
                parser.Atom{Value: dep},
                parser.String{Value: "1.0.0"}, // Default version
            },
        })
    }
    
    terms = append(terms, parser.Tuple{
        Elements: []parser.Term{
            parser.Atom{Value: "deps"},
            parser.List{Elements: depElements},
        },
    })
    
    // Create config and format
    config := &parser.RebarConfig{Terms: terms}
    return config.Format(2)
}

// Usage
template := generateConfigTemplate("my_app", []string{"cowboy", "jsx", "lager"})
fmt.Println("Generated template:")
fmt.Println(template)
```

## Formatting Best Practices

### Consistent Formatting

```go
func formatConsistently(configs []string) {
    const standardIndent = 2
    
    for i, configStr := range configs {
        config, err := parser.Parse(configStr)
        if err != nil {
            fmt.Printf("Config %d: Parse error - %v\n", i+1, err)
            continue
        }
        
        fmt.Printf("=== Config %d (formatted) ===\n", i+1)
        fmt.Println(config.Format(standardIndent))
        fmt.Println()
    }
}
```

### Validation Before Formatting

```go
func formatWithValidation(configPath string) {
    config, err := parser.ParseFile(configPath)
    if err != nil {
        log.Fatalf("Parse error: %v", err)
    }
    
    // Validate configuration
    if len(config.Terms) == 0 {
        fmt.Println("Warning: Configuration is empty")
        return
    }
    
    // Check for common sections
    hasErlOpts := false
    hasDeps := false
    
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                switch atom.Value {
                case "erl_opts":
                    hasErlOpts = true
                case "deps":
                    hasDeps = true
                }
            }
        }
    }
    
    if !hasErlOpts {
        fmt.Println("Note: No erl_opts found")
    }
    if !hasDeps {
        fmt.Println("Note: No deps found")
    }
    
    // Format and display
    fmt.Println("Formatted configuration:")
    fmt.Println(config.Format(2))
}
```

## Complete Example

Here's a complete example that demonstrates various formatting scenarios:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Sample messy configuration
    messyConfig := `{erl_opts,[debug_info,warnings_as_errors,{parse_transform,lager_transform}]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"},{lager,"3.9.2"}]}.{profiles,[{dev,[{deps,[{sync,"0.1.3"}]}]},{test,[{deps,[{proper,"1.3.0"}]}]}]}.`
    
    config, err := parser.Parse(messyConfig)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Println("=== Pretty Printing Demo ===")
    
    fmt.Println("\n1. Original (messy):")
    fmt.Println(messyConfig)
    
    fmt.Println("\n2. Formatted (2-space indent):")
    fmt.Println(config.Format(2))
    
    fmt.Println("\n3. Formatted (4-space indent):")
    fmt.Println(config.Format(4))
    
    // Save formatted version
    formatted := config.Format(2)
    err = os.WriteFile("formatted_rebar.config", []byte(formatted), 0644)
    if err != nil {
        log.Printf("Failed to save formatted config: %v", err)
    } else {
        fmt.Println("\n4. Saved formatted version to 'formatted_rebar.config'")
    }
    
    // Demonstrate different use cases
    demonstrateUseCases(config)
}

func demonstrateUseCases(config *parser.RebarConfig) {
    fmt.Println("\n=== Use Case Demonstrations ===")
    
    // Use case 1: Configuration summary
    fmt.Println("\n--- Configuration Summary ---")
    fmt.Printf("Total terms: %d\n", len(config.Terms))
    
    // Use case 2: Extract and format specific sections
    fmt.Println("\n--- Formatted Dependencies ---")
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        // Create a temporary config with just dependencies
        depsConfig := &parser.RebarConfig{
            Terms: []parser.Term{
                parser.Tuple{
                    Elements: []parser.Term{
                        parser.Atom{Value: "deps"},
                        deps[0],
                    },
                },
            },
        }
        fmt.Println(depsConfig.Format(2))
    }
    
    // Use case 3: Different formatting styles
    fmt.Println("\n--- Different Styles ---")
    styles := map[string]int{
        "Compact": 0,
        "Standard": 2,
        "Wide": 4,
        "Extra Wide": 8,
    }
    
    for name, indent := range styles {
        fmt.Printf("\n%s style (%d spaces):\n", name, indent)
        // Show just the first term for brevity
        if len(config.Terms) > 0 {
            singleTermConfig := &parser.RebarConfig{
                Terms: []parser.Term{config.Terms[0]},
            }
            fmt.Println(singleTermConfig.Format(indent))
        }
    }
}
```

This comprehensive example shows how to use the formatting capabilities effectively for various scenarios, from simple cleanup to complex configuration management tasks.
