# Complex Analysis

This example demonstrates advanced parsing and analysis scenarios using the Erlang Rebar Config Parser for real-world applications.

## Overview

This section covers complex use cases such as:
- Multi-file configuration analysis
- Dependency tree analysis
- Configuration validation and linting
- Migration tools
- Performance analysis

## Multi-File Configuration Analysis

### Project Structure Analysis

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

type ProjectAnalysis struct {
    MainConfig    *parser.RebarConfig
    ProfileConfigs map[string]*parser.RebarConfig
    Dependencies  []Dependency
    Profiles      []string
    Warnings      []string
}

type Dependency struct {
    Name     string
    Version  string
    Source   string
    Profiles []string
}

func analyzeProject(projectPath string) (*ProjectAnalysis, error) {
    analysis := &ProjectAnalysis{
        ProfileConfigs: make(map[string]*parser.RebarConfig),
        Dependencies:   []Dependency{},
        Profiles:       []string{},
        Warnings:       []string{},
    }
    
    // Parse main rebar.config
    mainConfigPath := filepath.Join(projectPath, "rebar.config")
    if _, err := os.Stat(mainConfigPath); err == nil {
        config, err := parser.ParseFile(mainConfigPath)
        if err != nil {
            return nil, fmt.Errorf("failed to parse main config: %w", err)
        }
        analysis.MainConfig = config
        
        // Extract dependencies from main config
        analysis.extractDependencies("main")
        
        // Extract profiles
        analysis.extractProfiles()
    } else {
        analysis.Warnings = append(analysis.Warnings, "No main rebar.config found")
    }
    
    // Look for profile-specific configs
    configDir := filepath.Join(projectPath, "config")
    if _, err := os.Stat(configDir); err == nil {
        analysis.analyzeProfileConfigs(configDir)
    }
    
    return analysis, nil
}

func (a *ProjectAnalysis) extractDependencies(profile string) {
    if a.MainConfig == nil {
        return
    }
    
    if deps, ok := a.MainConfig.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if name, ok := tuple.Elements[0].(parser.Atom); ok {
                        dependency := Dependency{
                            Name:     name.Value,
                            Profiles: []string{profile},
                        }
                        
                        // Extract version
                        switch version := tuple.Elements[1].(type) {
                        case parser.String:
                            dependency.Version = version.Value
                        case parser.Tuple:
                            dependency.Version = version.String()
                            // Could be a git dependency or version spec
                            if len(version.Elements) > 0 {
                                if atom, ok := version.Elements[0].(parser.Atom); ok {
                                    dependency.Source = atom.Value
                                }
                            }
                        }
                        
                        a.Dependencies = append(a.Dependencies, dependency)
                    }
                }
            }
        }
    }
}

func (a *ProjectAnalysis) extractProfiles() {
    if a.MainConfig == nil {
        return
    }
    
    if profiles, ok := a.MainConfig.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            for _, profile := range profilesList.Elements {
                if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if name, ok := tuple.Elements[0].(parser.Atom); ok {
                        a.Profiles = append(a.Profiles, name.Value)
                    }
                }
            }
        }
    }
}

func (a *ProjectAnalysis) analyzeProfileConfigs(configDir string) {
    filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if strings.HasSuffix(path, ".config") && !info.IsDir() {
            profileName := strings.TrimSuffix(filepath.Base(path), ".config")
            
            config, err := parser.ParseFile(path)
            if err != nil {
                a.Warnings = append(a.Warnings, 
                    fmt.Sprintf("Failed to parse %s: %v", path, err))
                return nil
            }
            
            a.ProfileConfigs[profileName] = config
        }
        
        return nil
    })
}

