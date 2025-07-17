# Installation

This guide covers different ways to install and set up the Erlang Rebar Config Parser in your Go project.

## Requirements

- **Go version**: 1.18 or later
- **Operating System**: Linux, macOS, Windows
- **Architecture**: amd64, arm64

## Standard Installation

### Using go get (Recommended)

The easiest way to install the library is using `go get`:

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

This will download the latest version and add it to your `go.mod` file.

### Specific Version

To install a specific version:

```bash
go get github.com/scagogogo/erlang-rebar-config-parser@v1.0.0
```

### Latest Development Version

To get the latest development version from the main branch:

```bash
go get github.com/scagogogo/erlang-rebar-config-parser@main
```

## Project Setup

### New Project

If you're starting a new project:

```bash
# Create a new directory
mkdir my-rebar-parser
cd my-rebar-parser

# Initialize Go module
go mod init my-rebar-parser

# Install the library
go get github.com/scagogogo/erlang-rebar-config-parser

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    config, err := parser.Parse(`{erl_opts, [debug_info]}.`)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Parsed %d terms\n", len(config.Terms))
}
EOF

# Run the program
go run main.go
```

### Existing Project

For an existing Go project, simply add the import and run:

```bash
go mod tidy
```

This will automatically download the library when you build or run your project.

## Import Statement

Add the import to your Go files:

```go
import "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
```

### Common Import Patterns

```go
// Standard import
import "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"

// With alias
import rebarparser "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"

// Multiple imports
import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)
```

## Verification

### Quick Test

Create a simple test to verify the installation:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Test basic parsing
    config, err := parser.Parse(`{test, "installation"}.`)
    if err != nil {
        log.Fatalf("Installation test failed: %v", err)
    }
    
    // Test helper methods
    if term, ok := config.GetTerm("test"); ok {
        fmt.Printf("✓ Installation successful! Parsed: %s\n", term.String())
    } else {
        log.Fatal("✗ Installation test failed: could not retrieve term")
    }
    
    // Test formatting
    formatted := config.Format(2)
    fmt.Printf("✓ Formatting works! Output:\n%s", formatted)
}
```

Run the test:

```bash
go run main.go
```

Expected output:
```
✓ Installation successful! Parsed: {test, "installation"}
✓ Formatting works! Output:
{test, "installation"}.
```

### Unit Test

Create a proper unit test:

```go
// main_test.go
package main

import (
    "testing"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func TestInstallation(t *testing.T) {
    config, err := parser.Parse(`{erl_opts, [debug_info]}.`)
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    if len(config.Terms) != 1 {
        t.Errorf("Expected 1 term, got %d", len(config.Terms))
    }
    
    opts, ok := config.GetErlOpts()
    if !ok {
        t.Error("Could not get erl_opts")
    }
    
    if len(opts) == 0 {
        t.Error("erl_opts is empty")
    }
}
```

Run the test:

```bash
go test
```

## Troubleshooting

### Common Issues

#### 1. Go Version Too Old

**Error**: `go: module requires Go 1.18`

**Solution**: Update Go to version 1.18 or later:

```bash
# Check current version
go version

# Update Go (method varies by OS)
# On macOS with Homebrew:
brew install go

# On Ubuntu:
sudo snap install go --classic

# Or download from https://golang.org/dl/
```

#### 2. Module Not Found

**Error**: `cannot find module providing package`

**Solution**: Ensure you're in a Go module directory:

```bash
# Check if go.mod exists
ls go.mod

# If not, initialize module
go mod init your-project-name

# Then install
go get github.com/scagogogo/erlang-rebar-config-parser
```

#### 3. Import Path Issues

**Error**: `package github.com/scagogogo/erlang-rebar-config-parser/pkg/parser is not in GOROOT`

**Solution**: Make sure you're using Go modules (not GOPATH mode):

```bash
# Check Go environment
go env GOMOD

# Should show path to go.mod file, not empty
# If empty, you're in GOPATH mode - create go.mod:
go mod init your-project
```

#### 4. Network Issues

**Error**: `dial tcp: lookup proxy.golang.org: no such host`

**Solution**: Configure Go proxy or use direct mode:

```bash
# Use direct mode (bypasses proxy)
export GOPROXY=direct

# Or configure proxy
export GOPROXY=https://proxy.golang.org,direct

# For China users
export GOPROXY=https://goproxy.cn,direct
```

### Dependency Conflicts

If you encounter dependency conflicts:

```bash
# Clean module cache
go clean -modcache

# Update dependencies
go get -u github.com/scagogogo/erlang-rebar-config-parser

# Tidy up
go mod tidy
```

## IDE Setup

### VS Code

Install the Go extension and ensure your workspace is properly configured:

1. Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
2. Open your project folder in VS Code
3. Press `Ctrl+Shift+P` and run "Go: Install/Update Tools"
4. The extension should automatically detect the imported package

### GoLand/IntelliJ

1. Open your project in GoLand
2. The IDE should automatically download dependencies
3. If not, go to File → Sync Dependencies

### Vim/Neovim

With vim-go plugin:

```vim
" In your .vimrc
Plugin 'fatih/vim-go'

" Then run
:PluginInstall
:GoInstallBinaries
```

## Docker Setup

If you're using Docker:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## Next Steps

Now that you have the library installed:

1. **[Basic Usage](./basic-usage)** - Learn common usage patterns
2. **[API Reference](../api/)** - Explore the complete API
3. **[Examples](../examples/)** - See real-world examples

## Getting Help

If you encounter issues during installation:

- Check the [GitHub Issues](https://github.com/scagogogo/erlang-rebar-config-parser/issues)
- Ask questions in [GitHub Discussions](https://github.com/scagogogo/erlang-rebar-config-parser/discussions)
- Review the [troubleshooting section](#troubleshooting) above
