---
layout: home

hero:
  name: "Erlang Rebar 配置解析器"
  text: "用于解析 Erlang rebar 配置文件的 Go 库"
  tagline: "轻松解析、访问和格式化 Erlang rebar.config 文件"
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/getting-started
    - theme: alt
      text: API 参考
      link: /zh/api/
    - theme: alt
      text: 示例
      link: /zh/examples/
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/scagogogo/erlang-rebar-config-parser

features:
  - icon: 🚀
    title: 易于使用
    details: 简单的 API 用于将 rebar.config 文件解析为结构化的 Go 对象，提供全面的辅助方法。
  
  - icon: 🔧
    title: 功能完整
    details: 支持所有常见的 Erlang 术语类型，包括元组、列表、原子、字符串、数字和嵌套结构。
  
  - icon: 📝
    title: 美化输出
    details: 格式化和美化输出 rebar 配置文件，支持可配置的缩进以提高可读性。
  
  - icon: ⚡
    title: 高性能
    details: 高效解析，98% 测试覆盖率和全面的错误处理，适用于生产环境。
  
  - icon: 🌍
    title: 多语言支持
    details: 提供中英文完整文档和示例。
  
  - icon: 🔍
    title: 术语比较
    details: 内置比较功能，用于检查不同 Erlang 术语和配置之间的相等性。
---

## 快速示例

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 解析 rebar.config 文件
    config, err := parser.ParseFile("path/to/rebar.config")
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 获取并打印依赖项
    deps, ok := config.GetDeps()
    if ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("找到 %d 个依赖项\n", len(depsList.Elements))
            
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                        fmt.Printf("- 依赖项: %s\n", atom.Value)
                    }
                }
            }
        }
    }
    
    // 格式化并打印配置，使用美观的缩进
    fmt.Println("\n格式化的配置:")
    fmt.Println(config.Format(2))
}
```

## 安装

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

## 特性

- **解析 rebar.config 文件**为结构化的 Go 对象
- **支持所有常见的 Erlang 术语类型**（元组、列表、原子、字符串、数字）
- **辅助方法**轻松访问常见配置元素
- **完全支持嵌套数据结构**
- **正确处理注释和空白字符**
- **美化输出**支持可配置缩进
- **比较功能**检查术语相等性
- **持续集成**通过 GitHub Actions
- **全面的文档**提供中英文示例
- **98% 测试覆盖率**包含全面的边缘情况测试

## 支持的 Erlang 术语类型

| Erlang 类型 | 示例 | Go 表示 |
|-------------|------|---------|
| 原子 | `atom_name`, `'quoted-atom'` | `Atom{Value: "atom_name", IsQuoted: false}` |
| 字符串 | `"hello world"` | `String{Value: "hello world"}` |
| 整数 | `123`, `-42` | `Integer{Value: 123}` |
| 浮点数 | `3.14`, `-1.5e-3` | `Float{Value: 3.14}` |
| 元组 | `{key, value}` | `Tuple{Elements: []Term{...}}` |
| 列表 | `[1, 2, 3]` | `List{Elements: []Term{...}}` |
