# Advanced Usage

This guide covers advanced usage patterns and best practices for the Erlang Rebar Config Parser, including performance optimization, error handling strategies, and complex use cases.

## Performance Optimization

### Memory Management

When working with large configurations or processing many files, consider memory usage:

```go
func processLargeConfigs(configPaths []string) error {
    for _, path := range configPaths {
        // Process one config at a time to limit memory usage
        config, err := parser.ParseFile(path)
        if err != nil {
            log.Printf("Failed to parse %s: %v", path, err)
            continue
        }
        
        // Process the configuration
        processConfig(config)
        
        // Allow garbage collection
        config = nil
        runtime.GC()
    }
    
    return nil
}
```

### Concurrent Processing

For processing multiple configurations concurrently:

```go
func processConcurrently(configPaths []string, maxWorkers int) {
    jobs := make(chan string, len(configPaths))
    results := make(chan ConfigResult, len(configPaths))
    
    // Start workers
    for w := 0; w < maxWorkers; w++ {
        go configWorker(jobs, results)
    }
    
    // Send jobs
    for _, path := range configPaths {
        jobs <- path
    }
    close(jobs)
    
    // Collect results
    for i := 0; i < len(configPaths); i++ {
        result := <-results
        if result.Error != nil {
            log.Printf("Error processing %s: %v", result.Path, result.Error)
        } else {
            log.Printf("Successfully processed %s", result.Path)
        }
    }
}

type ConfigResult struct {
    Path   string
    Config *parser.RebarConfig
    Error  error
}

func configWorker(jobs <-chan string, results chan<- ConfigResult) {
    for path := range jobs {
        config, err := parser.ParseFile(path)
        results <- ConfigResult{
            Path:   path,
            Config: config,
            Error:  err,
        }
    }
}
```

## Advanced Error Handling

### Custom Error Types

Create custom error types for better error handling:

```go
type ConfigError struct {
    Type    string
    Path    string
    Line    int
    Column  int
    Message string
    Cause   error
}

func (e *ConfigError) Error() string {
    if e.Path != "" {
        return fmt.Sprintf("%s error in %s: %s", e.Type, e.Path, e.Message)
    }
    return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

func (e *ConfigError) Unwrap() error {
    return e.Cause
}

func parseWithDetailedErrors(path string) (*parser.RebarConfig, error) {
    config, err := parser.ParseFile(path)
    if err != nil {
        // Enhance error with context
        configErr := &ConfigError{
            Type:    "Parse",
            Path:    path,
            Message: err.Error(),
            Cause:   err,
        }
        
        // Try to extract line/column information if available
        if strings.Contains(err.Error(), "line") {
            // Parse line information from error message
            // This is implementation-specific
        }
        
        return nil, configErr
    }
    
    return config, nil
}
```

### Error Recovery Strategies

Implement strategies to recover from partial parse failures:

```go
func parseWithRecovery(content string) (*parser.RebarConfig, []error) {
    var errors []error
    
    // Try to parse the entire content first
    config, err := parser.Parse(content)
    if err == nil {
        return config, nil
    }
    
    errors = append(errors, err)
    
    // If parsing fails, try to parse individual terms
    lines := strings.Split(content, "\n")
    var validTerms []parser.Term
    
    for i, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "%") {
            continue // Skip empty lines and comments
        }
        
        // Try to parse individual line as a term
        if strings.HasSuffix(line, ".") {
            termConfig, err := parser.Parse(line)
            if err != nil {
                errors = append(errors, fmt.Errorf("line %d: %w", i+1, err))
            } else if len(termConfig.Terms) > 0 {
                validTerms = append(validTerms, termConfig.Terms...)
            }
        }
    }
    
    if len(validTerms) > 0 {
        // Return partial configuration
        return &parser.RebarConfig{
            Terms: validTerms,
            Raw:   content,
        }, errors
    }
    
    return nil, errors
}
```

## Configuration Transformation

### Configuration Migration

Migrate configurations between different formats or versions:

