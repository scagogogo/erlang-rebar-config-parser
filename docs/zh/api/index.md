# API 参考

Erlang Rebar 配置解析器提供了全面的 API 用于解析和操作 Erlang rebar 配置文件。本节记录了库中所有可用的公共函数、类型和方法。

## 包概述

主包 `github.com/scagogogo/erlang-rebar-config-parser/pkg/parser` 提供以下功能：

- **解析 rebar.config 文件** - 从各种来源（文件、字符串、io.Reader）
- **访问配置元素** - 通过便捷的辅助方法
- **格式化和美化输出** - 可自定义缩进的配置格式化
- **比较 Erlang 术语** - 类型感知的相等性比较
- **处理所有常见的 Erlang 数据类型** - 原子、字符串、数字、元组、列表
- **处理转义序列** - 字符串和引号原子中的转义字符
- **验证配置结构** - 提供详细的错误消息

## 快速导航

- **[核心函数](./core-functions)** - 主要解析函数（`ParseFile`、`Parse`、`ParseReader`）
- **[类型定义](./types)** - 数据类型定义（`RebarConfig`、`Term`、`Atom`、`String` 等）
- **[RebarConfig 方法](./rebar-config)** - 配置访问方法（`GetDeps`、`GetErlOpts`、`Format` 等）
- **[Term 接口](./term-interface)** - 术语类型和操作（`String()`、`Compare()`）

## 完整 API 概览

### 核心解析函数

| 函数 | 描述 | 输入 | 输出 |
|------|------|------|------|
| `ParseFile(path string)` | 从文件解析 rebar.config | 文件路径 | `*RebarConfig`, `error` |
| `Parse(input string)` | 从字符串解析 rebar.config | 配置字符串 | `*RebarConfig`, `error` |
| `ParseReader(r io.Reader)` | 从读取器解析 rebar.config | io.Reader | `*RebarConfig`, `error` |
| `NewParser(input string)` | 创建新的解析器实例 | 输入字符串 | `*Parser` |

### RebarConfig 方法

| 方法 | 描述 | 返回值 |
|------|------|--------|
| `GetTerm(name string)` | 根据名称获取特定术语 | `Term`, `bool` |
| `GetTupleElements(name string)` | 获取元组元素（不包括名称） | `[]Term`, `bool` |
| `GetDeps()` | 获取依赖项配置 | `[]Term`, `bool` |
| `GetErlOpts()` | 获取 Erlang 编译选项 | `[]Term`, `bool` |
| `GetAppName()` | 获取应用程序名称 | `string`, `bool` |
| `GetPlugins()` | 获取插件配置 | `[]Term`, `bool` |
| `GetRelxConfig()` | 获取 relx（发布）配置 | `[]Term`, `bool` |
| `GetProfilesConfig()` | 获取配置文件配置 | `[]Term`, `bool` |
| `Format(indent int)` | 使用缩进格式化配置 | `string` |

### 术语类型

| 类型 | 描述 | 字段/方法 |
|------|------|-----------|
| `Term` | 所有 Erlang 术语的基础接口 | `String()`, `Compare(Term)` |
| `Atom` | Erlang 原子（符号） | `Value string`, `IsQuoted bool` |
| `String` | Erlang 字符串（双引号） | `Value string` |
| `Integer` | Erlang 整数 | `Value int64` |
| `Float` | Erlang 浮点数 | `Value float64` |
| `Tuple` | Erlang 元组 `{a, b, c}` | `Elements []Term` |
| `List` | Erlang 列表 `[a, b, c]` | `Elements []Term` |

### 实用函数

| 函数 | 描述 | 用途 |
|------|------|------|
| `processEscapes(s string)` | 处理转义序列 | 处理 `\"`, `\\`, `\n`, `\r`, `\t` |
| `isDigit(ch byte)` | 检查字符是否为数字 | 字符分类 |
| `isAtomStart(ch byte)` | 检查是否为有效原子起始字符 | 原子验证 |
| `isAtomChar(ch byte)` | 检查是否为有效原子字符 | 原子验证 |
| `isSimpleTerm(term Term)` | 检查术语是否简单 | 格式化决策 |
| `allSimpleTerms(terms []Term)` | 检查所有术语是否都简单 | 格式化决策 |

## 基本使用模式

