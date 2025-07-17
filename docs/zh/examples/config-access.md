# 配置访问

此示例演示如何使用库提供的辅助方法轻松访问 rebar.config 文件中的常见配置元素。

## 概述

`RebarConfig` 类型提供了几个辅助方法，使访问常用配置部分变得容易，无需手动解析元组和列表。

## 可用的辅助方法

- `GetDeps()` - 获取依赖项
- `GetErlOpts()` - 获取 Erlang 编译选项
- `GetAppName()` - 获取应用程序名称
- `GetPlugins()` - 获取插件
- `GetProfilesConfig()` - 获取构建配置文件
- `GetRelxConfig()` - 获取发布配置

## 基本配置访问

### 示例配置

```erlang
{app_name, my_awesome_app}.

{erl_opts, [
    debug_info,
    warnings_as_errors,
    {parse_transform, lager_transform}
]}.

{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"},
    {lager, "3.9.2"}
]}.

{plugins, [
    rebar3_hex,
    rebar3_auto
]}.
```

### 访问配置

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
    
    // 获取应用程序名称
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("应用程序: %s\n", appName)
    } else {
        fmt.Println("未找到应用程序名称")
    }
    
    // 获取依赖项
    if deps, ok := config.GetDeps(); ok {
        fmt.Println("找到依赖项！")
        processDependencies(deps)
    } else {
        fmt.Println("未找到依赖项")
    }
    
    // 获取 Erlang 选项
    if erlOpts, ok := config.GetErlOpts(); ok {
        fmt.Println("找到 Erlang 选项！")
        processErlangOptions(erlOpts)
    } else {
        fmt.Println("未找到 Erlang 选项")
    }
    
    // 获取插件
    if plugins, ok := config.GetPlugins(); ok {
        fmt.Println("找到插件！")
        processPlugins(plugins)
    } else {
        fmt.Println("未找到插件")
    }
}

func processDependencies(deps []parser.Term) {
    if len(deps) == 0 {
        return
    }
    
    if depsList, ok := deps[0].(parser.List); ok {
        fmt.Printf("  找到 %d 个依赖项:\n", len(depsList.Elements))
        
        for _, dep := range depsList.Elements {
            if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                if name, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("    - %s", name.Value)
                    
                    // 获取版本
                    switch version := tuple.Elements[1].(type) {
                    case parser.String:
                        fmt.Printf(" (版本: %s)", version.Value)
                    case parser.Tuple:
                        fmt.Printf(" (版本规范: %s)", version.String())
                    }
                    fmt.Println()
                }
            }
        }
    }
}

func processErlangOptions(erlOpts []parser.Term) {
    if len(erlOpts) == 0 {
        return
    }
    
    if optsList, ok := erlOpts[0].(parser.List); ok {
        fmt.Printf("  找到 %d 个 Erlang 选项:\n", len(optsList.Elements))
        
        for _, opt := range optsList.Elements {
            switch o := opt.(type) {
            case parser.Atom:
                fmt.Printf("    - %s\n", o.Value)
            case parser.Tuple:
                fmt.Printf("    - %s\n", o.String())
            default:
                fmt.Printf("    - %s\n", o.String())
            }
        }
    }
}

func processPlugins(plugins []parser.Term) {
    if len(plugins) == 0 {
        return
    }
    
    if pluginsList, ok := plugins[0].(parser.List); ok {
        fmt.Printf("  找到 %d 个插件:\n", len(pluginsList.Elements))
        
        for _, plugin := range pluginsList.Elements {
            if atom, ok := plugin.(parser.Atom); ok {
                fmt.Printf("    - %s\n", atom.Value)
            }
        }
    }
}
```

## 处理配置文件

配置文件允许为不同环境（dev、test、prod）使用不同的配置。

### 示例配置文件配置

```erlang
{profiles, [
    {dev, [
        {deps, [
            {sync, "0.1.3"},
            {observer_cli, "1.7.3"}
        ]},
        {erl_opts, [debug_info]}
    ]},
    {test, [
        {deps, [
            {proper, "1.3.0"},
            {meck, "0.9.0"}
        ]},
        {erl_opts, [debug_info, export_all]}
    ]},
    {prod, [
        {erl_opts, [warnings_as_errors, no_debug_info]}
    ]}
]}.
```

### 访问配置文件配置

```go
func analyzeProfiles(config *parser.RebarConfig) {
    profiles, ok := config.GetProfilesConfig()
    if !ok {
        fmt.Println("未找到配置文件")
        return
    }
    
    if len(profiles) == 0 {
        return
    }
    
    if profilesList, ok := profiles[0].(parser.List); ok {
        fmt.Printf("找到 %d 个配置文件:\n", len(profilesList.Elements))
        
        for _, profile := range profilesList.Elements {
            if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                if name, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("\n配置文件: %s\n", name.Value)
                    
                    if configList, ok := tuple.Elements[1].(parser.List); ok {
                        analyzeProfileConfig(configList)
                    }
                }
            }
        }
    }
}