func (a *ProjectAnalysis) GenerateReport() string {
    var report strings.Builder
    
    report.WriteString("=== Project Analysis Report ===\n\n")
    
    // Basic information
    if a.MainConfig != nil {
        report.WriteString(fmt.Sprintf("Main configuration: %d terms\n", len(a.MainConfig.Terms)))
        
        if appName, ok := a.MainConfig.GetAppName(); ok {
            report.WriteString(fmt.Sprintf("Application: %s\n", appName))
        }
    }
    
    report.WriteString(fmt.Sprintf("Profiles: %d (%v)\n", len(a.Profiles), a.Profiles))
    report.WriteString(fmt.Sprintf("Profile configs: %d\n", len(a.ProfileConfigs)))
    report.WriteString(fmt.Sprintf("Dependencies: %d\n", len(a.Dependencies)))
    
    // Dependencies analysis
    if len(a.Dependencies) > 0 {
        report.WriteString("\n--- Dependencies ---\n")
        for _, dep := range a.Dependencies {
            report.WriteString(fmt.Sprintf("- %s: %s", dep.Name, dep.Version))
            if dep.Source != "" {
                report.WriteString(fmt.Sprintf(" (source: %s)", dep.Source))
            }
            report.WriteString(fmt.Sprintf(" [profiles: %v]\n", dep.Profiles))
        }
    }
    
    // Warnings
    if len(a.Warnings) > 0 {
        report.WriteString("\n--- Warnings ---\n")
        for _, warning := range a.Warnings {
            report.WriteString(fmt.Sprintf("⚠ %s\n", warning))
        }
    }
    
    return report.String()
}
```

## Dependency Tree Analysis

### Dependency Graph Builder

```go
type DependencyGraph struct {
    Nodes map[string]*DependencyNode
    Edges map[string][]string
}

type DependencyNode struct {
    Name         string
    Version      string
    Dependencies []string
    Dependents   []string
}

func buildDependencyGraph(configs []*parser.RebarConfig) *DependencyGraph {
    graph := &DependencyGraph{
        Nodes: make(map[string]*DependencyNode),
        Edges: make(map[string][]string),
    }
    
    // First pass: collect all dependencies
    for _, config := range configs {
        if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
            if depsList, ok := deps[0].(parser.List); ok {
                for _, dep := range depsList.Elements {
                    if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                        if name, ok := tuple.Elements[0].(parser.Atom); ok {
                            node := &DependencyNode{
                                Name:         name.Value,
                                Dependencies: []string{},
                                Dependents:   []string{},
                            }
                            
                            if version, ok := tuple.Elements[1].(parser.String); ok {
                                node.Version = version.Value
                            }
                            
                            graph.Nodes[name.Value] = node
                        }
                    }
                }
            }
        }
    }
    
    return graph
}

func (g *DependencyGraph) FindCircularDependencies() [][]string {
    var cycles [][]string
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    
    for node := range g.Nodes {
        if !visited[node] {
            if cycle := g.dfsForCycle(node, visited, recStack, []string{}); len(cycle) > 0 {
                cycles = append(cycles, cycle)
            }
        }
    }
    
    return cycles
}

func (g *DependencyGraph) dfsForCycle(node string, visited, recStack map[string]bool, path []string) []string {
    visited[node] = true
    recStack[node] = true
    path = append(path, node)
    
    for _, neighbor := range g.Edges[node] {
        if !visited[neighbor] {
            if cycle := g.dfsForCycle(neighbor, visited, recStack, path); len(cycle) > 0 {
                return cycle
            }
        } else if recStack[neighbor] {
            // Found cycle
            cycleStart := -1
            for i, n := range path {
                if n == neighbor {
                    cycleStart = i
                    break
                }
            }
            if cycleStart >= 0 {
                return append(path[cycleStart:], neighbor)
            }
        }
    }
    
    recStack[node] = false
    return nil
}