```go
import "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"

// 1. 从不同来源解析配置
config, err := parser.ParseFile("rebar.config")
if err != nil {
    log.Fatalf("解析错误: %v", err)
}

// 替代方案：从字符串解析
configStr := `{erl_opts, [debug_info]}.`
config, err = parser.Parse(configStr)

// 替代方案：从读取器解析
file, _ := os.Open("rebar.config")
defer file.Close()
config, err = parser.ParseReader(file)

// 2. 访问配置元素
if deps, ok := config.GetDeps(); ok {
    fmt.Println("找到依赖项！")
}

if appName, ok := config.GetAppName(); ok {
    fmt.Printf("应用程序: %s\n", appName)
}

// 3. 直接处理术语
for _, term := range config.Terms {
    switch t := term.(type) {
    case parser.Tuple:
        fmt.Printf("包含 %d 个元素的元组\n", len(t.Elements))
    case parser.List:
        fmt.Printf("包含 %d 个元素的列表\n", len(t.Elements))
    }
}

// 4. 格式化和显示
formatted := config.Format(2) // 2 个空格缩进
fmt.Println(formatted)
```

## 错误处理

解析器提供带有位置信息的详细错误消息：

```go
config, err := parser.Parse(`{invalid syntax`)
if err != nil {
    // 错误将包含行和列信息
    fmt.Printf("解析错误: %v\n", err)
    // 示例: "第 1 行第 15 列语法错误: 期望 '}'"
}
```

## 高级用法

### 自定义术语处理

```go
func processTerm(term parser.Term) {
    switch t := term.(type) {
    case parser.Atom:
        fmt.Printf("原子: %s (引号: %t)\n", t.Value, t.IsQuoted)
    case parser.String:
        fmt.Printf("字符串: %s\n", t.Value)
    case parser.Integer:
        fmt.Printf("整数: %d\n", t.Value)
    case parser.Float:
        fmt.Printf("浮点数: %f\n", t.Value)
    case parser.Tuple:
        fmt.Printf("包含 %d 个元素的元组:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    case parser.List:
        fmt.Printf("包含 %d 个元素的列表:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    }
}
```

### 术语比较

```go
atom1 := parser.Atom{Value: "debug_info", IsQuoted: false}
atom2 := parser.Atom{Value: "debug_info", IsQuoted: true}

// 比较忽略 IsQuoted 标志
if atom1.Compare(atom2) {
    fmt.Println("原子相等")
}

// 比较不同类型
str := parser.String{Value: "debug_info"}
if !atom1.Compare(str) {
    fmt.Println("不同类型不相等")
}
```

## 性能考虑

- **内存使用**: 大型配置被解析到内存中。对于非常大的文件，考虑流式处理方法。
- **解析速度**: 解析器针对典型的 rebar.config 文件进行了优化。复杂的嵌套结构可能需要更长时间。
- **格式化**: `Format()` 方法创建新字符串。对于大型配置，这可能使用大量内存。

## 线程安全

解析器类型**不是线程安全的**。如果需要从多个 goroutine 访问解析的配置，请使用适当的同步机制或创建单独的解析器实例。

## 错误处理

所有解析函数都返回错误作为第二个返回值。常见的错误场景包括：

- **文件未找到** - 当解析不存在的文件时
- **语法错误** - 当 Erlang 语法无效时
- **意外字符** - 当遇到不支持的语法时
- **未终止的字符串/原子** - 当引号未正确关闭时
- **无效数字** - 当数字格式不正确时

错误处理示例：

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    if strings.Contains(err.Error(), "no such file") {
        log.Fatal("配置文件未找到")
    } else if strings.Contains(err.Error(), "syntax error") {
        log.Fatalf("配置语法无效: %v", err)
    } else {
        log.Fatalf("解析配置失败: %v", err)
    }
}
```

## 类型安全

库使用 Go 的类型系统提供对 Erlang 术语的安全访问。使用类型断言来处理特定的术语类型：

```go
// 安全的类型断言
if atom, ok := term.(parser.Atom); ok {
    fmt.Println("原子值:", atom.Value)
}

// 处理嵌套结构
if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
    if list, ok := tuple.Elements[1].(parser.List); ok {
        fmt.Printf("列表有 %d 个元素\n", len(list.Elements))
    }
}
```

## 性能考虑

- **内存使用**: 解析器将整个配置加载到内存中
- **解析速度**: 相对于输入大小的线性时间复杂度
- **类型断言**: 类型检查的开销最小
- **字符串格式化**: 延迟求值 - 仅在需要时计算

## 线程安全

库是**读安全**的，但不是**写安全**的：

- 多个 goroutine 可以安全地从同一个 `RebarConfig` 实例读取
- 解析操作是独立的，可以并发执行
- 不要从多个 goroutine 修改 `RebarConfig` 或 `Term` 实例

## 兼容性

- **Go 版本**: 需要 Go 1.18 或更高版本
- **Erlang 兼容性**: 支持标准 Erlang 术语语法
- **Rebar 版本**: 兼容 rebar3 和 rebar2 配置格式
