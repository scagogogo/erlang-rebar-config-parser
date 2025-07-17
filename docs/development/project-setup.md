# Project Setup and Configuration

This document describes the comprehensive project setup and configuration files for the Erlang Rebar Config Parser project.

## Overview

The project includes several configuration files to ensure consistent development experience across different platforms, editors, and team members.

## Configuration Files

### 1. `.gitignore` - Version Control Exclusions

**Purpose**: Prevents unwanted files from being committed to version control.

**Key Features**:
- Comprehensive Go build artifact exclusion
- IDE and editor file filtering
- OS-specific file handling
- Documentation build artifact exclusion
- Security-sensitive file protection
- Project-specific output file handling

**Sections**:
- Go language files (binaries, test files, coverage reports)
- IDE and editor files (IntelliJ, VS Code, Vim, Emacs, etc.)
- Operating system files (macOS, Windows, Linux)
- Documentation build files (VitePress, Node.js)
- Project-specific files (examples, utilities, test files)
- Security files (environment variables, keys, certificates)
- Development and testing artifacts
- CI/CD and deployment files

**Testing**: Use `make check-gitignore` to verify effectiveness.

### 2. `.gitattributes` - Git Behavior Configuration

**Purpose**: Ensures consistent Git behavior across different platforms.

**Key Features**:
- Automatic line ending normalization
- File type-specific handling
- Binary file identification
- Export exclusions for archives
- Language detection hints for GitHub
- Merge strategy specifications

**Benefits**:
- Consistent line endings across Windows/macOS/Linux
- Proper handling of binary files
- Better GitHub language statistics
- Cleaner git archives

### 3. `.editorconfig` - Editor Configuration

**Purpose**: Maintains consistent coding style across different editors.

**Key Features**:
- UTF-8 encoding enforcement
- Consistent indentation (tabs for Go, spaces for others)
- Line ending normalization
- Trailing whitespace handling
- File-type specific settings

**Supported Editors**: VS Code, IntelliJ, Vim, Emacs, Sublime Text, and many others.

### 4. `Makefile` - Development Automation

**Purpose**: Provides convenient commands for common development tasks.

**Available Commands**:

#### Build and Test
```bash
make build              # Build the project
make test               # Run tests
make test-coverage      # Run tests with coverage
make test-race          # Run tests with race detection
make bench              # Run benchmarks
```

#### Code Quality
```bash
make fmt                # Format code
make vet                # Run go vet
make lint               # Run linter (requires golangci-lint)
make check              # Run all checks
```

#### Documentation
```bash
make docs-dev           # Start development server
make docs-build         # Build for production
make docs-preview       # Preview production build
```

#### Utilities
```bash
make clean              # Clean build artifacts
make examples           # Build example programs
make check-gitignore    # Test .gitignore effectiveness
make dev-setup          # Set up development environment
```

#### Comprehensive Workflows
```bash
make all                # Full build pipeline
make ci                 # CI pipeline
make pre-release        # Release preparation
```

### 5. Test Scripts

#### `scripts/test-gitignore.sh`

**Purpose**: Validates that `.gitignore` is working correctly.

**Features**:
- Creates test files that should be ignored
- Verifies each pattern works as expected
- Cleans up after testing
- Provides detailed pass/fail reporting

**Usage**:
```bash
./scripts/test-gitignore.sh
# or
make check-gitignore
```

## Development Workflow

### Initial Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/scagogogo/erlang-rebar-config-parser.git
   cd erlang-rebar-config-parser
   ```

2. **Set up development environment**:
   ```bash
   make dev-setup
   ```

3. **Verify setup**:
   ```bash
   make quick-check
   ```

### Daily Development

1. **Before starting work**:
   ```bash
   make quick-check        # Quick verification
   ```

2. **During development**:
   ```bash
   make test               # Run tests frequently
   make fmt                # Format code
   ```

3. **Before committing**:
   ```bash
   make check              # Full check
   make check-gitignore    # Verify .gitignore
   ```

### Documentation Development

1. **Start documentation server**:
   ```bash
   make docs-dev
   ```

2. **Build documentation**:
   ```bash
   make docs-build
   ```

3. **Preview production build**:
   ```bash
   make docs-preview
   ```

## Best Practices

### Git Workflow

1. **Always check status before committing**:
   ```bash
   git status
   ```

2. **Use meaningful commit messages**:
   ```bash
   git commit -m "feat: add new parsing feature for complex tuples"
   ```

3. **Verify .gitignore effectiveness**:
   ```bash
   make check-gitignore
   ```

### Code Quality

1. **Format code before committing**:
   ```bash
   make fmt
   ```

2. **Run full checks**:
   ```bash
   make check
   ```

3. **Include tests for new features**:
   ```bash
   make test-coverage
   ```

### Documentation

1. **Keep documentation up to date**
2. **Test documentation builds**:
   ```bash
   make docs-build
   ```

3. **Use consistent formatting** (enforced by `.editorconfig`)

## Troubleshooting

### Common Issues

1. **Files still being tracked after adding to .gitignore**:
   ```bash
   git rm --cached filename
   git commit -m "Remove tracked file from repository"
   ```

2. **Line ending issues**:
   - Ensure `.gitattributes` is committed
   - Run `git add --renormalize .`

3. **Editor not respecting .editorconfig**:
   - Install EditorConfig plugin for your editor
   - Restart editor after installing plugin

### Validation Commands

```bash
# Check what files would be ignored
git ls-files --others --ignored --exclude-standard

# Check if specific file is ignored
git check-ignore -v filename

# Verify .gitignore effectiveness
make check-gitignore

# Check line endings
git ls-files --eol
```

## Maintenance

### Regular Tasks

1. **Update .gitignore** when adding new tools or file types
2. **Test .gitignore** after updates:
   ```bash
   make check-gitignore
   ```

3. **Review and update Makefile** commands as project evolves
4. **Keep documentation** current with project changes

### Adding New File Types

1. **Update .gitignore** with appropriate patterns
2. **Update .gitattributes** with file type handling
3. **Update .editorconfig** with formatting rules
4. **Test changes** with `make check-gitignore`

This comprehensive setup ensures a consistent, professional development environment that works well for both individual developers and teams.
