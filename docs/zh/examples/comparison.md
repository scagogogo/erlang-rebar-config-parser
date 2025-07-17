# 术语比较

此示例演示如何使用库的比较功能来比较 Erlang 术语和配置。

## 概述

库通过所有术语类型实现的 `Compare()` 方法提供比较功能。这允许您检查不同术语和配置之间的相等性。

## 基本术语比较

### 简单比较

```go
package main

import (
    "fmt"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 创建一些术语进行比较
    atom1 := parser.Atom{Value: "debug_info", IsQuoted: false}
    atom2 := parser.Atom{Value: "debug_info", IsQuoted: true}  // 不同引号
    atom3 := parser.Atom{Value: "warnings_as_errors", IsQuoted: false}
    
    str1 := parser.String{Value: "hello"}
    str2 := parser.String{Value: "hello"}
    str3 := parser.String{Value: "world"}
    
    int1 := parser.Integer{Value: 42}
    int2 := parser.Integer{Value: 42}
    int3 := parser.Integer{Value: 24}
    
    // 原子比较（忽略 IsQuoted）
    fmt.Printf("atom1.Compare(atom2): %t\n", atom1.Compare(atom2)) // true
    fmt.Printf("atom1.Compare(atom3): %t\n", atom1.Compare(atom3)) // false
    
    // 字符串比较
    fmt.Printf("str1.Compare(str2): %t\n", str1.Compare(str2)) // true
    fmt.Printf("str1.Compare(str3): %t\n", str1.Compare(str3)) // false
    
    // 整数比较
    fmt.Printf("int1.Compare(int2): %t\n", int1.Compare(int2)) // true
    fmt.Printf("int1.Compare(int3): %t\n", int1.Compare(int3)) // false
    
    // 跨类型比较（总是 false）
    fmt.Printf("atom1.Compare(str1): %t\n", atom1.Compare(str1)) // false
    fmt.Printf("str1.Compare(int1): %t\n", str1.Compare(int1))   // false
}
```

### 复杂结构比较

```go
func demonstrateComplexComparisons() {
    // 创建相同的元组
    tuple1 := parser.Tuple{
        Elements: []parser.Term{
            parser.Atom{Value: "cowboy"},
            parser.String{Value: "2.9.0"},
        },
    }
    
    tuple2 := parser.Tuple{
        Elements: []parser.Term{
            parser.Atom{Value: "cowboy"},
            parser.String{Value: "2.9.0"},
        },
    }
    
    tuple3 := parser.Tuple{
        Elements: []parser.Term{
            parser.Atom{Value: "cowboy"},
            parser.String{Value: "2.8.0"}, // 不同版本
        },
    }
    
    fmt.Printf("tuple1.Compare(tuple2): %t\n", tuple1.Compare(tuple2)) // true
    fmt.Printf("tuple1.Compare(tuple3): %t\n", tuple1.Compare(tuple3)) // false
    
    // 创建相同的列表
    list1 := parser.List{
        Elements: []parser.Term{
            parser.Atom{Value: "debug_info"},
            parser.Atom{Value: "warnings_as_errors"},
        },
    }
    
    list2 := parser.List{
        Elements: []parser.Term{
            parser.Atom{Value: "debug_info"},
            parser.Atom{Value: "warnings_as_errors"},
        },
    }
    
    list3 := parser.List{
        Elements: []parser.Term{
            parser.Atom{Value: "debug_info"},
            // 缺少第二个元素
        },
    }
    
    fmt.Printf("list1.Compare(list2): %t\n", list1.Compare(list2)) // true
    fmt.Printf("list1.Compare(list3): %t\n", list1.Compare(list3)) // false
}
```

## 配置比较

### 比较整个配置

