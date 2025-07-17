# 安装

本指南涵盖了在 Go 项目中安装和设置 Erlang Rebar 配置解析器的不同方法。

## 要求

- **Go 版本**: 1.18 或更高版本
- **操作系统**: Linux、macOS、Windows
- **架构**: amd64、arm64

## 标准安装

### 使用 go get（推荐）

安装库的最简单方法是使用 `go get`：

```bash
go get github.com/scagogogo/erlang-rebar-config-parser
```

这将下载最新版本并将其添加到您的 `go.mod` 文件中。

### 特定版本

要安装特定版本：

```bash
go get github.com/scagogogo/erlang-rebar-config-parser@v1.0.0
```

### 最新开发版本

要从主分支获取最新开发版本：

```bash
go get github.com/scagogogo/erlang-rebar-config-parser@main
```

## 项目设置

### 新项目

如果您正在开始一个新项目：

```bash
# 创建新目录
mkdir my-rebar-parser
cd my-rebar-parser

# 初始化 Go 模块
go mod init my-rebar-parser

# 安装库
go get github.com/scagogogo/erlang-rebar-config-parser

# 创建 main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    config, err := parser.Parse(`{erl_opts, [debug_info]}.`)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("解析了 %d 个项\n", len(config.Terms))
}
EOF

# 运行程序
go run main.go
```

### 现有项目

对于现有的 Go 项目，只需添加导入并运行：

```bash
go mod tidy
```

这将在您构建或运行项目时自动下载库。

## 导入语句

将导入添加到您的 Go 文件中：

```go
import "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
```

### 常见导入模式

```go
// 标准导入
import "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"

// 使用别名
import rebarparser "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"

// 多个导入
import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)
```

## 验证

### 快速测试

创建一个简单的测试来验证安装：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 测试基本解析
    config, err := parser.Parse(`{test, "installation"}.`)
    if err != nil {
        log.Fatalf("安装测试失败: %v", err)
    }
    
    // 测试辅助方法
    if term, ok := config.GetTerm("test"); ok {
        fmt.Printf("✓ 安装成功！解析了: %s\n", term.String())
    } else {
        log.Fatal("✗ 安装测试失败: 无法检索项")
    }
    
    // 测试格式化
    formatted := config.Format(2)
    fmt.Printf("✓ 格式化工作正常！输出:\n%s", formatted)
}
```

运行测试：

```bash
go run main.go
```

预期输出：
```
✓ 安装成功！解析了: {test, "installation"}
✓ 格式化工作正常！输出:
{test, "installation"}.
```

## 故障排除

### 常见问题

#### 1. Go 版本太旧

**错误**: `go: module requires Go 1.18`

**解决方案**: 将 Go 更新到版本 1.18 或更高版本：

```bash
# 检查当前版本
go version

# 更新 Go（方法因操作系统而异）
# 在 macOS 上使用 Homebrew:
brew install go

# 在 Ubuntu 上:
sudo snap install go --classic

# 或从 https://golang.org/dl/ 下载
```

#### 2. 找不到模块

**错误**: `cannot find module providing package`

**解决方案**: 确保您在 Go 模块目录中：

```bash
# 检查 go.mod 是否存在
ls go.mod

# 如果不存在，初始化模块
go mod init your-project-name

# 然后安装
go get github.com/scagogogo/erlang-rebar-config-parser
```

#### 3. 导入路径问题

**错误**: `package github.com/scagogogo/erlang-rebar-config-parser/pkg/parser is not in GOROOT`

**解决方案**: 确保您使用的是 Go 模块（而不是 GOPATH 模式）：

```bash
# 检查 Go 环境
go env GOMOD

# 应该显示 go.mod 文件的路径，而不是空
# 如果为空，您处于 GOPATH 模式 - 创建 go.mod:
go mod init your-project
```

#### 4. 网络问题

**错误**: `dial tcp: lookup proxy.golang.org: no such host`

**解决方案**: 配置 Go 代理或使用直接模式：

```bash
# 使用直接模式（绕过代理）
export GOPROXY=direct

# 或配置代理
export GOPROXY=https://proxy.golang.org,direct

# 对于中国用户
export GOPROXY=https://goproxy.cn,direct
```

### 依赖冲突

如果遇到依赖冲突：

```bash
# 清理模块缓存
go clean -modcache

# 更新依赖项
go get -u github.com/scagogogo/erlang-rebar-config-parser

# 整理
go mod tidy
```

## IDE 设置

### VS Code

安装 Go 扩展并确保您的工作区配置正确：

1. 安装 [Go 扩展](https://marketplace.visualstudio.com/items?itemName=golang.Go)
2. 在 VS Code 中打开您的项目文件夹
3. 按 `Ctrl+Shift+P` 并运行 "Go: Install/Update Tools"
4. 扩展应该自动检测导入的包

### GoLand/IntelliJ

1. 在 GoLand 中打开您的项目
2. IDE 应该自动下载依赖项
3. 如果没有，转到 File → Sync Dependencies

## Docker 设置

如果您使用 Docker：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 下一步

现在您已经安装了库：

1. **[基本用法](./basic-usage)** - 学习常见使用模式
2. **[API 参考](../api/)** - 探索完整的 API
3. **[示例](../examples/)** - 查看实际示例

## 获取帮助

如果您在安装过程中遇到问题：

- 检查 [GitHub Issues](https://github.com/scagogogo/erlang-rebar-config-parser/issues)
- 在 [GitHub 讨论](https://github.com/scagogogo/erlang-rebar-config-parser/discussions) 中提问
- 查看上面的[故障排除部分](#故障排除)