func (g *DependencyGraph) GenerateReport() string {
    var report strings.Builder
    
    report.WriteString("=== Dependency Graph Analysis ===\n\n")
    report.WriteString(fmt.Sprintf("Total dependencies: %d\n", len(g.Nodes)))
    
    // Find root dependencies (no dependents)
    var roots []string
    for name, node := range g.Nodes {
        if len(node.Dependents) == 0 {
            roots = append(roots, name)
        }
    }
    
    if len(roots) > 0 {
        report.WriteString(fmt.Sprintf("Root dependencies: %v\n", roots))
    }
    
    // Check for circular dependencies
    cycles := g.FindCircularDependencies()
    if len(cycles) > 0 {
        report.WriteString("\n⚠ Circular dependencies found:\n")
        for i, cycle := range cycles {
            report.WriteString(fmt.Sprintf("  %d. %s\n", i+1, strings.Join(cycle, " -> ")))
        }
    } else {
        report.WriteString("\n✓ No circular dependencies found\n")
    }
    
    return report.String()
}
```

## Configuration Validation and Linting

### Advanced Validator

```go
type ConfigValidator struct {
    Rules []ValidationRule
}

type ValidationRule interface {
    Validate(config *parser.RebarConfig) []ValidationIssue
    Name() string
}

type ValidationIssue struct {
    Level   string // "error", "warning", "info"
    Rule    string
    Message string
    Term    parser.Term
}

// Rule: Check for required sections
type RequiredSectionsRule struct{}

func (r RequiredSectionsRule) Name() string {
    return "required-sections"
}

func (r RequiredSectionsRule) Validate(config *parser.RebarConfig) []ValidationIssue {
    var issues []ValidationIssue
    
    requiredSections := []string{"erl_opts", "deps"}
    
    for _, section := range requiredSections {
        if _, ok := config.GetTerm(section); !ok {
            issues = append(issues, ValidationIssue{
                Level:   "warning",
                Rule:    r.Name(),
                Message: fmt.Sprintf("Missing recommended section: %s", section),
            })
        }
    }
    
    return issues
}

// Rule: Check for deprecated options
type DeprecatedOptionsRule struct{}

func (r DeprecatedOptionsRule) Name() string {
    return "deprecated-options"
}

func (r DeprecatedOptionsRule) Validate(config *parser.RebarConfig) []ValidationIssue {
    var issues []ValidationIssue
    
    deprecatedOptions := map[string]string{
        "no_debug_info": "Use {debug_info, false} instead",
        "debug_info":    "Consider using {debug_info, true} for clarity",
    }
    
    if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
        if optsList, ok := erlOpts[0].(parser.List); ok {
            for _, opt := range optsList.Elements {
                if atom, ok := opt.(parser.Atom); ok {
                    if suggestion, deprecated := deprecatedOptions[atom.Value]; deprecated {
                        issues = append(issues, ValidationIssue{
                            Level:   "info",
                            Rule:    r.Name(),
                            Message: fmt.Sprintf("Deprecated option '%s': %s", atom.Value, suggestion),
                            Term:    opt,
                        })
                    }
                }
            }
        }
    }
    
    return issues
}

// Rule: Check dependency versions
type DependencyVersionRule struct{}

func (r DependencyVersionRule) Name() string {
    return "dependency-versions"
}

func (r DependencyVersionRule) Validate(config *parser.RebarConfig) []ValidationIssue {
    var issues []ValidationIssue
    
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if name, ok := tuple.Elements[0].(parser.Atom); ok {
                        if version, ok := tuple.Elements[1].(parser.String); ok {
                            // Check for loose version constraints
                            if strings.Contains(version.Value, "*") {
                                issues = append(issues, ValidationIssue{
                                    Level:   "warning",
                                    Rule:    r.Name(),
                                    Message: fmt.Sprintf("Dependency '%s' uses loose version constraint: %s", name.Value, version.Value),
                                    Term:    dep,
                                })
                            }
                            
                            // Check for very old versions (example heuristic)
                            if strings.HasPrefix(version.Value, "0.") || strings.HasPrefix(version.Value, "1.") {
                                issues = append(issues, ValidationIssue{
                                    Level:   "info",
                                    Rule:    r.Name(),
                                    Message: fmt.Sprintf("Dependency '%s' may be using an old version: %s", name.Value, version.Value),
                                    Term:    dep,
                                })
                            }
                        }
                    }
                }
            }
        }
    }
    
    return issues
}