func analyzeProfileConfig(configList parser.List) {
    for _, configItem := range configList.Elements {
        if tuple, ok := configItem.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
            if key, ok := tuple.Elements[0].(parser.Atom); ok {
                fmt.Printf("  %s: %s\n", key.Value, tuple.Elements[1].String())
            }
        }
    }
}
```

## 完整示例

这是一个演示全面配置访问的完整示例：

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
    
    fmt.Println("=== 配置分析 ===")
    
    // 基本信息
    analyzeBasicInfo(config)
    
    // 依赖项
    analyzeDependencies(config)
    
    // Erlang 选项
    analyzeErlangOptions(config)
    
    // 插件
    analyzePlugins(config)
    
    // 配置文件
    analyzeProfiles(config)
    
    // 验证
    validateAndReport(config)
}

func analyzeBasicInfo(config *parser.RebarConfig) {
    fmt.Println("\n--- 基本信息 ---")
    
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("应用程序: %s\n", appName)
    } else {
        fmt.Println("应用程序: 未指定")
    }
    
    fmt.Printf("总配置项: %d\n", len(config.Terms))
}

func analyzeDependencies(config *parser.RebarConfig) {
    fmt.Println("\n--- 依赖项 ---")
    
    deps, ok := config.GetDeps()
    if !ok {
        fmt.Println("未找到依赖项")
        return
    }
    
    if len(deps) == 0 {
        fmt.Println("依赖项部分为空")
        return
    }
    
    if depsList, ok := deps[0].(parser.List); ok {
        fmt.Printf("找到 %d 个依赖项:\n", len(depsList.Elements))
        
        for i, dep := range depsList.Elements {
            fmt.Printf("%d. %s\n", i+1, dep.String())
        }
    }
}

func analyzeErlangOptions(config *parser.RebarConfig) {
    fmt.Println("\n--- Erlang 选项 ---")
    
    erlOpts, ok := config.GetErlOpts()
    if !ok {
        fmt.Println("未找到 Erlang 选项")
        return
    }
    
    if len(erlOpts) == 0 {
        fmt.Println("Erlang 选项部分为空")
        return
    }
    
    if optsList, ok := erlOpts[0].(parser.List); ok {
        fmt.Printf("找到 %d 个 Erlang 选项:\n", len(optsList.Elements))
        
        for i, opt := range optsList.Elements {
            fmt.Printf("%d. %s\n", i+1, opt.String())
        }
    }
}

func analyzePlugins(config *parser.RebarConfig) {
    fmt.Println("\n--- 插件 ---")
    
    plugins, ok := config.GetPlugins()
    if !ok {
        fmt.Println("未找到插件")
        return
    }
    
    if len(plugins) == 0 {
        fmt.Println("插件部分为空")
        return
    }
    
    if pluginsList, ok := plugins[0].(parser.List); ok {
        fmt.Printf("找到 %d 个插件:\n", len(pluginsList.Elements))
        
        for i, plugin := range pluginsList.Elements {
            fmt.Printf("%d. %s\n", i+1, plugin.String())
        }
    }
}

func validateAndReport(config *parser.RebarConfig) {
    fmt.Println("\n--- 验证 ---")
    
    warnings := validateConfiguration(config)
    if len(warnings) == 0 {
        fmt.Println("✓ 配置看起来不错！")
    } else {
        fmt.Printf("发现 %d 个警告:\n", len(warnings))
        for i, warning := range warnings {
            fmt.Printf("%d. %s\n", i+1, warning)
        }
    }
}

func validateConfiguration(config *parser.RebarConfig) []string {
    var warnings []string
    
    // 检查必需的部分
    if _, ok := config.GetAppName(); !ok {
        warnings = append(warnings, "缺少应用程序名称")
    }
    
    if _, ok := config.GetDeps(); !ok {
        warnings = append(warnings, "未定义依赖项")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        warnings = append(warnings, "未定义 Erlang 选项")
    }
    
    // 检查推荐选项
    if erlOpts, ok := config.GetErlOpts(); ok && len(erlOpts) > 0 {
        if optsList, ok := erlOpts[0].(parser.List); ok {
            hasDebugInfo := false
            for _, opt := range optsList.Elements {
                if atom, ok := opt.(parser.Atom); ok && atom.Value == "debug_info" {
                    hasDebugInfo = true
                    break
                }
            }
            if !hasDebugInfo {
                warnings = append(warnings, "推荐: 将 debug_info 添加到 erl_opts")
            }
        }
    }
    
    return warnings
}
```

此示例提供了对 rebar.config 文件的全面分析，演示了如何有效使用所有辅助方法。

## 下一步

现在您了解了配置访问：

1. **[基本解析](./basic-parsing)** - 回顾基本解析概念
2. **[美化输出](./pretty-printing)** - 学习格式化配置
3. **[术语比较](./comparison)** - 比较配置和术语
4. **[复杂分析](./complex-analysis)** - 高级分析技术
