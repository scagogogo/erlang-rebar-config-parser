# 基本用法

本指南涵盖了 Erlang Rebar 配置解析器的基本使用模式。阅读后，您将了解如何解析配置、访问常见元素以及处理不同的数据类型。

## 解析配置

### 从文件

最常见的用例是解析现有的 rebar.config 文件：

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
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Printf("成功解析了 %d 个配置项\n", len(config.Terms))
}
```

### 从字符串

当您将配置内容作为字符串时：

```go
configContent := `
{erl_opts, [debug_info, warnings_as_errors]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
`

config, err := parser.Parse(configContent)
if err != nil {
    log.Fatalf("解析配置失败: %v", err)
}
```

### 从读取器

用于从任何 io.Reader 读取（文件、HTTP 响应等）：

```go
import (
    "net/http"
    "os"
)

// 从文件
file, err := os.Open("rebar.config")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

config, err := parser.ParseReader(file)

// 从 HTTP 响应
resp, err := http.Get("https://example.com/rebar.config")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

config, err = parser.ParseReader(resp.Body)
```

## 访问配置元素

### 使用辅助方法

库为常见配置部分提供了便捷的辅助方法：

```go
config, _ := parser.ParseFile("rebar.config")

// 获取依赖项
if deps, ok := config.GetDeps(); ok {
    fmt.Println("找到依赖项！")
    // deps 包含依赖项术语
}

// 获取 Erlang 编译选项
if erlOpts, ok := config.GetErlOpts(); ok {
    fmt.Println("找到 Erlang 选项！")
    // erlOpts 包含编译选项
}

// 获取应用程序名称
if appName, ok := config.GetAppName(); ok {
    fmt.Printf("应用程序名称: %s\n", appName)
}

// 获取插件
if plugins, ok := config.GetPlugins(); ok {
    fmt.Println("找到插件！")
}

// 获取配置文件
if profiles, ok := config.GetProfilesConfig(); ok {
    fmt.Println("找到配置文件！")
}
```

### 手动术语访问

对于自定义配置部分，使用通用术语访问：

```go
// 获取任何命名术语
if term, ok := config.GetTerm("custom_config"); ok {
    fmt.Printf("自定义配置: %s\n", term.String())
}

// 获取元组元素
if elements, ok := config.GetTupleElements("my_tuple"); ok {
    fmt.Printf("元组有 %d 个元素\n", len(elements))
}
```

## 处理不同术语类型

### 类型检查和转换

处理术语时始终使用安全的类型断言：

```go
func processTerm(term parser.Term) {
    switch t := term.(type) {
    case parser.Atom:
        fmt.Printf("原子: %s", t.Value)
        if t.IsQuoted {
            fmt.Print(" (引号)")
        }
        fmt.Println()
        
    case parser.String:
        fmt.Printf("字符串: %s\n", t.Value)
        
    case parser.Integer:
        fmt.Printf("整数: %d\n", t.Value)
        
    case parser.Float:
        fmt.Printf("浮点数: %f\n", t.Value)
        
    case parser.Tuple:
        fmt.Printf("元组，包含 %d 个元素\n", len(t.Elements))
        
    case parser.List:
        fmt.Printf("列表，包含 %d 个元素\n", len(t.Elements))
        
    default:
        fmt.Printf("未知术语类型: %T\n", t)
    }
}
```

### 处理集合

#### 处理列表

```go
func processList(list parser.List) {
    fmt.Printf("处理包含 %d 个元素的列表:\n", len(list.Elements))
    
    for i, element := range list.Elements {
        fmt.Printf("  [%d]: %s\n", i, element.String())
    }
}

// 示例: 处理 erl_opts
if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
    if optsList, ok := erlOpts[0].(parser.List); ok {
        processList(optsList)
    }
}
```

#### 处理元组

```go
func processTuple(tuple parser.Tuple) {
    fmt.Printf("处理包含 %d 个元素的元组:\n", len(tuple.Elements))
    
    for i, element := range tuple.Elements {
        fmt.Printf("  [%d]: %s\n", i, element.String())
    }
}