func NewConfigValidator() *ConfigValidator {
    return &ConfigValidator{
        Rules: []ValidationRule{
            RequiredSectionsRule{},
            DeprecatedOptionsRule{},
            DependencyVersionRule{},
        },
    }
}

func (v *ConfigValidator) ValidateConfig(config *parser.RebarConfig) []ValidationIssue {
    var allIssues []ValidationIssue
    
    for _, rule := range v.Rules {
        issues := rule.Validate(config)
        allIssues = append(allIssues, issues...)
    }
    
    return allIssues
}

func (v *ConfigValidator) GenerateReport(issues []ValidationIssue) string {
    var report strings.Builder
    
    report.WriteString("=== Configuration Validation Report ===\n\n")
    
    if len(issues) == 0 {
        report.WriteString("✓ No issues found\n")
        return report.String()
    }
    
    // Group by level
    errorCount := 0
    warningCount := 0
    infoCount := 0
    
    for _, issue := range issues {
        switch issue.Level {
        case "error":
            errorCount++
        case "warning":
            warningCount++
        case "info":
            infoCount++
        }
    }
    
    report.WriteString(fmt.Sprintf("Found %d issues: %d errors, %d warnings, %d info\n\n", 
        len(issues), errorCount, warningCount, infoCount))
    
    // List issues by level
    levels := []string{"error", "warning", "info"}
    symbols := map[string]string{"error": "✗", "warning": "⚠", "info": "ℹ"}
    
    for _, level := range levels {
        levelIssues := filterIssuesByLevel(issues, level)
        if len(levelIssues) > 0 {
            report.WriteString(fmt.Sprintf("--- %s (%d) ---\n", strings.ToUpper(level), len(levelIssues)))
            for _, issue := range levelIssues {
                report.WriteString(fmt.Sprintf("%s [%s] %s\n", symbols[level], issue.Rule, issue.Message))
                if issue.Term != nil {
                    report.WriteString(fmt.Sprintf("    %s\n", issue.Term.String()))
                }
            }
            report.WriteString("\n")
        }
    }
    
    return report.String()
}

