# 复杂分析

此示例演示了使用 Erlang Rebar 配置解析器进行实际应用的高级解析和分析场景。

## 概述

本节涵盖复杂用例，例如：
- 多文件配置分析
- 依赖项树分析
- 配置验证和检查
- 迁移工具
- 性能分析

## 多文件配置分析

### 项目结构分析

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
    
    // 解析主 rebar.config
    mainConfigPath := filepath.Join(projectPath, "rebar.config")
    if _, err := os.Stat(mainConfigPath); err == nil {
        config, err := parser.ParseFile(mainConfigPath)
        if err != nil {
            return nil, fmt.Errorf("解析主配置失败: %w", err)
        }
        analysis.MainConfig = config
        
        // 从主配置提取依赖项
        analysis.extractDependencies("main")
        
        // 提取配置文件
        analysis.extractProfiles()
    } else {
        analysis.Warnings = append(analysis.Warnings, "未找到主 rebar.config")
    }
    
    // 查找特定配置文件的配置
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
                        
                        // 提取版本
                        switch version := tuple.Elements[1].(type) {
                        case parser.String:
                            dependency.Version = version.Value
                        case parser.Tuple:
                            dependency.Version = version.String()
                            // 可能是 git 依赖项或版本规范
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

func (a *ProjectAnalysis) GenerateReport() string {
    var report strings.Builder
    
    report.WriteString("=== 项目分析报告 ===\n\n")
    
    // 基本信息
    if a.MainConfig != nil {
        report.WriteString(fmt.Sprintf("主配置: %d 个项\n", len(a.MainConfig.Terms)))
        
        if appName, ok := a.MainConfig.GetAppName(); ok {
            report.WriteString(fmt.Sprintf("应用程序: %s\n", appName))
        }
    }
    
    report.WriteString(fmt.Sprintf("配置文件: %d (%v)\n", len(a.Profiles), a.Profiles))
    report.WriteString(fmt.Sprintf("配置文件配置: %d\n", len(a.ProfileConfigs)))
    report.WriteString(fmt.Sprintf("依赖项: %d\n", len(a.Dependencies)))
    
    // 依赖项分析
    if len(a.Dependencies) > 0 {
        report.WriteString("\n--- 依赖项 ---\n")
        for _, dep := range a.Dependencies {
            report.WriteString(fmt.Sprintf("- %s: %s", dep.Name, dep.Version))
            if dep.Source != "" {
                report.WriteString(fmt.Sprintf(" (来源: %s)", dep.Source))
            }
            report.WriteString(fmt.Sprintf(" [配置文件: %v]\n", dep.Profiles))
        }
    }
    
    // 警告
    if len(a.Warnings) > 0 {
        report.WriteString("\n--- 警告 ---\n")
        for _, warning := range a.Warnings {
            report.WriteString(fmt.Sprintf("⚠ %s\n", warning))
        }
    }
    
    return report.String()
}
```

## 配置验证和检查

### 高级验证器

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

// 规则：检查必需部分
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
                Message: fmt.Sprintf("缺少推荐部分: %s", section),
            })
        }
    }
    
    return issues
}

// 规则：检查已弃用选项
type DeprecatedOptionsRule struct{}

func (r DeprecatedOptionsRule) Name() string {
    return "deprecated-options"
}

func (r DeprecatedOptionsRule) Validate(config *parser.RebarConfig) []ValidationIssue {
    var issues []ValidationIssue
    
    deprecatedOptions := map[string]string{
        "no_debug_info": "使用 {debug_info, false} 代替",
        "debug_info":    "考虑使用 {debug_info, true} 以提高清晰度",
    }
    
    if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
        if optsList, ok := erlOpts[0].(parser.List); ok {
            for _, opt := range optsList.Elements {
                if atom, ok := opt.(parser.Atom); ok {
                    if suggestion, deprecated := deprecatedOptions[atom.Value]; deprecated {
                        issues = append(issues, ValidationIssue{
                            Level:   "info",
                            Rule:    r.Name(),
                            Message: fmt.Sprintf("已弃用选项 '%s': %s", atom.Value, suggestion),
                            Term:    opt,
                        })
                    }
                }
            }
        }
    }
    
    return issues
}

// 规则：检查依赖项版本
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
                            // 检查宽松版本约束
                            if strings.Contains(version.Value, "*") {
                                issues = append(issues, ValidationIssue{
                                    Level:   "warning",
                                    Rule:    r.Name(),
                                    Message: fmt.Sprintf("依赖项 '%s' 使用宽松版本约束: %s", name.Value, version.Value),
                                    Term:    dep,
                                })
                            }
                            
                            // 检查非常旧的版本（示例启发式）
                            if strings.HasPrefix(version.Value, "0.") || strings.HasPrefix(version.Value, "1.") {
                                issues = append(issues, ValidationIssue{
                                    Level:   "info",
                                    Rule:    r.Name(),
                                    Message: fmt.Sprintf("依赖项 '%s' 可能使用旧版本: %s", name.Value, version.Value),
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
    
    report.WriteString("=== 配置验证报告 ===\n\n")
    
    if len(issues) == 0 {
        report.WriteString("✓ 未发现问题\n")
        return report.String()
    }
    
    // 按级别分组
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
    
    report.WriteString(fmt.Sprintf("发现 %d 个问题: %d 个错误，%d 个警告，%d 个信息\n\n", 
        len(issues), errorCount, warningCount, infoCount))
    
    // 按级别列出问题
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

## 完整分析示例

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 示例复杂配置
    complexConfig := `
    {app_name, complex_app}.
    
    {erl_opts, [
        debug_info,
        warnings_as_errors,
        {parse_transform, lager_transform},
        no_debug_info  % 已弃用选项
    ]}.
    
    {deps, [
        {cowboy, "2.9.0"},
        {jsx, "3.*"},  % 宽松版本约束
        {lager, "1.2.0"},  % 可能是旧版本
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
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Println("=== 复杂配置分析 ===")
    
    // 1. 基本分析
    performBasicAnalysis(config)
    
    // 2. 验证
    performValidation(config)
    
    // 3. 依赖项分析
    performDependencyAnalysis(config)
    
    // 4. 配置文件分析
    performProfileAnalysis(config)
}

func performBasicAnalysis(config *parser.RebarConfig) {
    fmt.Println("\n--- 基本分析 ---")
    fmt.Printf("总项数: %d\n", len(config.Terms))
    
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("应用程序: %s\n", appName)
    }
    
    // 计算不同部分类型
    sectionCounts := make(map[string]int)
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                sectionCounts[atom.Value]++
            }
        }
    }
    
    fmt.Println("部分:")
    for section, count := range sectionCounts {
        fmt.Printf("  %s: %d\n", section, count)
    }
}

