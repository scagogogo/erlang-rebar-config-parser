# 示例

本节提供在实际场景中使用 Erlang Rebar 配置解析器的实用示例。每个示例都包含完整的、可运行的代码和解释。

## 概述

示例按复杂性和用例组织：

- **[基本解析](./basic-parsing)** - 简单解析和术语访问
- **[配置访问](./config-access)** - 使用辅助方法访问常见配置
- **[美化输出](./pretty-printing)** - 格式化和显示配置
- **[术语比较](./comparison)** - 比较配置和术语
- **[复杂分析](./complex-analysis)** - 高级解析和分析场景

## 快速示例

### 解析和显示

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
        log.Fatal(err)
    }
    
    fmt.Printf("配置有 %d 个项\n", len(config.Terms))
    fmt.Println(config.Format(2))
}
```

### 提取依赖项

```go
func extractDependencies(configPath string) ([]string, error) {
    config, err := parser.ParseFile(configPath)
    if err != nil {
        return nil, err
    }
    
    var deps []string
    if depsTerms, ok := config.GetDeps(); ok && len(depsTerms) > 0 {
        if depsList, ok := depsTerms[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                        deps = append(deps, atom.Value)
                    }
                }
            }
        }
    }
    
    return deps, nil
}
```

### 配置验证

```go
func validateConfig(config *parser.RebarConfig) []string {
    var warnings []string
    
    // 检查必需的部分
    if _, ok := config.GetDeps(); !ok {
        warnings = append(warnings, "未定义依赖项")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        warnings = append(warnings, "未定义 Erlang 选项")
    }
    
    if _, ok := config.GetAppName(); !ok {
        warnings = append(warnings, "未定义应用程序名称")
    }
    
    return warnings
}
```

## 常见模式

### 错误处理

```go
func parseWithErrorHandling(path string) (*parser.RebarConfig, error) {
    config, err := parser.ParseFile(path)
    if err != nil {
        // 为不同错误类型提供上下文
        if strings.Contains(err.Error(), "no such file") {
            return nil, fmt.Errorf("配置文件 '%s' 未找到", path)
        } else if strings.Contains(err.Error(), "syntax error") {
            return nil, fmt.Errorf("'%s' 中的语法无效: %w", path, err)
        }
        return nil, fmt.Errorf("解析 '%s' 失败: %w", path, err)
    }
    return config, nil
}
```

### 安全类型检查

```go
func safelyAccessTuple(term parser.Term, minElements int) (parser.Tuple, bool) {
    if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= minElements {
        return tuple, true
    }
    return parser.Tuple{}, false
}

func getAtomValue(term parser.Term) (string, bool) {
    if atom, ok := term.(parser.Atom); ok {
        return atom.Value, true
    }
    return "", false
}
```

### 处理集合

```go
func processErlangList(list parser.List, processor func(parser.Term)) {
    for _, element := range list.Elements {
        processor(element)
    }
}

func findInList(list parser.List, predicate func(parser.Term) bool) (parser.Term, bool) {
    for _, element := range list.Elements {
        if predicate(element) {
            return element, true
        }
    }
    return nil, false
}
```

## 实际用例

### 1. 依赖项分析器

分析项目依赖项及其版本：

```go
type Dependency struct {
    Name    string
    Version string
    Source  string
}

func analyzeDependencies(configPath string) ([]Dependency, error) {
    config, err := parser.ParseFile(configPath)
    if err != nil {
        return nil, err
    }
    
    var deps []Dependency
    if depsTerms, ok := config.GetDeps(); ok && len(depsTerms) > 0 {
        if depsList, ok := depsTerms[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if d := parseDependency(dep); d != nil {
                    deps = append(deps, *d)
                }
            }
        }
    }
    
    return deps, nil
}
```

### 2. 配置合并器

合并多个 rebar.config 文件：

```go
func mergeConfigs(paths ...string) (*parser.RebarConfig, error) {
    var allTerms []parser.Term
    
    for _, path := range paths {
        config, err := parser.ParseFile(path)
        if err != nil {
            return nil, fmt.Errorf("解析 %s 失败: %w", path, err)
        }
        allTerms = append(allTerms, config.Terms...)
    }
    
    return &parser.RebarConfig{
        Terms: allTerms,
        Raw:   "", // 合并的原始内容会很复杂
    }, nil
}
```

### 3. 配置生成器

以编程方式生成 rebar.config：

```go
func generateConfig(appName string, deps []Dependency) string {
    config := &parser.RebarConfig{
        Terms: []parser.Term{
            parser.Tuple{
                Elements: []parser.Term{
                    parser.Atom{Value: "app_name"},
                    parser.Atom{Value: appName},
                },
            },
            generateDepsConfig(deps),
            generateErlOptsConfig(),
        },
    }
    
    return config.Format(2)
}
```

## 测试示例

### 使用解析器进行单元测试

```go
func TestConfigParsing(t *testing.T) {
    configContent := `
    {erl_opts, [debug_info]}.
    {deps, [{cowboy, "2.9.0"}]}.
    `
    
    config, err := parser.Parse(configContent)
    if err != nil {
        t.Fatalf("解析失败: %v", err)
    }
    
    // 测试 erl_opts
    opts, ok := config.GetErlOpts()
    if !ok {
        t.Error("未找到 erl_opts")
    }
    
    // 测试 deps
    deps, ok := config.GetDeps()
    if !ok {
        t.Error("未找到 deps")
    }
    
    // 验证结构
    if len(config.Terms) != 2 {
        t.Errorf("期望 2 个项，得到 %d", len(config.Terms))
    }
}
```

### 基准测试

```go
func BenchmarkParsing(b *testing.B) {
    configContent := generateLargeConfig() // 您的测试数据
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := parser.Parse(configContent)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 下一步

探索详细示例：

1. **[基本解析](./basic-parsing)** - 从简单解析示例开始
2. **[配置访问](./config-access)** - 学习访问特定配置
3. **[美化输出](./pretty-printing)** - 美观地格式化配置
4. **[术语比较](./comparison)** - 比较和验证配置
5. **[复杂分析](./complex-analysis)** - 高级用例和模式

每个示例页面包括：
- 完整的、可运行的代码
- 逐步解释
- 常见陷阱及如何避免
- 性能考虑
- 测试策略
