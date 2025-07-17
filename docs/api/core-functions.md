# Core Functions

The core functions provide the main entry points for parsing Erlang rebar configuration files from different sources.

## ParseFile

```go
func ParseFile(path string) (*RebarConfig, error)
```

Parses a rebar.config file from the given file path.

### Parameters

- `path` (string): The file path to the rebar.config file

### Returns

- `*RebarConfig`: The parsed configuration object
- `error`: Any error that occurred during parsing

### Example

```go
config, err := parser.ParseFile("./rebar.config")
if err != nil {
    log.Fatalf("Failed to parse config: %v", err)
}

fmt.Printf("Parsed %d configuration terms\n", len(config.Terms))
```

### Error Cases

- File does not exist
- Permission denied
- Invalid Erlang syntax in the file

---

## Parse

```go
func Parse(input string) (*RebarConfig, error)
```

Parses a rebar.config from the given string content.

### Parameters

- `input` (string): String containing Erlang configuration syntax

### Returns

- `*RebarConfig`: The parsed configuration object
- `error`: Any error that occurred during parsing

### Example

```go
configStr := `
{erl_opts, [debug_info]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
`

config, err := parser.Parse(configStr)
if err != nil {
    log.Fatalf("Failed to parse config: %v", err)
}

// Access parsed configuration
deps, ok := config.GetDeps()
if ok {
    fmt.Printf("Found dependencies: %v\n", deps)
}
```

### Error Cases

- Invalid Erlang syntax
- Unterminated strings or atoms
- Malformed numbers
- Missing required punctuation (dots, commas)

---

## ParseReader

```go
func ParseReader(r io.Reader) (*RebarConfig, error)
```

Parses a rebar.config from the given io.Reader.

### Parameters

- `r` (io.Reader): Reader containing Erlang configuration data

### Returns

- `*RebarConfig`: The parsed configuration object
- `error`: Any error that occurred during parsing

### Example

```go
// Parse from file
file, err := os.Open("rebar.config")
if err != nil {
    log.Fatalf("Failed to open file: %v", err)
}
defer file.Close()

config, err := parser.ParseReader(file)
if err != nil {
    log.Fatalf("Failed to parse config: %v", err)
}

// Parse from HTTP response
resp, err := http.Get("https://example.com/rebar.config")
if err != nil {
    log.Fatalf("Failed to fetch config: %v", err)
}
defer resp.Body.Close()

config, err = parser.ParseReader(resp.Body)
if err != nil {
    log.Fatalf("Failed to parse remote config: %v", err)
}
```

### Error Cases

- Read errors from the underlying reader
- Invalid Erlang syntax in the input
- Network errors (when reading from network sources)

---

## NewParser

```go
func NewParser(input string) *Parser
```

Creates a new parser instance for the given input string. This is primarily used internally but can be useful for advanced use cases.

### Parameters

- `input` (string): The input string to parse

### Returns

- `*Parser`: A new parser instance

### Example

```go
// Advanced usage - typically not needed
parser := parser.NewParser(`{erl_opts, [debug_info]}.`)
// Use parser methods directly...
```

---

## Usage Patterns

### Parsing from Different Sources

```go
// From file
config1, err := parser.ParseFile("rebar.config")

// From string
config2, err := parser.Parse(`{deps, []}.`)

// From any reader
var buf bytes.Buffer
buf.WriteString(`{erl_opts, [debug_info]}.`)
config3, err := parser.ParseReader(&buf)
```

### Error Handling Best Practices

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    // Check for specific error types
    if os.IsNotExist(err) {
        log.Fatal("Configuration file not found")
    } else if strings.Contains(err.Error(), "syntax error") {
        log.Fatalf("Invalid syntax in configuration: %v", err)
    } else {
        log.Fatalf("Unexpected error: %v", err)
    }
}
```

### Validation After Parsing

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    log.Fatalf("Parse error: %v", err)
}

// Validate that required sections exist
if _, ok := config.GetDeps(); !ok {
    log.Println("Warning: No dependencies found")
}

if _, ok := config.GetErlOpts(); !ok {
    log.Println("Warning: No Erlang options found")
}
```
