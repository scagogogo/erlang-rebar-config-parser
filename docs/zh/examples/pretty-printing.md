# 美化输出

此示例演示如何使用库的格式化功能来格式化和美化输出 Erlang rebar 配置。

## 概述

库提供了一个 `Format()` 方法，可以使用可配置的缩进美化输出配置，使其更易读且格式正确。

## 基本格式化

### 简单示例

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 解析紧凑配置
    compactConfig := `{erl_opts,[debug_info,warnings_as_errors]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"}]}.`
    
    config, err := parser.Parse(compactConfig)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Println("原始（紧凑）:")
    fmt.Println(compactConfig)
    
    fmt.Println("\n使用 2 个空格缩进格式化:")
    fmt.Println(config.Format(2))
    
    fmt.Println("\n使用 4 个空格缩进格式化:")
    fmt.Println(config.Format(4))
}
```

### 预期输出

```
原始（紧凑）:
{erl_opts,[debug_info,warnings_as_errors]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"}]}.

使用 2 个空格缩进格式化:
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
  {cowboy, "2.9.0"},
  {jsx, "3.1.0"}
]}.

使用 4 个空格缩进格式化:
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
```

## 高级格式化示例

### 复杂配置

```go
func demonstrateComplexFormatting() {
    complexConfig := `{erl_opts,[debug_info,warnings_as_errors,{parse_transform,lager_transform}]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"},{lager,"3.9.2"}]}.{profiles,[{dev,[{deps,[{sync,"0.1.3"}]}]},{test,[{deps,[{proper,"1.3.0"},{meck,"0.9.0"}]}]}]}.{relx,[{release,{my_app,"0.1.0"},[my_app,sasl]},{dev_mode,true},{include_erts,false}]}.`
    
    config, err := parser.Parse(complexConfig)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Println("=== 复杂配置格式化 ===")
    fmt.Println("\n原始（全部在一行）:")
    fmt.Println(complexConfig)
    
    fmt.Println("\n美观格式化:")
    fmt.Println(config.Format(2))
}
```

## 比较不同缩进样式

### 比较缩进级别

```go
func compareIndentationStyles(config *parser.RebarConfig) {
    indentLevels := []int{2, 4, 8}
    
    for _, indent := range indentLevels {
        fmt.Printf("=== %d 个空格缩进 ===\n", indent)
        fmt.Println(config.Format(indent))
        fmt.Println()
    }
}
```

### 自定义格式化函数

```go
func formatWithCustomStyle(config *parser.RebarConfig) {
    // 您可以为不同样式创建包装函数
    
    // 紧凑样式（最小间距）
    fmt.Println("紧凑样式:")
    compact := config.Format(0) // 无缩进
    fmt.Println(compact)
    
    // 标准样式（2 个空格）
    fmt.Println("\n标准样式:")
    standard := config.Format(2)
    fmt.Println(standard)
    
    // 宽样式（4 个空格）
    fmt.Println("\n宽样式:")
    wide := config.Format(4)
    fmt.Println(wide)
    
    // 超宽样式（8 个空格）
    fmt.Println("\n超宽样式:")
    extraWide := config.Format(8)
    fmt.Println(extraWide)
}
```

## 实际用例

### 配置文件清理

```go
func cleanupConfigFile(inputPath, outputPath string) error {
    // 读取并解析配置
    config, err := parser.ParseFile(inputPath)
    if err != nil {
        return fmt.Errorf("解析 %s 失败: %w", inputPath, err)
    }
    
    // 使用标准 2 个空格缩进格式化
    formatted := config.Format(2)
    
    // 写入格式化的配置
    err = os.WriteFile(outputPath, []byte(formatted), 0644)
    if err != nil {
        return fmt.Errorf("写入 %s 失败: %w", outputPath, err)
    }
    
    fmt.Printf("清理了 %s -> %s\n", inputPath, outputPath)
    return nil
}

