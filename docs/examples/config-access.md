# Configuration Access

This example demonstrates how to use the helper methods provided by the library to easily access common configuration elements in rebar.config files.

## Overview

The `RebarConfig` type provides several helper methods that make it easy to access commonly used configuration sections without manually parsing tuples and lists.

## Available Helper Methods

- `GetDeps()` - Get dependencies
- `GetErlOpts()` - Get Erlang compilation options
- `GetAppName()` - Get application name
- `GetPlugins()` - Get plugins
- `GetProfilesConfig()` - Get build profiles
- `GetRelxConfig()` - Get release configuration

## Basic Configuration Access

### Sample Configuration

```erlang
{app_name, my_awesome_app}.

{erl_opts, [
    debug_info,
    warnings_as_errors,
    {parse_transform, lager_transform}
]}.

{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"},
    {lager, "3.9.2"}
]}.

{plugins, [
    rebar3_hex,
    rebar3_auto
]}.
```

### Accessing Configuration

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
    
    // Get application name
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("Application: %s\n", appName)
    } else {
        fmt.Println("No application name found")
    }
    
    // Get dependencies
    if deps, ok := config.GetDeps(); ok {
        fmt.Println("Dependencies found!")
        processDependencies(deps)
    } else {
        fmt.Println("No dependencies found")
    }
    
    // Get Erlang options
    if erlOpts, ok := config.GetErlOpts(); ok {
        fmt.Println("Erlang options found!")
        processErlangOptions(erlOpts)
    } else {
        fmt.Println("No Erlang options found")
    }
    
    // Get plugins
    if plugins, ok := config.GetPlugins(); ok {
        fmt.Println("Plugins found!")
        processPlugins(plugins)
    } else {
        fmt.Println("No plugins found")
    }
}

func processDependencies(deps []parser.Term) {
    if len(deps) == 0 {
        return
    }
    
    if depsList, ok := deps[0].(parser.List); ok {
        fmt.Printf("  Found %d dependencies:\n", len(depsList.Elements))
        
        for _, dep := range depsList.Elements {
            if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                if name, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("    - %s", name.Value)
                    
                    // Get version
                    switch version := tuple.Elements[1].(type) {
                    case parser.String:
                        fmt.Printf(" (version: %s)", version.Value)
                    case parser.Tuple:
                        fmt.Printf(" (version spec: %s)", version.String())
                    }
                    fmt.Println()
                }
            }
        }
    }
}

func processErlangOptions(erlOpts []parser.Term) {
    if len(erlOpts) == 0 {
        return
    }
    
    if optsList, ok := erlOpts[0].(parser.List); ok {
        fmt.Printf("  Found %d Erlang options:\n", len(optsList.Elements))
        
        for _, opt := range optsList.Elements {
            switch o := opt.(type) {
            case parser.Atom:
                fmt.Printf("    - %s\n", o.Value)
            case parser.Tuple:
                fmt.Printf("    - %s\n", o.String())
            default:
                fmt.Printf("    - %s\n", o.String())
            }
        }
    }
}

func processPlugins(plugins []parser.Term) {
    if len(plugins) == 0 {
        return
    }
    
    if pluginsList, ok := plugins[0].(parser.List); ok {
        fmt.Printf("  Found %d plugins:\n", len(pluginsList.Elements))
        
        for _, plugin := range pluginsList.Elements {
            if atom, ok := plugin.(parser.Atom); ok {
                fmt.Printf("    - %s\n", atom.Value)
            }
        }
    }
}
```

## Working with Profiles

Profiles allow different configurations for different environments (dev, test, prod).

### Sample Profile Configuration

```erlang
{profiles, [
    {dev, [
        {deps, [
            {sync, "0.1.3"},
            {observer_cli, "1.7.3"}
        ]},
        {erl_opts, [debug_info]}
    ]},
    {test, [
        {deps, [
            {proper, "1.3.0"},
            {meck, "0.9.0"}
        ]},
        {erl_opts, [debug_info, export_all]}
    ]},
    {prod, [
        {erl_opts, [warnings_as_errors, no_debug_info]}
    ]}
]}.
```

### Accessing Profile Configuration

```go
func analyzeProfiles(config *parser.RebarConfig) {
    profiles, ok := config.GetProfilesConfig()
    if !ok {
        fmt.Println("No profiles found")
        return
    }
    
    if len(profiles) == 0 {
        return
    }
    
    if profilesList, ok := profiles[0].(parser.List); ok {
        fmt.Printf("Found %d profiles:\n", len(profilesList.Elements))
        
        for _, profile := range profilesList.Elements {
            if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                if name, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("\nProfile: %s\n", name.Value)
                    
                    if configList, ok := tuple.Elements[1].(parser.List); ok {
                        analyzeProfileConfig(configList)
                    }
                }
            }
        }
    }
}