```go
func compareConfigurations(config1Str, config2Str string) {
    config1, err := parser.Parse(config1Str)
    if err != nil {
        fmt.Printf("解析 config1 失败: %v\n", err)
        return
    }
    
    config2, err := parser.Parse(config2Str)
    if err != nil {
        fmt.Printf("解析 config2 失败: %v\n", err)
        return
    }
    
    fmt.Println("=== 配置比较 ===")
    
    // 比较项数量
    if len(config1.Terms) != len(config2.Terms) {
        fmt.Printf("项数量不同: %d vs %d\n", 
            len(config1.Terms), len(config2.Terms))
        return
    }
    
    // 比较每个项
    allEqual := true
    for i, term1 := range config1.Terms {
        term2 := config2.Terms[i]
        if !term1.Compare(term2) {
            fmt.Printf("项 %d 不同:\n", i+1)
            fmt.Printf("  Config1: %s\n", term1.String())
            fmt.Printf("  Config2: %s\n", term2.String())
            allEqual = false
        }
    }
    
    if allEqual {
        fmt.Println("✓ 配置相同")
    } else {
        fmt.Println("✗ 配置不同")
    }
}

// 用法示例
func main() {
    config1 := `{erl_opts, [debug_info]}.{deps, [{cowboy, "2.9.0"}]}.`
    config2 := `{erl_opts, [debug_info]}.{deps, [{cowboy, "2.9.0"}]}.`
    config3 := `{erl_opts, [debug_info]}.{deps, [{cowboy, "2.8.0"}]}.`
    
    compareConfigurations(config1, config2) // 应该相同
    compareConfigurations(config1, config3) // 应该不同
}
```

### 比较特定部分

```go
func compareSpecificSections(config1, config2 *parser.RebarConfig) {
    fmt.Println("=== 逐部分比较 ===")
    
    // 比较依赖项
    deps1, ok1 := config1.GetDeps()
    deps2, ok2 := config2.GetDeps()
    
    if ok1 && ok2 {
        if len(deps1) > 0 && len(deps2) > 0 {
            if deps1[0].Compare(deps2[0]) {
                fmt.Println("✓ 依赖项相同")
            } else {
                fmt.Println("✗ 依赖项不同")
                fmt.Printf("  Config1 deps: %s\n", deps1[0].String())
                fmt.Printf("  Config2 deps: %s\n", deps2[0].String())
            }
        }
    } else if ok1 != ok2 {
        fmt.Println("✗ 一个配置有依赖项，另一个没有")
    } else {
        fmt.Println("- 两个配置都没有依赖项")
    }
    
    // 比较 Erlang 选项
    opts1, ok1 := config1.GetErlOpts()
    opts2, ok2 := config2.GetErlOpts()
    
    if ok1 && ok2 {
        if len(opts1) > 0 && len(opts2) > 0 {
            if opts1[0].Compare(opts2[0]) {
                fmt.Println("✓ Erlang 选项相同")
            } else {
                fmt.Println("✗ Erlang 选项不同")
                fmt.Printf("  Config1 erl_opts: %s\n", opts1[0].String())
                fmt.Printf("  Config2 erl_opts: %s\n", opts2[0].String())
            }
        }
    } else if ok1 != ok2 {
        fmt.Println("✗ 一个配置有 erl_opts，另一个没有")
    } else {
        fmt.Println("- 两个配置都没有 erl_opts")
    }
    
    // 比较应用程序名称
    app1, ok1 := config1.GetAppName()
    app2, ok2 := config2.GetAppName()
    
    if ok1 && ok2 {
        if app1 == app2 {
            fmt.Printf("✓ 应用程序名称相同: %s\n", app1)
        } else {
            fmt.Printf("✗ 应用程序名称不同: %s vs %s\n", app1, app2)
        }
    } else if ok1 != ok2 {
        fmt.Println("✗ 一个配置有 app_name，另一个没有")
    } else {
        fmt.Println("- 两个配置都没有 app_name")
    }
}
```

## 高级比较实用程序

### 依赖项比较

