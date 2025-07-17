# Basic Parsing

This example demonstrates the fundamental parsing capabilities of the Erlang Rebar Config Parser. You'll learn how to parse configurations from different sources and access the parsed data.

## Simple File Parsing

Let's start with a basic example that parses a rebar.config file:

### Sample rebar.config

```erlang
%% Basic rebar.config
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"},
    {lager, "3.9.2"}
]}.

{app_name, my_awesome_app}.
```

### Parsing Code

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Parse the configuration file
    config, err := parser.ParseFile("rebar.config")
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Display basic information
    fmt.Printf("Successfully parsed configuration!\n")
    fmt.Printf("Number of top-level terms: %d\n", len(config.Terms))
    
    // Iterate through all terms
    for i, term := range config.Terms {
        fmt.Printf("Term %d: %s\n", i+1, term.String())
    }
}
```

### Expected Output

```
Successfully parsed configuration!
Number of top-level terms: 3
Term 1: {erl_opts, [debug_info, warnings_as_errors]}
Term 2: {deps, [{cowboy, "2.9.0"}, {jsx, "3.1.0"}, {lager, "3.9.2"}]}
Term 3: {app_name, my_awesome_app}
```

## Parsing from String

Sometimes you need to parse configuration content from a string:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Configuration as string
    configContent := `
    {erl_opts, [debug_info]}.
    {deps, [
        {cowboy, "2.9.0"}
    ]}.
    `
    
    // Parse from string
    config, err := parser.Parse(configContent)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Printf("Parsed %d terms from string\n", len(config.Terms))
    
    // Access the raw content
    fmt.Printf("Raw content length: %d characters\n", len(config.Raw))
}
```

## Parsing from Reader

For more advanced scenarios, you can parse from any `io.Reader`:

```go
package main

import (
    "bytes"
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Example 1: Parse from string reader
    configStr := `{erl_opts, [debug_info]}.`
    reader := strings.NewReader(configStr)
    
    config, err := parser.ParseReader(reader)
    if err != nil {
        log.Fatalf("Failed to parse from reader: %v", err)
    }
    
    fmt.Printf("Parsed from string reader: %d terms\n", len(config.Terms))
    
    // Example 2: Parse from bytes buffer
    var buf bytes.Buffer
    buf.WriteString(`{deps, [{jsx, "3.1.0"}]}.`)
    
    config2, err := parser.ParseReader(&buf)
    if err != nil {
        log.Fatalf("Failed to parse from buffer: %v", err)
    }
    
    fmt.Printf("Parsed from buffer: %d terms\n", len(config2.Terms))
}
```

## Understanding Term Types

Let's examine different Erlang term types and how they're represented:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Configuration with various term types
    configContent := `
    {atom_example, simple_atom}.
    {quoted_atom_example, 'quoted-atom'}.
    {string_example, "hello world"}.
    {integer_example, 42}.
    {float_example, 3.14}.
    {negative_number, -123}.
    {scientific_notation, 1.5e-3}.
    {tuple_example, {key, value, 123}}.
    {list_example, [item1, item2, item3]}.
    {nested_example, {deps, [{cowboy, "2.9.0"}]}}.
    `
    
    config, err := parser.Parse(configContent)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Examine each term
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
            name := tuple.Elements[0]
            value := tuple.Elements[1]
            
            fmt.Printf("Configuration: %s\n", name.String())
            fmt.Printf("  Type: %T\n", value)
            fmt.Printf("  Value: %s\n", value.String())
            
            // Detailed type analysis
            analyzeTermType(value)
            fmt.Println()
        }
    }
}

func analyzeTermType(term parser.Term) {
    switch t := term.(type) {
    case parser.Atom:
        fmt.Printf("  Atom details: value='%s', quoted=%t\n", t.Value, t.IsQuoted)
    case parser.String:
        fmt.Printf("  String details: value='%s', length=%d\n", t.Value, len(t.Value))
    case parser.Integer:
        fmt.Printf("  Integer details: value=%d\n", t.Value)
    case parser.Float:
        fmt.Printf("  Float details: value=%f\n", t.Value)
    case parser.Tuple:
        fmt.Printf("  Tuple details: %d elements\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("    [%d]: %s (%T)\n", i, elem.String(), elem)
        }
    case parser.List:
        fmt.Printf("  List details: %d elements\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("    [%d]: %s (%T)\n", i, elem.String(), elem)
        }
    }
}
```

## Error Handling Examples

### Handling Parse Errors

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Example of invalid configurations
    invalidConfigs := []string{
        `{incomplete_tuple`,           // Missing closing brace
        `{unterminated, "string}`,     // Unterminated string
        `{invalid_number, 12.34.56}`,  // Invalid number format
        `{missing_comma [item1 item2]}`, // Missing comma in list
    }
    
    for i, configStr := range invalidConfigs {
        fmt.Printf("Testing invalid config %d:\n", i+1)
        fmt.Printf("Content: %s\n", configStr)
        
        _, err := parser.Parse(configStr)
        if err != nil {
            fmt.Printf("Error (expected): %v\n", err)
            
            // Categorize error types
            if strings.Contains(err.Error(), "syntax error") {
                fmt.Println("  Category: Syntax Error")
            } else if strings.Contains(err.Error(), "unexpected") {
                fmt.Println("  Category: Unexpected Character")
            } else {
                fmt.Println("  Category: Other Parse Error")
            }
        } else {
            fmt.Println("Unexpected: No error occurred!")
        }
        fmt.Println()
    }
}
```

### Robust File Parsing

```go
package main

