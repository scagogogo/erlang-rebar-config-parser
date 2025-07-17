# Makefile for Erlang Rebar Config Parser

.PHONY: help build test test-verbose test-coverage clean fmt lint vet mod-tidy mod-verify docs docs-dev docs-build docs-preview install-tools check-gitignore examples all

# Default target
help: ## Show this help message
	@echo "Erlang Rebar Config Parser - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make test              # Run all tests"
	@echo "  make test-coverage     # Run tests with coverage"
	@echo "  make docs-dev          # Start documentation development server"
	@echo "  make all               # Run full build pipeline"

# Build targets
build: ## Build the project
	@echo "🔨 Building project..."
	go build ./...

test: ## Run tests
	@echo "🧪 Running tests..."
	go test ./...

test-verbose: ## Run tests with verbose output
	@echo "🧪 Running tests (verbose)..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	@echo "🧪 Running tests with race detection..."
	go test -race ./...

# Code quality targets
fmt: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

# Module management
mod-tidy: ## Tidy go modules
	@echo "📦 Tidying modules..."
	go mod tidy

mod-verify: ## Verify go modules
	@echo "🔍 Verifying modules..."
	go mod verify

mod-download: ## Download go modules
	@echo "📥 Downloading modules..."
	go mod download

# Documentation targets
docs: docs-dev ## Start documentation development server (alias)

docs-dev: ## Start documentation development server
	@echo "📚 Starting documentation development server..."
	@cd docs && npm run docs:dev

docs-build: ## Build documentation for production
	@echo "📚 Building documentation..."
	@cd docs && npm run docs:build

docs-preview: ## Preview built documentation
	@echo "📚 Previewing documentation..."
	@cd docs && npm run docs:preview

docs-install: ## Install documentation dependencies
	@echo "📚 Installing documentation dependencies..."
	@cd docs && npm install

# Utility targets
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	go clean ./...
	rm -f coverage.out coverage.html coverage_*.html
	rm -f *.test *.prof *.pprof
	rm -f prettyprint *.formatted.config *.parsed.config
	@if [ -d "docs/.vitepress/dist" ]; then rm -rf docs/.vitepress/dist; fi
	@if [ -d "docs/.vitepress/cache" ]; then rm -rf docs/.vitepress/cache; fi

install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

check-gitignore: ## Test .gitignore effectiveness
	@echo "🧪 Testing .gitignore..."
	@./scripts/test-gitignore.sh

# Example targets
examples: ## Build example programs
	@echo "🎯 Building examples..."
	@cd examples && go build -o prettyprint ./prettyprint/
	@echo "✅ Examples built successfully"

examples-clean: ## Clean example binaries
	@echo "🧹 Cleaning example binaries..."
	@cd examples && rm -f prettyprint

# Comprehensive targets
check: fmt vet lint test ## Run all checks (format, vet, lint, test)

ci: mod-verify check test-race test-coverage ## Run CI pipeline

all: clean mod-tidy check examples docs-build ## Run full build pipeline

# Development workflow
dev-setup: install-tools docs-install ## Set up development environment
	@echo "🚀 Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  make test              # Run tests"
	@echo "  make docs-dev          # Start documentation server"
	@echo "  make examples          # Build examples"

# Release preparation
pre-release: clean mod-tidy check test-coverage docs-build examples ## Prepare for release
	@echo "🚀 Pre-release checks complete!"
	@echo ""
	@echo "Ready for release. Don't forget to:"
	@echo "  1. Update version in README.md"
	@echo "  2. Update CHANGELOG.md"
	@echo "  3. Create and push git tag"

# Quick development commands
quick-test: ## Quick test (no coverage)
	@go test ./pkg/parser

quick-check: fmt vet quick-test ## Quick check (format, vet, basic test)

# Benchmarks
bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...

bench-cpu: ## Run CPU profiling benchmark
	@echo "⚡ Running CPU profiling benchmark..."
	go test -bench=. -cpuprofile=cpu.prof ./...

bench-mem: ## Run memory profiling benchmark
	@echo "⚡ Running memory profiling benchmark..."
	go test -bench=. -memprofile=mem.prof ./...

# Help for specific areas
help-docs: ## Show documentation commands
	@echo "📚 Documentation Commands:"
	@echo "  make docs-dev          # Start development server (http://localhost:5173)"
	@echo "  make docs-build        # Build for production"
	@echo "  make docs-preview      # Preview production build"
	@echo "  make docs-install      # Install dependencies"

help-test: ## Show testing commands
	@echo "🧪 Testing Commands:"
	@echo "  make test              # Run all tests"
	@echo "  make test-verbose      # Run tests with verbose output"
	@echo "  make test-coverage     # Run tests with coverage report"
	@echo "  make test-race         # Run tests with race detection"
	@echo "  make bench             # Run benchmarks"

help-dev: ## Show development commands
	@echo "🛠️  Development Commands:"
	@echo "  make dev-setup         # Set up development environment"
	@echo "  make quick-check       # Quick format, vet, and test"
	@echo "  make check             # Full check (format, vet, lint, test)"
	@echo "  make examples          # Build example programs"