func filterIssuesByLevel(issues []ValidationIssue, level string) []ValidationIssue {
    var filtered []ValidationIssue
    for _, issue := range issues {
        if issue.Level == level {
            filtered = append(filtered, issue)
        }
    }
    return filtered
}
```

## Complete Analysis Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // Sample complex configuration
    complexConfig := `
    {app_name, complex_app}.
    
    {erl_opts, [
        debug_info,
        warnings_as_errors,
        {parse_transform, lager_transform},
        no_debug_info  % deprecated option
    ]}.
    
    {deps, [
        {cowboy, "2.9.0"},
        {jsx, "3.*"},  % loose version constraint
        {lager, "1.2.0"},  % potentially old version
        {custom_dep, {git, "https://github.com/user/custom_dep.git", {branch, "master"}}}
    ]}.
    
    {profiles, [
        {dev, [
            {deps, [{sync, "0.1.3"}]},
            {erl_opts, [debug_info]}
        ]},
        {test, [
            {deps, [{proper, "1.3.0"}, {meck, "0.9.0"}]},
            {erl_opts, [debug_info, export_all]}
        ]},
        {prod, [
            {erl_opts, [warnings_as_errors, no_debug_info]}
        ]}
    ]}.
    
    {relx, [
        {release, {complex_app, "0.1.0"}, [complex_app, sasl]},
        {dev_mode, true},
        {include_erts, false}
    ]}.
    `
    
    config, err := parser.Parse(complexConfig)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    fmt.Println("=== Complex Configuration Analysis ===")
    
    // 1. Basic analysis
    performBasicAnalysis(config)
    
    // 2. Validation
    performValidation(config)
    
    // 3. Dependency analysis
    performDependencyAnalysis(config)
    
    // 4. Profile analysis
    performProfileAnalysis(config)
}

func performBasicAnalysis(config *parser.RebarConfig) {
    fmt.Println("\n--- Basic Analysis ---")
    fmt.Printf("Total terms: %d\n", len(config.Terms))
    
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("Application: %s\n", appName)
    }
    
    // Count different section types
    sectionCounts := make(map[string]int)
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                sectionCounts[atom.Value]++
            }
        }
    }
    
    fmt.Println("Sections:")
    for section, count := range sectionCounts {
        fmt.Printf("  %s: %d\n", section, count)
    }
}

func performValidation(config *parser.RebarConfig) {
    fmt.Println("\n--- Validation ---")
    
    validator := NewConfigValidator()
    issues := validator.ValidateConfig(config)
    
    report := validator.GenerateReport(issues)
    fmt.Print(report)
}

func performDependencyAnalysis(config *parser.RebarConfig) {
    fmt.Println("\n--- Dependency Analysis ---")
    
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("Main dependencies: %d\n", len(depsList.Elements))
            
            for _, dep := range depsList.Elements {
                analyzeDependency(dep)
            }
        }
    }
    
    // Analyze profile dependencies
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            fmt.Printf("\nProfile dependencies:\n")
            for _, profile := range profilesList.Elements {
                analyzeProfileDependencies(profile)
            }
        }
    }
}

func analyzeDependency(dep parser.Term) {
    if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
        if name, ok := tuple.Elements[0].(parser.Atom); ok {
            fmt.Printf("  %s: ", name.Value)
            
            switch version := tuple.Elements[1].(type) {
            case parser.String:
                fmt.Printf("version %s", version.Value)
            case parser.Tuple:
                if len(version.Elements) > 0 {
                    if source, ok := version.Elements[0].(parser.Atom); ok {
                        fmt.Printf("source %s", source.Value)
                    }
                }
            }
            fmt.Println()
        }
    }
}

func analyzeProfileDependencies(profile parser.Term) {
    if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
        if name, ok := tuple.Elements[0].(parser.Atom); ok {
            fmt.Printf("  Profile %s: ", name.Value)
            
            if configList, ok := tuple.Elements[1].(parser.List); ok {
                depCount := 0
                for _, configItem := range configList.Elements {
                    if configTuple, ok := configItem.(parser.Tuple); ok && len(configTuple.Elements) >= 2 {
                        if key, ok := configTuple.Elements[0].(parser.Atom); ok && key.Value == "deps" {
                            if depsList, ok := configTuple.Elements[1].(parser.List); ok {
                                depCount = len(depsList.Elements)
                            }
                        }
                    }
                }
                fmt.Printf("%d dependencies\n", depCount)
            }
        }
    }
}

func performProfileAnalysis(config *parser.RebarConfig) {
    fmt.Println("\n--- Profile Analysis ---")
    
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            fmt.Printf("Total profiles: %d\n", len(profilesList.Elements))
            
            for _, profile := range profilesList.Elements {
                if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if name, ok := tuple.Elements[0].(parser.Atom); ok {
                        fmt.Printf("\nProfile: %s\n", name.Value)
                        
                        if configList, ok := tuple.Elements[1].(parser.List); ok {
                            for _, configItem := range configList.Elements {
                                if configTuple, ok := configItem.(parser.Tuple); ok && len(configTuple.Elements) >= 2 {
                                    if key, ok := configTuple.Elements[0].(parser.Atom); ok {
                                        fmt.Printf("  %s: %s\n", key.Value, configTuple.Elements[1].String())
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
```

This comprehensive example demonstrates advanced analysis capabilities including project structure analysis, dependency graph building, configuration validation, and detailed reporting.
