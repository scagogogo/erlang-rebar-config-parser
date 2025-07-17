# RebarConfig Methods

The `RebarConfig` type provides convenient methods for accessing common configuration elements and formatting the configuration.

## Configuration Access Methods

### GetTerm

```go
func (c *RebarConfig) GetTerm(name string) (Term, bool)
```

Retrieves a specific top-level term from the configuration by name.

#### Parameters

- `name` (string): The name of the term to find

#### Returns

- `Term`: The found term
- `bool`: Whether the term was found

#### Example

```go
config, _ := parser.Parse(`{erl_opts, [debug_info]}.`)

term, ok := config.GetTerm("erl_opts")
if ok {
    fmt.Printf("Found erl_opts: %s\n", term.String())
}
```

---

### GetTupleElements

```go
func (c *RebarConfig) GetTupleElements(name string) ([]Term, bool)
```

Gets the elements of a named tuple (excluding the name itself).

#### Parameters

- `name` (string): The name of the tuple

#### Returns

- `[]Term`: List of tuple elements (excluding the name)
- `bool`: Whether the tuple was found and has elements

#### Example

```go
config, _ := parser.Parse(`{deps, [{cowboy, "2.9.0"}]}.`)

elements, ok := config.GetTupleElements("deps")
if ok {
    fmt.Printf("deps has %d elements\n", len(elements))
    // elements[0] would be the list [{cowboy, "2.9.0"}]
}
```

---

### GetDeps

```go
func (c *RebarConfig) GetDeps() ([]Term, bool)
```

Retrieves the dependencies configuration.

#### Returns

- `[]Term`: List of dependency terms
- `bool`: Whether dependencies were found

#### Example

```go
config, _ := parser.Parse(`
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
`)

deps, ok := config.GetDeps()
if ok && len(deps) > 0 {
    if depsList, ok := deps[0].(parser.List); ok {
        fmt.Printf("Found %d dependencies\n", len(depsList.Elements))
        
        for _, dep := range depsList.Elements {
            if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("- %s\n", atom.Value)
                }
            }
        }
    }
}
```

---

### GetErlOpts

```go
func (c *RebarConfig) GetErlOpts() ([]Term, bool)
```

Retrieves the Erlang compilation options.

#### Returns

- `[]Term`: List of Erlang compilation options
- `bool`: Whether erl_opts were found

#### Example

```go
config, _ := parser.Parse(`{erl_opts, [debug_info, warnings_as_errors]}.`)

opts, ok := config.GetErlOpts()
if ok && len(opts) > 0 {
    if optsList, ok := opts[0].(parser.List); ok {
        fmt.Printf("Erlang options:\n")
        for _, opt := range optsList.Elements {
            fmt.Printf("- %s\n", opt.String())
        }
    }
}
```

---

### GetAppName

```go
func (c *RebarConfig) GetAppName() (string, bool)
```

Retrieves the application name.

#### Returns

- `string`: The application name
- `bool`: Whether the application name was found

#### Example

```go
config, _ := parser.Parse(`{app_name, "my_awesome_app"}.`)

appName, ok := config.GetAppName()
if ok {
    fmt.Printf("Application name: %s\n", appName)
}

// Also works with atom values
config2, _ := parser.Parse(`{app_name, my_app}.`)
appName2, ok := config2.GetAppName()
if ok {
    fmt.Printf("Application name: %s\n", appName2)
}
```

---

### GetPlugins

```go
func (c *RebarConfig) GetPlugins() ([]Term, bool)
```

Retrieves the plugins configuration.

#### Returns

- `[]Term`: List of plugin terms
- `bool`: Whether plugins were found

#### Example

```go
config, _ := parser.Parse(`{plugins, [rebar3_hex, rebar3_auto]}.`)

plugins, ok := config.GetPlugins()
if ok && len(plugins) > 0 {
    if pluginsList, ok := plugins[0].(parser.List); ok {
        fmt.Printf("Plugins:\n")
        for _, plugin := range pluginsList.Elements {
            if atom, ok := plugin.(parser.Atom); ok {
                fmt.Printf("- %s\n", atom.Value)
            }
        }
    }
}
```

---

### GetRelxConfig

```go
func (c *RebarConfig) GetRelxConfig() ([]Term, bool)
```

Retrieves the relx (release) configuration.

#### Returns

- `[]Term`: List of relx configuration terms
- `bool`: Whether relx config was found

#### Example

```go
config, _ := parser.Parse(`
{relx, [
    {release, {my_app, "0.1.0"}, [my_app, sasl]},
    {dev_mode, true},
    {include_erts, false}
]}.
`)

relx, ok := config.GetRelxConfig()
if ok && len(relx) > 0 {
    if relxList, ok := relx[0].(parser.List); ok {
        fmt.Printf("Relx configuration has %d options\n", len(relxList.Elements))
    }
}
```

---

### GetProfilesConfig

```go
func (c *RebarConfig) GetProfilesConfig() ([]Term, bool)
```

Retrieves the profiles configuration for different environments.

#### Returns

- `[]Term`: List of profile configuration terms
- `bool`: Whether profiles config was found

#### Example

```go
config, _ := parser.Parse(`
{profiles, [
    {dev, [{deps, [{sync, "0.1.3"}]}]},
    {test, [{deps, [{proper, "1.3.0"}]}]}
]}.
`)

profiles, ok := config.GetProfilesConfig()
if ok && len(profiles) > 0 {
    if profilesList, ok := profiles[0].(parser.List); ok {
        fmt.Printf("Found %d profiles\n", len(profilesList.Elements))
        
        for _, profile := range profilesList.Elements {
            if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("- Profile: %s\n", atom.Value)
                }
            }
        }
    }
}
```

---

## Formatting Methods

### Format

```go
func (c *RebarConfig) Format(indent int) string
```

Returns a formatted string representation of the configuration with the specified indentation.

#### Parameters

- `indent` (int): Number of spaces for indentation (typically 2 or 4)

#### Returns

- `string`: Formatted configuration string

#### Example

```go
config, _ := parser.Parse(`{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"}]}.`)

// Format with 2-space indentation
formatted := config.Format(2)
fmt.Println(formatted)
// Output:
// {deps, [
//   {cowboy, "2.9.0"},
//   {jsx, "3.1.0"}
// ]}.

// Format with 4-space indentation
formatted4 := config.Format(4)
fmt.Println(formatted4)
// Output:
// {deps, [
//     {cowboy, "2.9.0"},
//     {jsx, "3.1.0"}
// ]}.
```

---

## Usage Patterns

### Checking for Required Configuration

```go
func validateConfig(config *parser.RebarConfig) error {
    if _, ok := config.GetDeps(); !ok {
        return fmt.Errorf("missing required 'deps' configuration")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        return fmt.Errorf("missing required 'erl_opts' configuration")
    }
    
    return nil
}
```

### Extracting Specific Information

```go
func analyzeConfig(config *parser.RebarConfig) {
    // Check application name
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("Application: %s\n", appName)
    }
    
    // Count dependencies
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("Dependencies: %d\n", len(depsList.Elements))
        }
    }
    
    // Check for development profile
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            for _, profile := range profilesList.Elements {
                if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == "dev" {
                        fmt.Println("Development profile found")
                        break
                    }
                }
            }
        }
    }
}
```

### Configuration Modification and Formatting

```go
func modifyAndFormat(config *parser.RebarConfig) string {
    // Modify configuration (create new terms as needed)
    // ... modification logic ...
    
    // Format with consistent indentation
    return config.Format(4)
}
```