func analyzeProfileConfig(configList parser.List) {
    for _, configItem := range configList.Elements {
        if tuple, ok := configItem.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
            if key, ok := tuple.Elements[0].(parser.Atom); ok {
                fmt.Printf("  %s: %s\n", key.Value, tuple.Elements[1].String())
            }
        }
    }
}
```

## Advanced Configuration Access

### Custom Configuration Sections

For configuration sections not covered by helper methods, use the generic `GetTerm` method:

```go
func accessCustomConfig(config *parser.RebarConfig) {
    // Access custom configuration
    if term, ok := config.GetTerm("custom_config"); ok {
        fmt.Printf("Custom config: %s\n", term.String())
    }
    
    // Access relx configuration
    if term, ok := config.GetTerm("relx"); ok {
        fmt.Printf("Relx config: %s\n", term.String())
        processRelxConfig(term)
    }
    
    // Access shell configuration
    if term, ok := config.GetTerm("shell"); ok {
        fmt.Printf("Shell config: %s\n", term.String())
    }
}

func processRelxConfig(term parser.Term) {
    if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
        if configList, ok := tuple.Elements[1].(parser.List); ok {
            fmt.Println("  Relx configuration items:")
            for _, item := range configList.Elements {
                fmt.Printf("    - %s\n", item.String())
            }
        }
    }
}
```

### Configuration Validation

```go
func validateConfiguration(config *parser.RebarConfig) []string {
    var warnings []string
    
    // Check for required sections
    if _, ok := config.GetAppName(); !ok {
        warnings = append(warnings, "Missing application name")
    }
    
    if _, ok := config.GetDeps(); !ok {
        warnings = append(warnings, "No dependencies defined")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        warnings = append(warnings, "No Erlang options defined")
    }
    
    // Check for recommended options
    if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
        if optsList, ok := erlOpts[0].(parser.List); ok {
            hasDebugInfo := false
            for _, opt := range optsList.Elements {
                if atom, ok := opt.(parser.Atom); ok && atom.Value == "debug_info" {
                    hasDebugInfo = true
                    break
                }
            }
            if !hasDebugInfo {
                warnings = append(warnings, "Recommended: add debug_info to erl_opts")
            }
        }
    }
    
    return warnings
}

// Usage
warnings := validateConfiguration(config)
if len(warnings) > 0 {
    fmt.Println("Configuration warnings:")
    for _, warning := range warnings {
        fmt.Printf("  - %s\n", warning)
    }
}
```

## Complete Example

Here's a complete example that demonstrates comprehensive configuration access:

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
    
    fmt.Println("=== Configuration Analysis ===")
    
    // Basic information
    analyzeBasicInfo(config)
    
    // Dependencies
    analyzeDependencies(config)
    
    // Erlang options
    analyzeErlangOptions(config)
    
    // Plugins
    analyzePlugins(config)
    
    // Profiles
    analyzeProfiles(config)
    
    // Validation
    validateAndReport(config)
}

func analyzeBasicInfo(config *parser.RebarConfig) {
    fmt.Println("\n--- Basic Information ---")
    
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("Application: %s\n", appName)
    } else {
        fmt.Println("Application: Not specified")
    }
    
    fmt.Printf("Total configuration terms: %d\n", len(config.Terms))
}

func analyzeDependencies(config *parser.RebarConfig) {
    fmt.Println("\n--- Dependencies ---")
    
    deps, ok := config.GetDeps()
    if !ok {
        fmt.Println("No dependencies found")
        return
    }
    
    if len(deps) == 0 {
        fmt.Println("Dependencies section is empty")
        return
    }
    
    if depsList, ok := deps[0].(parser.List); ok {
        fmt.Printf("Found %d dependencies:\n", len(depsList.Elements))
        
        for i, dep := range depsList.Elements {
            fmt.Printf("%d. %s\n", i+1, dep.String())
        }
    }
}

func analyzeErlangOptions(config *parser.RebarConfig) {
    fmt.Println("\n--- Erlang Options ---")
    
    erlOpts, ok := config.GetErlOpts()
    if !ok {
        fmt.Println("No Erlang options found")
        return
    }
    
    if len(erlOpts) == 0 {
        fmt.Println("Erlang options section is empty")
        return
    }
    
    if optsList, ok := erlOpts[0].(parser.List); ok {
        fmt.Printf("Found %d Erlang options:\n", len(optsList.Elements))
        
        for i, opt := range optsList.Elements {
            fmt.Printf("%d. %s\n", i+1, opt.String())
        }
    }
}

func analyzePlugins(config *parser.RebarConfig) {
    fmt.Println("\n--- Plugins ---")
    
    plugins, ok := config.GetPlugins()
    if !ok {
        fmt.Println("No plugins found")
        return
    }
    
    if len(plugins) == 0 {
        fmt.Println("Plugins section is empty")
        return
    }
    
    if pluginsList, ok := plugins[0].(parser.List); ok {
        fmt.Printf("Found %d plugins:\n", len(pluginsList.Elements))
        
        for i, plugin := range pluginsList.Elements {
            fmt.Printf("%d. %s\n", i+1, plugin.String())
        }
    }
}

func validateAndReport(config *parser.RebarConfig) {
    fmt.Println("\n--- Validation ---")
    
    warnings := validateConfiguration(config)
    if len(warnings) == 0 {
        fmt.Println("✓ Configuration looks good!")
    } else {
        fmt.Printf("Found %d warnings:\n", len(warnings))
        for i, warning := range warnings {
            fmt.Printf("%d. %s\n", i+1, warning)
        }
    }
}
```

This example provides a comprehensive analysis of a rebar.config file, demonstrating how to use all the helper methods effectively.