import (
    "fmt"
    "os"
    "strings"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func parseConfigFile(filename string) {
    fmt.Printf("Attempting to parse: %s\n", filename)
    
    config, err := parser.ParseFile(filename)
    if err != nil {
        // Handle different types of errors
        if os.IsNotExist(err) {
            fmt.Printf("  Error: File '%s' does not exist\n", filename)
        } else if os.IsPermission(err) {
            fmt.Printf("  Error: Permission denied for file '%s'\n", filename)
        } else if strings.Contains(err.Error(), "syntax error") {
            fmt.Printf("  Error: Invalid syntax in '%s': %v\n", filename, err)
        } else {
            fmt.Printf("  Error: Failed to parse '%s': %v\n", filename, err)
        }
        return
    }
    
    fmt.Printf("  Success: Parsed %d terms\n", len(config.Terms))
    
    // Display summary
    if len(config.Terms) > 0 {
        fmt.Printf("  First term: %s\n", config.Terms[0].String())
    }
}

func main() {
    // Test with various files
    testFiles := []string{
        "rebar.config",           // Existing file
        "nonexistent.config",     // Non-existent file
        "/etc/passwd",            // File with wrong format
        "test_configs/valid.config", // Valid test file
    }
    
    for _, filename := range testFiles {
        parseConfigFile(filename)
        fmt.Println()
    }
}
```

## Performance Considerations

### Measuring Parse Time

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Large configuration for performance testing
    largeConfig := generateLargeConfig(1000) // Generate config with 1000 dependencies
    
    // Measure parsing time
    start := time.Now()
    config, err := parser.Parse(largeConfig)
    duration := time.Since(start)
    
    if err != nil {
        fmt.Printf("Parse failed: %v\n", err)
        return
    }
    
    fmt.Printf("Parsed %d terms in %v\n", len(config.Terms), duration)
    fmt.Printf("Average time per term: %v\n", duration/time.Duration(len(config.Terms)))
    
    // Memory usage estimation
    fmt.Printf("Raw content size: %d bytes\n", len(config.Raw))
    fmt.Printf("Parsed terms: %d\n", len(config.Terms))
}

func generateLargeConfig(numDeps int) string {
    var builder strings.Builder
    
    builder.WriteString("{erl_opts, [debug_info]}.\n")
    builder.WriteString("{deps, [\n")
    
    for i := 0; i < numDeps; i++ {
        if i > 0 {
            builder.WriteString(",\n")
        }
        builder.WriteString(fmt.Sprintf("    {dep%d, \"1.0.%d\"}", i, i))
    }
    
    builder.WriteString("\n]}.\n")
    return builder.String()
}
```

## Next Steps

Now that you understand basic parsing:

1. **[Configuration Access](./config-access)** - Learn to use helper methods
2. **[Pretty Printing](./pretty-printing)** - Format configurations beautifully
3. **[Term Comparison](./comparison)** - Compare configurations
4. **[Complex Analysis](./complex-analysis)** - Advanced parsing scenarios

## Complete Example

Here's a complete example that demonstrates all basic parsing concepts:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Check if config file exists
    configFile := "rebar.config"
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        // Create a sample config if it doesn't exist
        createSampleConfig(configFile)
    }
    
    // Parse the configuration
    config, err := parser.ParseFile(configFile)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Display results
    fmt.Printf("✓ Successfully parsed %s\n", configFile)
    fmt.Printf("  Terms: %d\n", len(config.Terms))
    fmt.Printf("  Raw size: %d bytes\n", len(config.Raw))
    
    // Analyze each term
    fmt.Println("\nTerm analysis:")
    for i, term := range config.Terms {
        fmt.Printf("  [%d] %T: %s\n", i+1, term, term.String())
    }
    
    fmt.Println("\nFormatted output:")
    fmt.Println(config.Format(2))
}

func createSampleConfig(filename string) {
    content := `{erl_opts, [debug_info, warnings_as_errors]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
{app_name, sample_app}.
`
    
    err := os.WriteFile(filename, []byte(content), 0644)
    if err != nil {
        log.Fatalf("Failed to create sample config: %v", err)
    }
    
    fmt.Printf("Created sample config: %s\n", filename)
}
```