```go
type ConfigMigrator struct {
    migrations []Migration
}

type Migration interface {
    Name() string
    Description() string
    Apply(config *parser.RebarConfig) (*parser.RebarConfig, error)
    CanApply(config *parser.RebarConfig) bool
}

// Example migration: Convert old-style dependencies to new format
type DependencyFormatMigration struct{}

func (m DependencyFormatMigration) Name() string {
    return "dependency-format-v2"
}

func (m DependencyFormatMigration) Description() string {
    return "Convert old-style dependency format to new format"
}

func (m DependencyFormatMigration) CanApply(config *parser.RebarConfig) bool {
    // Check if config uses old format
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                // Check for old format patterns
                if atom, ok := dep.(parser.Atom); ok {
                    // Old format: just atom names
                    _ = atom
                    return true
                }
            }
        }
    }
    return false
}

func (m DependencyFormatMigration) Apply(config *parser.RebarConfig) (*parser.RebarConfig, error) {
    newTerms := make([]parser.Term, len(config.Terms))
    copy(newTerms, config.Terms)
    
    // Find and transform dependency terms
    for i, term := range newTerms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == "deps" {
                newDeps, err := m.transformDependencies(tuple.Elements[1])
                if err != nil {
                    return nil, err
                }
                
                newTerms[i] = parser.Tuple{
                    Elements: []parser.Term{
                        tuple.Elements[0],
                        newDeps,
                    },
                }
            }
        }
    }
    
    return &parser.RebarConfig{
        Terms: newTerms,
        Raw:   "", // Raw will be regenerated when formatted
    }, nil
}

func (m DependencyFormatMigration) transformDependencies(deps parser.Term) (parser.Term, error) {
    if depsList, ok := deps.(parser.List); ok {
        var newElements []parser.Term
        
        for _, dep := range depsList.Elements {
            if atom, ok := dep.(parser.Atom); ok {
                // Convert atom to {atom, "latest"} tuple
                newDep := parser.Tuple{
                    Elements: []parser.Term{
                        atom,
                        parser.String{Value: "latest"},
                    },
                }
                newElements = append(newElements, newDep)
            } else {
                // Keep existing format
                newElements = append(newElements, dep)
            }
        }
        
        return parser.List{Elements: newElements}, nil
    }
    
    return deps, nil
}

func (cm *ConfigMigrator) Migrate(config *parser.RebarConfig) (*parser.RebarConfig, error) {
    current := config
    
    for _, migration := range cm.migrations {
        if migration.CanApply(current) {
            log.Printf("Applying migration: %s", migration.Name())
            
            migrated, err := migration.Apply(current)
            if err != nil {
                return nil, fmt.Errorf("migration %s failed: %w", migration.Name(), err)
            }
            
            current = migrated
        }
    }
    
    return current, nil
}
```

### Configuration Merging

Merge multiple configurations intelligently:

```go
type ConfigMerger struct {
    strategy MergeStrategy
}

type MergeStrategy interface {
    Merge(base, overlay *parser.RebarConfig) (*parser.RebarConfig, error)
}

// Strategy: Override - overlay completely replaces base sections
type OverrideStrategy struct{}

func (s OverrideStrategy) Merge(base, overlay *parser.RebarConfig) (*parser.RebarConfig, error) {
    baseTerms := createTermMap(base)
    overlayTerms := createTermMap(overlay)
    
    // Start with base terms
    result := make(map[string]parser.Term)
    for key, term := range baseTerms {
        result[key] = term
    }
    
    // Override with overlay terms
    for key, term := range overlayTerms {
        result[key] = term
    }
    
    // Convert back to slice
    var terms []parser.Term
    for _, term := range result {
        terms = append(terms, term)
    }
    
    return &parser.RebarConfig{Terms: terms}, nil
}

// Strategy: Merge - intelligently merge compatible sections
type MergeStrategy struct{}

func (s MergeStrategy) Merge(base, overlay *parser.RebarConfig) (*parser.RebarConfig, error) {
    baseTerms := createTermMap(base)
    overlayTerms := createTermMap(overlay)
    
    result := make(map[string]parser.Term)
    
    // Copy base terms
    for key, term := range baseTerms {
        result[key] = term
    }
    
    // Merge overlay terms
    for key, overlayTerm := range overlayTerms {
        if baseTerm, exists := result[key]; exists {
            // Try to merge compatible terms
            merged, err := s.mergeTerms(baseTerm, overlayTerm)
            if err != nil {
                // If merging fails, use overlay term
                result[key] = overlayTerm
            } else {
                result[key] = merged
            }
        } else {
            result[key] = overlayTerm
        }
    }
    
    var terms []parser.Term
    for _, term := range result {
        terms = append(terms, term)
    }
    
    return &parser.RebarConfig{Terms: terms}, nil
}

func (s MergeStrategy) mergeTerms(base, overlay parser.Term) (parser.Term, error) {
    // Only merge if both are tuples with list values (like deps, erl_opts)
    baseTuple, baseOk := base.(parser.Tuple)
    overlayTuple, overlayOk := overlay.(parser.Tuple)
    
    if !baseOk || !overlayOk || len(baseTuple.Elements) < 2 || len(overlayTuple.Elements) < 2 {
        return overlay, nil // Use overlay as-is
    }
    
    baseList, baseListOk := baseTuple.Elements[1].(parser.List)
    overlayList, overlayListOk := overlayTuple.Elements[1].(parser.List)
    
    if !baseListOk || !overlayListOk {
        return overlay, nil // Use overlay as-is
    }
    
    // Merge lists (remove duplicates)
    merged := append([]parser.Term{}, baseList.Elements...)
    
    for _, overlayElement := range overlayList.Elements {
        if !containsTerm(merged, overlayElement) {
            merged = append(merged, overlayElement)
        }
    }
    
    return parser.Tuple{
        Elements: []parser.Term{
            baseTuple.Elements[0], // Keep the key
            parser.List{Elements: merged},
        },
    }, nil
}

func containsTerm(terms []parser.Term, target parser.Term) bool {
    for _, term := range terms {
        if term.Compare(target) {
            return true
        }
    }
    return false
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

## Plugin System

### Extensible Analysis Framework

Create a plugin system for custom analysis:

```go
type AnalysisPlugin interface {
    Name() string
    Description() string
    Analyze(config *parser.RebarConfig) (AnalysisResult, error)
}