func performValidation(config *parser.RebarConfig) {
    fmt.Println("\n--- 验证 ---")
    
    validator := NewConfigValidator()
    issues := validator.ValidateConfig(config)
    
    report := validator.GenerateReport(issues)
    fmt.Print(report)
}

func performDependencyAnalysis(config *parser.RebarConfig) {
    fmt.Println("\n--- 依赖项分析 ---")
    
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("主依赖项: %d\n", len(depsList.Elements))
            
            for _, dep := range depsList.Elements {
                analyzeDependency(dep)
            }
        }
    }
    
    // 分析配置文件依赖项
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            fmt.Printf("\n配置文件依赖项:\n")
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
                fmt.Printf("版本 %s", version.Value)
            case parser.Tuple:
                if len(version.Elements) > 0 {
                    if source, ok := version.Elements[0].(parser.Atom); ok {
                        fmt.Printf("来源 %s", source.Value)
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
            fmt.Printf("  配置文件 %s: ", name.Value)
            
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
                fmt.Printf("%d 个依赖项\n", depCount)
            }
        }
    }
}

func performProfileAnalysis(config *parser.RebarConfig) {
    fmt.Println("\n--- 配置文件分析 ---")
    
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            fmt.Printf("总配置文件: %d\n", len(profilesList.Elements))
            
            for _, profile := range profilesList.Elements {
                if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if name, ok := tuple.Elements[0].(parser.Atom); ok {
                        fmt.Printf("\n配置文件: %s\n", name.Value)
                        
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

这个全面的示例演示了高级分析功能，包括项目结构分析、依赖项图构建、配置验证和详细报告。

## 下一步

现在您了解了复杂分析：

1. **[基本解析](./basic-parsing)** - 回顾基本解析概念
2. **[配置访问](./config-access)** - 学习访问配置元素
3. **[美化输出](./pretty-printing)** - 学习格式化配置
4. **[术语比较](./comparison)** - 比较配置和术语
5. **[API 参考](../api/)** - 查看完整的 API 文档
