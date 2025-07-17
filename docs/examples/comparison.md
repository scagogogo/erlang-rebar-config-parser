# Term Comparison

This example demonstrates how to compare Erlang terms and configurations using the library's comparison functionality.

## Overview

The library provides comparison capabilities through the `Compare()` method implemented by all term types. This allows you to check equality between different terms and configurations.

## Basic Term Comparison

### Simple Comparisons

```go
package main

import (
    "fmt"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Create some terms for comparison
    atom1 := parser.Atom{Value: "debug_info", IsQuoted: false}
    atom2 := parser.Atom{Value: "debug_info", IsQuoted: true}  // Different quoting
    atom3 := parser.Atom{Value: "warnings_as_errors", IsQuoted: false}
    
    str1 := parser.String{Value: "hello"}
    str2 := parser.String{Value: "hello"}
    str3 := parser.String{Value: "world"}
    
    int1 := parser.Integer{Value: 42}
    int2 := parser.Integer{Value: 42}
    int3 := parser.Integer{Value: 24}
    
    // Atom comparisons (ignores IsQuoted)
    fmt.Printf("atom1.Compare(atom2): %t\n", atom1.Compare(atom2)) // true
    fmt.Printf("atom1.Compare(atom3): %t\n", atom1.Compare(atom3)) // false
    
    // String comparisons
    fmt.Printf("str1.Compare(str2): %t\n", str1.Compare(str2)) // true
    fmt.Printf("str1.Compare(str3): %t\n", str1.Compare(str3)) // false
    
    // Integer comparisons
    fmt.Printf("int1.Compare(int2): %t\n", int1.Compare(int2)) // true
    fmt.Printf("int1.Compare(int3): %t\n", int1.Compare(int3)) // false
    
    // Cross-type comparisons (always false)
    fmt.Printf("atom1.Compare(str1): %t\n", atom1.Compare(str1)) // false
    fmt.Printf("str1.Compare(int1): %t\n", str1.Compare(int1))   // false
}
```

### Complex Structure Comparisons

```go
func demonstrateComplexComparisons() {
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
    
    fmt.Printf("tuple1.Compare(tuple2): %t\n", tuple1.Compare(tuple2)) // true
    fmt.Printf("tuple1.Compare(tuple3): %t\n", tuple1.Compare(tuple3)) // false
    
    // Create identical lists
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
    
    fmt.Printf("list1.Compare(list2): %t\n", list1.Compare(list2)) // true
    fmt.Printf("list1.Compare(list3): %t\n", list1.Compare(list3)) // false
}
```

## Configuration Comparison

### Comparing Entire Configurations

```go
func compareConfigurations(config1Str, config2Str string) {
    config1, err := parser.Parse(config1Str)
    if err != nil {
        fmt.Printf("Failed to parse config1: %v\n", err)
        return
    }
    
    config2, err := parser.Parse(config2Str)
    if err != nil {
        fmt.Printf("Failed to parse config2: %v\n", err)
        return
    }
    
    fmt.Println("=== Configuration Comparison ===")
    
    // Compare number of terms
    if len(config1.Terms) != len(config2.Terms) {
        fmt.Printf("Different number of terms: %d vs %d\n", 
            len(config1.Terms), len(config2.Terms))
        return
    }
    
    // Compare each term
    allEqual := true
    for i, term1 := range config1.Terms {
        term2 := config2.Terms[i]
        if !term1.Compare(term2) {
            fmt.Printf("Term %d differs:\n", i+1)
            fmt.Printf("  Config1: %s\n", term1.String())
            fmt.Printf("  Config2: %s\n", term2.String())
            allEqual = false
        }
    }
    
    if allEqual {
        fmt.Println("✓ Configurations are identical")
    } else {
        fmt.Println("✗ Configurations differ")
    }
}

// Usage example
func main() {
    config1 := `{erl_opts, [debug_info]}.{deps, [{cowboy, "2.9.0"}]}.`
    config2 := `{erl_opts, [debug_info]}.{deps, [{cowboy, "2.9.0"}]}.`
    config3 := `{erl_opts, [debug_info]}.{deps, [{cowboy, "2.8.0"}]}.`
    
    compareConfigurations(config1, config2) // Should be identical
    compareConfigurations(config1, config3) // Should differ
}
```

### Comparing Specific Sections