// 用法
err := cleanupConfigFile("messy_rebar.config", "clean_rebar.config")
if err != nil {
    log.Fatal(err)
}
```

### 配置差异工具

```go
func compareConfigurations(path1, path2 string) {
    config1, err := parser.ParseFile(path1)
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", path1, err)
    }
    
    config2, err := parser.ParseFile(path2)
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", path2, err)
    }
    
    fmt.Printf("=== %s ===\n", path1)
    fmt.Println(config1.Format(2))
    
    fmt.Printf("\n=== %s ===\n", path2)
    fmt.Println(config2.Format(2))
    
    // 您可以在这里添加实际的差异逻辑
    if len(config1.Terms) != len(config2.Terms) {
        fmt.Printf("\n差异: %s 有 %d 个项，%s 有 %d 个项\n", 
            path1, len(config1.Terms), path2, len(config2.Terms))
    }
}
```

## 格式化最佳实践

### 一致格式化

```go
func formatConsistently(configs []string) {
    const standardIndent = 2
    
    for i, configStr := range configs {
        config, err := parser.Parse(configStr)
        if err != nil {
            fmt.Printf("配置 %d: 解析错误 - %v\n", i+1, err)
            continue
        }
        
        fmt.Printf("=== 配置 %d（格式化） ===\n", i+1)
        fmt.Println(config.Format(standardIndent))
        fmt.Println()
    }
}
```

### 格式化前验证

```go
func formatWithValidation(configPath string) {
    config, err := parser.ParseFile(configPath)
    if err != nil {
        log.Fatalf("解析错误: %v", err)
    }
    
    // 验证配置
    if len(config.Terms) == 0 {
        fmt.Println("警告: 配置为空")
        return
    }
    
    // 检查常见部分
    hasErlOpts := false
    hasDeps := false
    
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                switch atom.Value {
                case "erl_opts":
                    hasErlOpts = true
                case "deps":
                    hasDeps = true
                }
            }
        }
    }
    
    if !hasErlOpts {
        fmt.Println("注意: 未找到 erl_opts")
    }
    if !hasDeps {
        fmt.Println("注意: 未找到 deps")
    }
    
    // 格式化并显示
    fmt.Println("格式化的配置:")
    fmt.Println(config.Format(2))
}
```

## 完整示例

这是一个演示各种格式化场景的完整示例：

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 示例混乱配置
    messyConfig := `{erl_opts,[debug_info,warnings_as_errors,{parse_transform,lager_transform}]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"},{lager,"3.9.2"}]}.{profiles,[{dev,[{deps,[{sync,"0.1.3"}]}]},{test,[{deps,[{proper,"1.3.0"}]}]}]}.`
    
    config, err := parser.Parse(messyConfig)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Println("=== 美化输出演示 ===")
    
    fmt.Println("\n1. 原始（混乱）:")
    fmt.Println(messyConfig)
    
    fmt.Println("\n2. 格式化（2 个空格缩进）:")
    fmt.Println(config.Format(2))
    
    fmt.Println("\n3. 格式化（4 个空格缩进）:")
    fmt.Println(config.Format(4))
    
    // 保存格式化版本
    formatted := config.Format(2)
    err = os.WriteFile("formatted_rebar.config", []byte(formatted), 0644)
    if err != nil {
        log.Printf("保存格式化配置失败: %v", err)
    } else {
        fmt.Println("\n4. 已将格式化版本保存到 'formatted_rebar.config'")
    }
    
    // 演示不同用例
    demonstrateUseCases(config)
}

func demonstrateUseCases(config *parser.RebarConfig) {
    fmt.Println("\n=== 用例演示 ===")
    
    // 用例 1: 配置摘要
    fmt.Println("\n--- 配置摘要 ---")
    fmt.Printf("总项数: %d\n", len(config.Terms))
    
    // 用例 2: 提取并格式化特定部分
    fmt.Println("\n--- 格式化的依赖项 ---")
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        // 创建仅包含依赖项的临时配置
        depsConfig := &parser.RebarConfig{
            Terms: []parser.Term{
                parser.Tuple{
                    Elements: []parser.Term{
                        parser.Atom{Value: "deps"},
                        deps[0],
                    },
                },
            },
        }
        fmt.Println(depsConfig.Format(2))
    }
    
    // 用例 3: 不同格式化样式
    fmt.Println("\n--- 不同样式 ---")
    styles := map[string]int{
        "紧凑": 0,
        "标准": 2,
        "宽": 4,
        "超宽": 8,
    }
    
    for name, indent := range styles {
        fmt.Printf("\n%s样式（%d 个空格）:\n", name, indent)
        // 为简洁起见，仅显示第一个项
        if len(config.Terms) > 0 {
            singleTermConfig := &parser.RebarConfig{
                Terms: []parser.Term{config.Terms[0]},
            }
            fmt.Println(singleTermConfig.Format(indent))
        }
    }
}
```

这个全面的示例展示了如何有效地使用格式化功能进行各种场景，从简单清理到复杂的配置管理任务。

## 下一步

现在您了解了美化输出：

1. **[基本解析](./basic-parsing)** - 回顾基本解析概念
2. **[配置访问](./config-access)** - 学习访问配置元素
3. **[术语比较](./comparison)** - 比较配置和术语
4. **[复杂分析](./complex-analysis)** - 高级分析技术