```go
func compareDependencies(config1, config2 *parser.RebarConfig) {
    fmt.Println("=== 依赖项分析 ===")
    
    deps1 := extractDependencyNames(config1)
    deps2 := extractDependencyNames(config2)
    
    // 查找共同依赖项
    common := findCommonDeps(deps1, deps2)
    if len(common) > 0 {
        fmt.Printf("共同依赖项 (%d): %v\n", len(common), common)
    }
    
    // 查找唯一依赖项
    unique1 := findUniqueDeps(deps1, deps2)
    if len(unique1) > 0 {
        fmt.Printf("仅在 config1 中 (%d): %v\n", len(unique1), unique1)
    }
    
    unique2 := findUniqueDeps(deps2, deps1)
    if len(unique2) > 0 {
        fmt.Printf("仅在 config2 中 (%d): %v\n", len(unique2), unique2)
    }
    
    // 比较共同依赖项的版本
    compareVersions(config1, config2, common)
}

func extractDependencyNames(config *parser.RebarConfig) []string {
    var names []string
    
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                        names = append(names, atom.Value)
                    }
                }
            }
        }
    }
    
    return names
}

func findCommonDeps(deps1, deps2 []string) []string {
    var common []string
    for _, dep1 := range deps1 {
        for _, dep2 := range deps2 {
            if dep1 == dep2 {
                common = append(common, dep1)
                break
            }
        }
    }
    return common
}

func findUniqueDeps(deps1, deps2 []string) []string {
    var unique []string
    for _, dep1 := range deps1 {
        found := false
        for _, dep2 := range deps2 {
            if dep1 == dep2 {
                found = true
                break
            }
        }
        if !found {
            unique = append(unique, dep1)
        }
    }
    return unique
}

func compareVersions(config1, config2 *parser.RebarConfig, commonDeps []string) {
    if len(commonDeps) == 0 {
        return
    }
    
    fmt.Println("\n--- 版本比较 ---")
    
    for _, depName := range commonDeps {
        version1 := getDependencyVersion(config1, depName)
        version2 := getDependencyVersion(config2, depName)
        
        if version1 == version2 {
            fmt.Printf("✓ %s: %s (相同)\n", depName, version1)
        } else {
            fmt.Printf("✗ %s: %s vs %s\n", depName, version1, version2)
        }
    }
}

func getDependencyVersion(config *parser.RebarConfig, depName string) string {
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == depName {
                        if str, ok := tuple.Elements[1].(parser.String); ok {
                            return str.Value
                        }
                        return tuple.Elements[1].String()
                    }
                }
            }
        }
    }
    return "unknown"
}
```

## 完整比较示例

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 用于比较的示例配置
    config1Str := `
    {app_name, my_app}.
    {erl_opts, [debug_info, warnings_as_errors]}.
    {deps, [
        {cowboy, "2.9.0"},
        {jsx, "3.1.0"}
    ]}.
    `
    
    config2Str := `
    {app_name, my_app}.
    {erl_opts, [debug_info, warnings_as_errors]}.
    {deps, [
        {cowboy, "2.8.0"},
        {jsx, "3.1.0"},
        {lager, "3.9.2"}
    ]}.
    `
    
    config1, err := parser.Parse(config1Str)
    if err != nil {
        log.Fatalf("解析 config1 失败: %v", err)
    }
    
    config2, err := parser.Parse(config2Str)
    if err != nil {
        log.Fatalf("解析 config2 失败: %v", err)
    }
    
    fmt.Println("=== 全面配置比较 ===")
    
    // 1. 基本比较
    compareConfigurations(config1Str, config2Str)
    
    fmt.Println()
    
    // 2. 逐部分比较
    compareSpecificSections(config1, config2)
    
    fmt.Println()
    
    // 3. 依赖项分析
    compareDependencies(config1, config2)
    
    // 4. 演示术语级比较
    fmt.Println("\n=== 术语级比较 ===")
    demonstrateTermComparisons()
}

func demonstrateTermComparisons() {
    // 创建各种术语进行比较
    terms := []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "debug_info", IsQuoted: true},
        parser.String{Value: "2.9.0"},
        parser.Integer{Value: 42},
        parser.Float{Value: 3.14},
        parser.Tuple{
            Elements: []parser.Term{
                parser.Atom{Value: "cowboy"},
                parser.String{Value: "2.9.0"},
            },
        },
        parser.List{
            Elements: []parser.Term{
                parser.Atom{Value: "debug_info"},
                parser.Atom{Value: "warnings_as_errors"},
            },
        },
    }
    
    // 将每个术语与自身和其他术语比较
    for i, term1 := range terms {
        for j, term2 := range terms {
            result := term1.Compare(term2)
            if i == j {
                fmt.Printf("✓ 术语 %d == 术语 %d: %t (自比较)\n", i+1, j+1, result)
            } else if result {
                fmt.Printf("✓ 术语 %d == 术语 %d: %t\n", i+1, j+1, result)
                fmt.Printf("  %s == %s\n", term1.String(), term2.String())
            }
        }
    }
}
```

这个全面的示例演示了术语和配置比较的所有方面，从简单的相等性检查到复杂的差异生成和依赖项分析。

## 下一步

现在您了解了术语比较：

1. **[基本解析](./basic-parsing)** - 回顾基本解析概念
2. **[配置访问](./config-access)** - 学习访问配置元素
3. **[美化输出](./pretty-printing)** - 学习格式化配置
4. **[复杂分析](./complex-analysis)** - 高级分析技术