```go
func compareSpecificSections(config1, config2 *parser.RebarConfig) {
    fmt.Println("=== Section-by-Section Comparison ===")
    
    // Compare dependencies
    deps1, ok1 := config1.GetDeps()
    deps2, ok2 := config2.GetDeps()
    
    if ok1 && ok2 {
        if len(deps1) > 0 && len(deps2) > 0 {
            if deps1[0].Compare(deps2[0]) {
                fmt.Println("✓ Dependencies are identical")
            } else {
                fmt.Println("✗ Dependencies differ")
                fmt.Printf("  Config1 deps: %s\n", deps1[0].String())
                fmt.Printf("  Config2 deps: %s\n", deps2[0].String())
            }
        }
    } else if ok1 != ok2 {
        fmt.Println("✗ One config has dependencies, the other doesn't")
    } else {
        fmt.Println("- Neither config has dependencies")
    }
    
    // Compare Erlang options
    opts1, ok1 := config1.GetErlOpts()
    opts2, ok2 := config2.GetErlOpts()
    
    if ok1 && ok2 {
        if len(opts1) > 0 && len(opts2) > 0 {
            if opts1[0].Compare(opts2[0]) {
                fmt.Println("✓ Erlang options are identical")
            } else {
                fmt.Println("✗ Erlang options differ")
                fmt.Printf("  Config1 erl_opts: %s\n", opts1[0].String())
                fmt.Printf("  Config2 erl_opts: %s\n", opts2[0].String())
            }
        }
    } else if ok1 != ok2 {
        fmt.Println("✗ One config has erl_opts, the other doesn't")
    } else {
        fmt.Println("- Neither config has erl_opts")
    }
    
    // Compare application names
    app1, ok1 := config1.GetAppName()
    app2, ok2 := config2.GetAppName()
    
    if ok1 && ok2 {
        if app1 == app2 {
            fmt.Printf("✓ Application names are identical: %s\n", app1)
        } else {
            fmt.Printf("✗ Application names differ: %s vs %s\n", app1, app2)
        }
    } else if ok1 != ok2 {
        fmt.Println("✗ One config has app_name, the other doesn't")
    } else {
        fmt.Println("- Neither config has app_name")
    }
}
```

## Advanced Comparison Utilities

### Dependency Comparison

```go
func compareDependencies(config1, config2 *parser.RebarConfig) {
    fmt.Println("=== Dependency Analysis ===")
    
    deps1 := extractDependencyNames(config1)
    deps2 := extractDependencyNames(config2)
    
    // Find common dependencies
    common := findCommonDeps(deps1, deps2)
    if len(common) > 0 {
        fmt.Printf("Common dependencies (%d): %v\n", len(common), common)
    }
    
    // Find unique dependencies
    unique1 := findUniqueDeps(deps1, deps2)
    if len(unique1) > 0 {
        fmt.Printf("Only in config1 (%d): %v\n", len(unique1), unique1)
    }
    
    unique2 := findUniqueDeps(deps2, deps1)
    if len(unique2) > 0 {
        fmt.Printf("Only in config2 (%d): %v\n", len(unique2), unique2)
    }
    
    // Compare versions for common dependencies
    compareVersions(config1, config2, common)
}

func extractDependencyNames(config *parser.RebarConfig) []string {
    var names []string
    
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                        names = append(names, atom.Value)
                    }
                }
            }
        }
    }
    
    return names
}

func findCommonDeps(deps1, deps2 []string) []string {
    var common []string
    for _, dep1 := range deps1 {
        for _, dep2 := range deps2 {
            if dep1 == dep2 {
                common = append(common, dep1)
                break
            }
        }
    }
    return common
}

func findUniqueDeps(deps1, deps2 []string) []string {
    var unique []string
    for _, dep1 := range deps1 {
        found := false
        for _, dep2 := range deps2 {
            if dep1 == dep2 {
                found = true
                break
            }
        }
        if !found {
            unique = append(unique, dep1)
        }
    }
    return unique
}

func compareVersions(config1, config2 *parser.RebarConfig, commonDeps []string) {
    if len(commonDeps) == 0 {
        return
    }
    
    fmt.Println("\n--- Version Comparison ---")
    
    for _, depName := range commonDeps {
        version1 := getDependencyVersion(config1, depName)
        version2 := getDependencyVersion(config2, depName)
        
        if version1 == version2 {
            fmt.Printf("✓ %s: %s (same)\n", depName, version1)
        } else {
            fmt.Printf("✗ %s: %s vs %s\n", depName, version1, version2)
        }
    }
}

func getDependencyVersion(config *parser.RebarConfig, depName string) string {
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == depName {
                        if str, ok := tuple.Elements[1].(parser.String); ok {
                            return str.Value
                        }
                        return tuple.Elements[1].String()
                    }
                }
            }
        }
    }
    return "unknown"
}
```