// 示例: 处理依赖项元组
func processDependency(dep parser.Term) {
    if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
        if name, ok := tuple.Elements[0].(parser.Atom); ok {
            fmt.Printf("依赖项: %s\n", name.Value)
            
            // 版本可以是字符串或元组
            switch version := tuple.Elements[1].(type) {
            case parser.String:
                fmt.Printf("  版本: %s\n", version.Value)
            case parser.Tuple:
                fmt.Printf("  版本规范: %s\n", version.String())
            }
        }
    }
}
```

## 常见模式

### 提取依赖项

```go
func extractDependencies(config *parser.RebarConfig) []string {
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
    
    return deps
}

// 用法
deps := extractDependencies(config)
fmt.Printf("找到依赖项: %v\n", deps)
```

### 检查特定选项

```go
func hasDebugInfo(config *parser.RebarConfig) bool {
    if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
        if optsList, ok := erlOpts[0].(parser.List); ok {
            for _, opt := range optsList.Elements {
                if atom, ok := opt.(parser.Atom); ok && atom.Value == "debug_info" {
                    return true
                }
            }
        }
    }
    return false
}

if hasDebugInfo(config) {
    fmt.Println("启用了调试信息")
}
```

### 查找特定配置文件

```go
func hasProfile(config *parser.RebarConfig, profileName string) bool {
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            for _, profile := range profilesList.Elements {
                if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == profileName {
                        return true
                    }
                }
            }
        }
    }
    return false
}

if hasProfile(config, "test") {
    fmt.Println("配置了测试配置文件")
}
```

## 错误处理

### 全面错误处理

```go
func parseConfigSafely(path string) (*parser.RebarConfig, error) {
    config, err := parser.ParseFile(path)
    if err != nil {
        // 为不同错误类型提供上下文
        if strings.Contains(err.Error(), "no such file") {
            return nil, fmt.Errorf("配置文件 '%s' 未找到", path)
        } else if strings.Contains(err.Error(), "permission denied") {
            return nil, fmt.Errorf("读取 '%s' 权限被拒绝", path)
        } else if strings.Contains(err.Error(), "syntax error") {
            return nil, fmt.Errorf("'%s' 中的语法无效: %w", path, err)
        }
        return nil, fmt.Errorf("解析 '%s' 失败: %w", path, err)
    }
    return config, nil
}
```

### 解析后验证

```go
func validateConfig(config *parser.RebarConfig) error {
    // 检查必需的部分
    if _, ok := config.GetDeps(); !ok {
        return fmt.Errorf("缺少必需的 'deps' 配置")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        return fmt.Errorf("缺少必需的 'erl_opts' 配置")
    }
    
    // 验证应用程序名称
    if appName, ok := config.GetAppName(); ok {
        if appName == "" {
            return fmt.Errorf("应用程序名称不能为空")
        }
    } else {
        return fmt.Errorf("缺少应用程序名称")
    }
    
    return nil
}

// 用法
config, err := parser.ParseFile("rebar.config")
if err != nil {
    log.Fatal(err)
}

if err := validateConfig(config); err != nil {
    log.Fatalf("配置验证失败: %v", err)
}
```

## 格式化和显示

### 美化输出

```go
// 使用不同缩进级别格式化
formatted2 := config.Format(2)  // 2 个空格缩进
formatted4 := config.Format(4)  // 4 个空格缩进

fmt.Println("2 个空格缩进:")
fmt.Println(formatted2)

fmt.Println("\n4 个空格缩进:")
fmt.Println(formatted4)
```

### 自定义显示函数

```go
func displayConfig(config *parser.RebarConfig) {
    fmt.Printf("配置摘要:\n")
    fmt.Printf("  总项数: %d\n", len(config.Terms))
    
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("  应用程序: %s\n", appName)
    }
    
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("  依赖项: %d\n", len(depsList.Elements))
        }
    }
    
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            fmt.Printf("  配置文件: %d\n", len(profilesList.Elements))
        }
    }
    
    fmt.Println("\n格式化的配置:")
    fmt.Println(config.Format(2))
}
```

## 下一步

现在您了解了基础知识：

1. **[高级用法](./advanced-usage)** - 了解复杂场景和最佳实践
2. **[API 参考](../api/)** - 探索完整的 API 文档
3. **[示例](../examples/)** - 查看实际示例和用例