type AnalysisResult struct {
    PluginName string
    Data       interface{}
    Issues     []Issue
}

type Issue struct {
    Severity string
    Message  string
    Location string
}

// Example plugin: Security analyzer
type SecurityAnalyzer struct{}

func (p SecurityAnalyzer) Name() string {
    return "security"
}

func (p SecurityAnalyzer) Description() string {
    return "Analyzes configuration for security issues"
}

func (p SecurityAnalyzer) Analyze(config *parser.RebarConfig) (AnalysisResult, error) {
    var issues []Issue
    
    // Check for insecure dependencies
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if name, ok := tuple.Elements[0].(parser.Atom); ok {
                        // Check for known insecure packages
                        if p.isInsecurePackage(name.Value) {
                            issues = append(issues, Issue{
                                Severity: "high",
                                Message:  fmt.Sprintf("Dependency '%s' has known security vulnerabilities", name.Value),
                                Location: "deps",
                            })
                        }
                        
                        // Check for git dependencies without pinned versions
                        if len(tuple.Elements) >= 2 {
                            if gitSpec, ok := tuple.Elements[1].(parser.Tuple); ok {
                                if len(gitSpec.Elements) > 0 {
                                    if source, ok := gitSpec.Elements[0].(parser.Atom); ok && source.Value == "git" {
                                        // Check if version is pinned
                                        if !p.hasVersionPin(gitSpec) {
                                            issues = append(issues, Issue{
                                                Severity: "medium",
                                                Message:  fmt.Sprintf("Git dependency '%s' should be pinned to a specific version", name.Value),
                                                Location: "deps",
                                            })
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    
    return AnalysisResult{
        PluginName: p.Name(),
        Data:       map[string]int{"total_issues": len(issues)},
        Issues:     issues,
    }, nil
}

func (p SecurityAnalyzer) isInsecurePackage(name string) bool {
    // This would typically check against a database of known vulnerabilities
    insecurePackages := []string{"vulnerable_package", "old_crypto_lib"}
    for _, pkg := range insecurePackages {
        if name == pkg {
            return true
        }
    }
    return false
}

func (p SecurityAnalyzer) hasVersionPin(gitSpec parser.Tuple) bool {
    // Check if git spec has a tag or specific commit
    for _, element := range gitSpec.Elements {
        if tuple, ok := element.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
            if key, ok := tuple.Elements[0].(parser.Atom); ok {
                if key.Value == "tag" || key.Value == "ref" {
                    return true
                }
            }
        }
    }
    return false
}

// Plugin manager
type AnalysisManager struct {
    plugins []AnalysisPlugin
}

func NewAnalysisManager() *AnalysisManager {
    return &AnalysisManager{
        plugins: []AnalysisPlugin{
            SecurityAnalyzer{},
            // Add more plugins here
        },
    }
}

func (am *AnalysisManager) RunAnalysis(config *parser.RebarConfig) ([]AnalysisResult, error) {
    var results []AnalysisResult
    
    for _, plugin := range am.plugins {
        result, err := plugin.Analyze(config)
        if err != nil {
            log.Printf("Plugin %s failed: %v", plugin.Name(), err)
            continue
        }
        results = append(results, result)
    }
    
    return results, nil
}

func (am *AnalysisManager) GenerateReport(results []AnalysisResult) string {
    var report strings.Builder
    
    report.WriteString("=== Analysis Report ===\n\n")
    
    totalIssues := 0
    for _, result := range results {
        totalIssues += len(result.Issues)
    }
    
    report.WriteString(fmt.Sprintf("Total issues found: %d\n\n", totalIssues))
    
    for _, result := range results {
        report.WriteString(fmt.Sprintf("--- %s ---\n", strings.ToUpper(result.PluginName)))
        
        if len(result.Issues) == 0 {
            report.WriteString("✓ No issues found\n\n")
            continue
        }
        
        for _, issue := range result.Issues {
            severity := strings.ToUpper(issue.Severity)
            report.WriteString(fmt.Sprintf("[%s] %s\n", severity, issue.Message))
            if issue.Location != "" {
                report.WriteString(fmt.Sprintf("    Location: %s\n", issue.Location))
            }
        }
        report.WriteString("\n")
    }
    
    return report.String()
}
```

## Best Practices Summary

### 1. Memory Management
- Process large files one at a time
- Use goroutines for concurrent processing
- Allow garbage collection between operations

### 2. Error Handling
- Create custom error types for better context
- Implement recovery strategies for partial failures
- Provide detailed error messages with location information

### 3. Configuration Management
- Use migration patterns for format changes
- Implement intelligent merging strategies
- Validate configurations before processing

### 4. Extensibility
- Design plugin systems for custom analysis
- Use interfaces for flexible implementations
- Provide comprehensive reporting mechanisms

### 5. Performance
- Cache parsed configurations when possible
- Use streaming for very large files
- Profile memory usage in production

This advanced guide provides the foundation for building robust, production-ready applications using the Erlang Rebar Config Parser.