### Configuration Diff Tool

```go
func generateConfigDiff(config1, config2 *parser.RebarConfig) {
    fmt.Println("=== Configuration Diff ===")
    
    // Compare basic structure
    fmt.Printf("Config1 terms: %d\n", len(config1.Terms))
    fmt.Printf("Config2 terms: %d\n", len(config2.Terms))
    
    // Create maps for easier comparison
    terms1 := createTermMap(config1)
    terms2 := createTermMap(config2)
    
    // Find all unique keys
    allKeys := make(map[string]bool)
    for key := range terms1 {
        allKeys[key] = true
    }
    for key := range terms2 {
        allKeys[key] = true
    }
    
    // Compare each section
    for key := range allKeys {
        term1, exists1 := terms1[key]
        term2, exists2 := terms2[key]
        
        if !exists1 {
            fmt.Printf("+ %s: %s (only in config2)\n", key, term2.String())
        } else if !exists2 {
            fmt.Printf("- %s: %s (only in config1)\n", key, term1.String())
        } else if !term1.Compare(term2) {
            fmt.Printf("~ %s:\n", key)
            fmt.Printf("  - %s\n", term1.String())
            fmt.Printf("  + %s\n", term2.String())
        } else {
            fmt.Printf("= %s: identical\n", key)
        }
    }
}

func createTermMap(config *parser.RebarConfig) map[string]parser.Term {
    termMap := make(map[string]parser.Term)
    
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                termMap[atom.Value] = term
            }
        }
    }
    
    return termMap
}
```

## Complete Comparison Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Sample configurations for comparison
    config1Str := `
    {app_name, my_app}.
    {erl_opts, [debug_info, warnings_as_errors]}.
    {deps, [
        {cowboy, "2.9.0"},
        {jsx, "3.1.0"}
    ]}.
    `
    
    config2Str := `
    {app_name, my_app}.
    {erl_opts, [debug_info, warnings_as_errors]}.
    {deps, [
        {cowboy, "2.8.0"},
        {jsx, "3.1.0"},
        {lager, "3.9.2"}
    ]}.
    `
    
    config1, err := parser.Parse(config1Str)
    if err != nil {
        log.Fatalf("Failed to parse config1: %v", err)
    }
    
    config2, err := parser.Parse(config2Str)
    if err != nil {
        log.Fatalf("Failed to parse config2: %v", err)
    }
    
    fmt.Println("=== Comprehensive Configuration Comparison ===")
    
    // 1. Basic comparison
    compareConfigurations(config1Str, config2Str)
    
    fmt.Println()
    
    // 2. Section-by-section comparison
    compareSpecificSections(config1, config2)
    
    fmt.Println()
    
    // 3. Dependency analysis
    compareDependencies(config1, config2)
    
    fmt.Println()
    
    // 4. Generate diff
    generateConfigDiff(config1, config2)
    
    // 5. Demonstrate term-level comparisons
    fmt.Println("\n=== Term-Level Comparisons ===")
    demonstrateTermComparisons()
}

func demonstrateTermComparisons() {
    // Create various terms for comparison
    terms := []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "debug_info", IsQuoted: true},
        parser.String{Value: "2.9.0"},
        parser.Integer{Value: 42},
        parser.Float{Value: 3.14},
        parser.Tuple{
            Elements: []parser.Term{
                parser.Atom{Value: "cowboy"},
                parser.String{Value: "2.9.0"},
            },
        },
        parser.List{
            Elements: []parser.Term{
                parser.Atom{Value: "debug_info"},
                parser.Atom{Value: "warnings_as_errors"},
            },
        },
    }
    
    // Compare each term with itself and others
    for i, term1 := range terms {
        for j, term2 := range terms {
            result := term1.Compare(term2)
            if i == j {
                fmt.Printf("✓ Term %d == Term %d: %t (self-comparison)\n", i+1, j+1, result)
            } else if result {
                fmt.Printf("✓ Term %d == Term %d: %t\n", i+1, j+1, result)
                fmt.Printf("  %s == %s\n", term1.String(), term2.String())
            }
        }
    }
}
```

This comprehensive example demonstrates all aspects of term and configuration comparison, from simple equality checks to complex diff generation and dependency analysis.
