# Erlang Rebar 配置解析器

**中文文档** | [English](README.md) | [📖 文档网站](https://scagogogo.github.io/erlang-rebar-config-parser/)

[![Go 测试和示例](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/go.yml/badge.svg)](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/go.yml)
[![文档部署](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/docs.yml/badge.svg)](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/docs.yml)
[![GoDoc](https://godoc.org/github.com/scagogogo/erlang-rebar-config-parser?status.svg)](https://godoc.org/github.com/scagogogo/erlang-rebar-config-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/erlang-rebar-config-parser)](https://goreportcard.com/report/github.com/scagogogo/erlang-rebar-config-parser)
[![Go 版本](https://img.shields.io/github/go-mod/go-version/scagogogo/erlang-rebar-config-parser)](https://github.com/scagogogo/erlang-rebar-config-parser/blob/main/go.mod)
[![许可证](https://img.shields.io/github/license/scagogogo/erlang-rebar-config-parser)](https://github.com/scagogogo/erlang-rebar-config-parser/blob/main/LICENSE)

一个用于解析 Erlang rebar 配置文件的 Go 库，支持完整的 Erlang 数据类型和高级功能。该库允许您将 `rebar.config` 文件解析为结构化的 Go 对象，便于程序化访问和操作 Erlang 项目配置。

## 📚 文档

**📖 [完整文档网站](https://scagogogo.github.io/erlang-rebar-config-parser/zh/)**

文档包括：
- **[快速开始指南](https://scagogogo.github.io/erlang-rebar-config-parser/zh/guide/getting-started.html)** - 快速介绍和基本用法
- **[安装说明](https://scagogogo.github.io/erlang-rebar-config-parser/zh/guide/installation.html)** - 详细设置指南
- **[API 参考](https://scagogogo.github.io/erlang-rebar-config-parser/zh/api/)** - 完整的 API 文档
- **[示例](https://scagogogo.github.io/erlang-rebar-config-parser/zh/examples/)** - 实际应用示例
- **[高级用法](https://scagogogo.github.io/erlang-rebar-config-parser/zh/guide/advanced-usage.html)** - 复杂场景和最佳实践

## 🌟 特性

- **完整的 Erlang 支持**: 解析所有常见的 Erlang 数据类型（原子、字符串、整数、浮点数、元组、列表）
- **多种输入源**: 从文件、字符串或任何 `io.Reader` 解析
- **辅助方法**: 便捷访问常见配置部分（deps、erl_opts、profiles 等）
- **美化输出**: 可配置缩进和格式化
- **术语比较**: 类型感知的 Erlang 术语相等性比较
- **转义序列处理**: 正确处理字符串和原子中的转义字符
- **错误报告**: 带有位置信息的详细错误消息
- **线程安全解析**: 解析过程安全（解析结果需要同步）

## 🚀 快速开始

### 安装

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

### 基本用法

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 从文件解析
    config, err := parser.ParseFile("rebar.config")
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 访问常见配置元素
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("应用程序: %s\n", appName)
    }
    
    if deps, ok := config.GetDeps(); ok {
        fmt.Println("找到依赖项！")
        // 处理依赖项...
    }
    
    if erlOpts, ok := config.GetErlOpts(); ok {
        fmt.Println("找到 Erlang 选项！")
        // 处理 Erlang 选项...
    }
    
    // 使用 2 个空格缩进美化输出
    fmt.Println("\n格式化的配置:")
    fmt.Println(config.Format(2))
}
```

### 从不同来源解析

```go
// 从字符串解析
configStr := `{erl_opts, [debug_info]}.`
config, err := parser.Parse(configStr)

// 从 io.Reader 解析
file, err := os.Open("rebar.config")
if err == nil {
    defer file.Close()
    config, err = parser.ParseReader(file)
}

// 从 HTTP 响应解析
resp, err := http.Get("https://example.com/rebar.config")
if err == nil {
    defer resp.Body.Close()
    config, err = parser.ParseReader(resp.Body)
}
```

## 📋 支持的配置元素

库为常见的 rebar.config 部分提供辅助方法：

| 方法 | 描述 | 示例 |
|------|------|------|
| `GetDeps()` | 依赖项 | `{deps, [{cowboy, "2.9.0"}]}` |
| `GetErlOpts()` | Erlang 编译器选项 | `{erl_opts, [debug_info]}` |
| `GetAppName()` | 应用程序名称 | `{app_name, my_app}` |
| `GetPlugins()` | Rebar3 插件 | `{plugins, [rebar3_hex]}` |
| `GetProfilesConfig()` | 构建配置文件 | `{profiles, [{test, [...]}]}` |
| `GetRelxConfig()` | 发布配置 | `{relx, [{release, {...}}]}` |

## 🔧 高级功能

### 术语类型处理

```go
// 处理不同的 Erlang 术语类型
for _, term := range config.Terms {
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
        fmt.Printf("包含 %d 个元素的元组\n", len(t.Elements))
    case parser.List:
        fmt.Printf("包含 %d 个元素的列表\n", len(t.Elements))
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
```

## 📁 示例

[examples/](examples/) 目录包含实用示例：

- **[基本解析](examples/basic-parsing/)** - 简单解析示例
- **[美化输出器](examples/prettyprint/)** - 配置格式化工具
- **[配置分析](examples/analysis/)** - 高级配置分析

运行示例：

```bash
# 构建并运行美化输出器
make examples
cd examples
./prettyprint ../testdata/sample.config
```

## 🛠️ 开发

### 前提条件

- Go 1.18 或更高版本
- Node.js 18+（用于文档）

### 开发设置

```bash
# 克隆仓库
git clone https://github.com/scagogogo/erlang-rebar-config-parser.git
cd erlang-rebar-config-parser

# 设置开发环境
make dev-setup

# 运行测试
make test

# 运行带覆盖率的测试
make test-coverage

# 启动文档开发服务器
make docs-dev
```

### 可用的 Make 命令

```bash
make help              # 显示所有可用命令
make test              # 运行测试
make test-coverage     # 运行带覆盖率的测试
make lint              # 运行代码检查
make fmt               # 格式化代码
make docs-dev          # 启动文档服务器
make docs-build        # 构建文档
make examples          # 构建示例程序
make clean             # 清理构建产物
```

## 🤝 贡献

欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详情。

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 链接

- **[文档网站](https://scagogogo.github.io/erlang-rebar-config-parser/zh/)** - 完整文档
- **[GitHub 仓库](https://github.com/scagogogo/erlang-rebar-config-parser)** - 源代码
- **[问题反馈](https://github.com/scagogogo/erlang-rebar-config-parser/issues)** - 错误报告和功能请求
- **[版本发布](https://github.com/scagogogo/erlang-rebar-config-parser/releases)** - 版本历史

## ⭐ Star 历史

如果您觉得这个项目有用，请考虑在 GitHub 上给它一个 star！
