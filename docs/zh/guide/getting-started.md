# 快速开始

欢迎使用 Erlang Rebar 配置解析器！本指南将帮助您快速上手解析和操作 Go 中的 Erlang rebar 配置文件。

## 这个库是什么？

Erlang Rebar 配置解析器是一个 Go 库，允许您：

- 将 Erlang rebar.config 文件解析为结构化的 Go 对象
- 通过便捷的辅助方法访问配置元素
- 格式化和美化输出配置
- 比较 Erlang 术语的相等性
- 处理所有常见的 Erlang 数据类型（原子、字符串、数字、元组、列表）

## 前提条件

- Go 1.18 或更高版本
- 对 Erlang 语法的基本了解（有帮助但不是必需的）
- 熟悉 Go 编程

## 安装

将库添加到您的 Go 项目：

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

## 您的第一个程序

让我们从一个解析 rebar.config 文件的简单示例开始：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 示例 rebar.config 内容
    configContent := `
    {erl_opts, [debug_info, warnings_as_errors]}.
    {deps, [
        {cowboy, "2.9.0"},
        {jsx, "3.1.0"}
    ]}.
    `
    
    // 解析配置
    config, err := parser.Parse(configContent)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Printf("成功解析了 %d 个配置项\n", len(config.Terms))
    
    // 访问依赖项
    deps, ok := config.GetDeps()
    if ok {
        fmt.Println("找到依赖项！")
    }
    
    // 美化输出配置
    fmt.Println("\n格式化的配置:")
    fmt.Println(config.Format(2))
}
```

## 理解输出

运行此程序时，您将看到：

```
成功解析了 2 个配置项
找到依赖项！

格式化的配置:
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
  {cowboy, "2.9.0"},
  {jsx, "3.1.0"}
]}.
```

## 关键概念

### 1. 术语

Erlang 中的一切都是"术语"。库将这些表示为实现 `Term` 接口的 Go 类型：

- **原子**: `debug_info`, `'quoted-atom'`
- **字符串**: `"hello world"`
- **数字**: `123`, `3.14`
- **元组**: `{key, value}`
- **列表**: `[item1, item2]`

### 2. 配置结构

rebar.config 文件由顶级术语组成，通常是元组，其中第一个元素是标识配置部分的原子：

```erlang
{erl_opts, [debug_info]}.        % Erlang 编译选项
{deps, [{cowboy, "2.9.0"}]}.     % 依赖项
{profiles, [{test, [...]}]}.     % 构建配置文件
```

### 3. 辅助方法

库提供便捷的方法来访问常见的配置部分：

```go
// 而不是手动解析元组
term, ok := config.GetTerm("deps")

// 使用辅助方法
deps, ok := config.GetDeps()
erlOpts, ok := config.GetErlOpts()
appName, ok := config.GetAppName()
```

## 常见模式

### 从不同来源解析

```go
// 从文件
config, err := parser.ParseFile("rebar.config")

// 从字符串
config, err := parser.Parse(configString)

// 从任何 io.Reader
config, err := parser.ParseReader(reader)
```

### 安全类型检查

```go
// 始终安全地使用类型断言
if atom, ok := term.(parser.Atom); ok {
    fmt.Printf("原子值: %s\n", atom.Value)
}

// 检查特定结构
if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
    // 安全访问 tuple.Elements[0] 和 tuple.Elements[1]
}
```

### 错误处理

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    // 处理不同的错误类型
    if strings.Contains(err.Error(), "no such file") {
        log.Fatal("配置文件未找到")
    } else if strings.Contains(err.Error(), "syntax error") {
        log.Fatalf("语法无效: %v", err)
    } else {
        log.Fatalf("解析错误: %v", err)
    }
}
```

## 下一步

现在您了解了基础知识，探索这些主题：

1. **[安装](./installation)** - 详细的安装说明
2. **[基本用法](./basic-usage)** - 常见使用模式和示例
3. **[高级用法](./advanced-usage)** - 复杂场景和最佳实践
4. **[API 参考](../api/)** - 完整的 API 文档
5. **[示例](../examples/)** - 实际示例和用例

## 快速参考

### 基本函数

```go
// 解析
config, err := parser.ParseFile("rebar.config")
config, err := parser.Parse(configString)
config, err := parser.ParseReader(reader)

// 访问配置
deps, ok := config.GetDeps()
erlOpts, ok := config.GetErlOpts()
appName, ok := config.GetAppName()

// 格式化
formatted := config.Format(2) // 2 个空格缩进
```

### 基本类型

```go
// 检查术语类型
switch t := term.(type) {
case parser.Atom:
    // t.Value, t.IsQuoted
case parser.String:
    // t.Value
case parser.Integer:
    // t.Value
case parser.Tuple:
    // t.Elements
case parser.List:
    // t.Elements
}
```

## 获取帮助

- **文档**: 浏览完整的 [API 参考](../api/)
- **示例**: 查看 [实际示例](../examples/)
- **问题**: 在 [GitHub](https://github.com/scagogogo/erlang-rebar-config-parser/issues) 上报告错误
- **讨论**: 在 [GitHub 讨论](https://github.com/scagogogo/erlang-rebar-config-parser/discussions) 中提问
