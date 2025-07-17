# 核心函数

核心函数提供了从不同来源解析 Erlang rebar 配置文件的主要入口点。

## ParseFile

```go
func ParseFile(path string) (*RebarConfig, error)
```

从给定的文件路径解析 rebar.config 文件。

### 参数

- `path` (string): rebar.config 文件的路径

### 返回值

- `*RebarConfig`: 解析后的配置对象
- `error`: 解析过程中发生的任何错误

### 示例

```go
config, err := parser.ParseFile("./rebar.config")
if err != nil {
    log.Fatalf("解析配置失败: %v", err)
}

fmt.Printf("解析了 %d 个配置项\n", len(config.Terms))
```

### 错误情况

- 文件不存在
- 权限被拒绝
- 文件中的 Erlang 语法无效

---

## Parse

```go
func Parse(input string) (*RebarConfig, error)
```

从给定的字符串内容解析 rebar.config。

### 参数

- `input` (string): 包含 Erlang 配置语法的字符串

### 返回值

- `*RebarConfig`: 解析后的配置对象
- `error`: 解析过程中发生的任何错误

### 示例

```go
configStr := `
{erl_opts, [debug_info]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
`

config, err := parser.Parse(configStr)
if err != nil {
    log.Fatalf("解析配置失败: %v", err)
}

// 访问解析后的配置
deps, ok := config.GetDeps()
if ok {
    fmt.Printf("找到依赖项: %v\n", deps)
}
```

### 错误情况

- 无效的 Erlang 语法
- 未终止的字符串或原子
- 格式错误的数字
- 缺少必需的标点符号（点号、逗号）

---

## ParseReader

```go
func ParseReader(r io.Reader) (*RebarConfig, error)
```

从给定的 io.Reader 解析 rebar.config。

### 参数

- `r` (io.Reader): 包含 Erlang 配置数据的读取器

### 返回值

- `*RebarConfig`: 解析后的配置对象
- `error`: 解析过程中发生的任何错误

### 示例

```go
// 从文件解析
file, err := os.Open("rebar.config")
if err != nil {
    log.Fatalf("打开文件失败: %v", err)
}
defer file.Close()

config, err := parser.ParseReader(file)
if err != nil {
    log.Fatalf("解析配置失败: %v", err)
}

// 从 HTTP 响应解析
resp, err := http.Get("https://example.com/rebar.config")
if err != nil {
    log.Fatalf("获取配置失败: %v", err)
}
defer resp.Body.Close()

config, err = parser.ParseReader(resp.Body)
if err != nil {
    log.Fatalf("解析远程配置失败: %v", err)
}
```

### 错误情况

- 底层读取器的读取错误
- 输入中的无效 Erlang 语法
- 网络错误（从网络源读取时）

---

## NewParser

```go
func NewParser(input string) *Parser
```

为给定的输入字符串创建新的解析器实例。这主要用于内部使用，但对于高级用例可能有用。

### 参数

- `input` (string): 要解析的输入字符串

### 返回值

- `*Parser`: 新的解析器实例

### 示例

```go
// 高级用法 - 通常不需要
parser := parser.NewParser(`{erl_opts, [debug_info]}.`)
// 直接使用解析器方法...
```

---

## 使用模式

### 从不同来源解析

```go
// 从文件
config1, err := parser.ParseFile("rebar.config")

// 从字符串
config2, err := parser.Parse(`{deps, []}.`)

// 从任何读取器
var buf bytes.Buffer
buf.WriteString(`{erl_opts, [debug_info]}.`)
config3, err := parser.ParseReader(&buf)
```

### 错误处理最佳实践

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    // 检查特定错误类型
    if os.IsNotExist(err) {
        log.Fatal("配置文件未找到")
    } else if strings.Contains(err.Error(), "syntax error") {
        log.Fatalf("配置语法无效: %v", err)
    } else {
        log.Fatalf("意外错误: %v", err)
    }
}
```

### 解析后验证

```go
config, err := parser.ParseFile("rebar.config")
if err != nil {
    log.Fatalf("解析错误: %v", err)
}

// 验证必需的部分是否存在
if _, ok := config.GetDeps(); !ok {
    log.Println("警告: 未找到依赖项")
}

if _, ok := config.GetErlOpts(); !ok {
    log.Println("警告: 未找到 Erlang 选项")
}
```
