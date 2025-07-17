# Erlang Rebar 配置解析器

[English Documentation](README.md) | [📖 完整文档](https://scagogogo.github.io/erlang-rebar-config-parser/zh/)

一个用于解析 Erlang rebar 配置文件的 Go 库。该库允许您将 `rebar.config` 文件解析为结构化的 Go 对象，使得以编程方式访问和操作 Erlang 项目配置变得简单。

## 📚 文档

- **[完整文档](https://scagogogo.github.io/erlang-rebar-config-parser/zh/)** - 完整的文档网站
- **[快速开始指南](https://scagogogo.github.io/erlang-rebar-config-parser/zh/guide/getting-started)** - 快速入门教程
- **[API 参考](https://scagogogo.github.io/erlang-rebar-config-parser/zh/api/)** - 完整的 API 文档
- **[示例](https://scagogogo.github.io/erlang-rebar-config-parser/zh/examples/)** - 实际应用示例

## 🌟 特性

- 将 rebar.config 文件解析为结构化的 Go 对象
- 支持所有常见的 Erlang 术语类型（元组、列表、原子、字符串、数字）
- 提供辅助方法轻松访问常见配置元素
- 完全支持嵌套数据结构
- 正确处理注释和空白字符
- 支持可配置缩进的美化输出
- 提供比较功能检查术语相等性
- 通过 GitHub Actions 进行持续集成
- 提供中英文完整文档和示例
- 98% 测试覆盖率，包含全面的边缘情况测试

## 📦 安装

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

## 🚀 快速开始

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
    
    // 格式化并打印配置
    fmt.Println("\n格式化的配置:")
    fmt.Println(config.Format(2))
}
```

## 📖 支持的 Erlang 术语类型

| Erlang 类型 | 示例 | Go 表示 |
|-------------|------|---------|
| 原子 | `atom_name`, `'quoted-atom'` | `Atom{Value: "atom_name", IsQuoted: false}` |
| 字符串 | `"hello world"` | `String{Value: "hello world"}` |
| 整数 | `123`, `-42` | `Integer{Value: 123}` |
| 浮点数 | `3.14`, `-1.5e-3` | `Float{Value: 3.14}` |
| 元组 | `{key, value}` | `Tuple{Elements: []Term{...}}` |
| 列表 | `[1, 2, 3]` | `List{Elements: []Term{...}}` |

## 🔧 主要功能

### 解析功能

```go
// 从文件解析
config, err := parser.ParseFile("rebar.config")

// 从字符串解析
config, err := parser.Parse(configString)

// 从 io.Reader 解析
config, err := parser.ParseReader(reader)
```

### 配置访问

```go
// 获取依赖项
deps, ok := config.GetDeps()

// 获取 Erlang 编译选项
erlOpts, ok := config.GetErlOpts()

// 获取应用程序名称
appName, ok := config.GetAppName()

// 获取插件
plugins, ok := config.GetPlugins()

// 获取配置文件
profiles, ok := config.GetProfilesConfig()
```

### 格式化输出

```go
// 使用 2 个空格缩进格式化
formatted := config.Format(2)
fmt.Println(formatted)
```

## 📋 示例配置文件

```erlang
%% rebar.config 示例
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

{profiles, [
    {dev, [
        {deps, [
            {sync, "0.1.3"}
        ]}
    ]},
    {test, [
        {deps, [
            {proper, "1.3.0"},
            {meck, "0.9.0"}
        ]}
    ]}
]}.

{relx, [
    {release, {my_app, "0.1.0"}, [my_app, sasl]},
    {dev_mode, true},
    {include_erts, false}
]}.
```

## 🧪 测试

运行测试套件：

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📊 项目状态

[![Go Tests and Examples](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/go.yml/badge.svg)](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/go.yml)
[![Documentation](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/docs.yml/badge.svg)](https://github.com/scagogogo/erlang-rebar-config-parser/actions/workflows/docs.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/erlang-rebar-config-parser)](https://goreportcard.com/report/github.com/scagogogo/erlang-rebar-config-parser)
[![Go Version](https://img.shields.io/github/go-mod/go-version/scagogogo/erlang-rebar-config-parser)](https://github.com/scagogogo/erlang-rebar-config-parser/blob/main/go.mod)
[![License](https://img.shields.io/github/license/scagogogo/erlang-rebar-config-parser)](https://github.com/scagogogo/erlang-rebar-config-parser/blob/main/LICENSE)

- ✅ **98% 测试覆盖率** - 全面的测试套件
- ✅ **持续集成** - 自动化测试和部署
- ✅ **完整文档** - 中英文双语文档
- ✅ **生产就绪** - 稳定可靠的 API

## 🤝 贡献

欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详情。

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- 感谢 Erlang/OTP 团队提供优秀的 Erlang 语言
- 感谢 rebar3 团队提供强大的构建工具
- 感谢所有贡献者和用户的支持

## 📞 联系方式

- **GitHub Issues**: [报告问题](https://github.com/scagogogo/erlang-rebar-config-parser/issues)
- **GitHub Discussions**: [讨论和问题](https://github.com/scagogogo/erlang-rebar-config-parser/discussions)
- **文档**: [在线文档](https://scagogogo.github.io/erlang-rebar-config-parser/zh/)
