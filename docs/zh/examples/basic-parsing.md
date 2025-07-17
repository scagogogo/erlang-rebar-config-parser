# 基本解析

此示例演示了 Erlang Rebar 配置解析器的基本解析功能。您将学习如何从不同来源解析配置并访问解析后的数据。

## 简单文件解析

让我们从解析 rebar.config 文件的基本示例开始：

### 示例 rebar.config

```erlang
%% 基本 rebar.config
{erl_opts, [debug_info, warnings_as_errors]}.

{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"},
    {lager, "3.9.2"}
]}.

{app_name, my_awesome_app}.
```

### 解析代码

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 解析配置文件
    config, err := parser.ParseFile("rebar.config")
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 显示基本信息
    fmt.Printf("成功解析配置！\n")
    fmt.Printf("顶级项数量: %d\n", len(config.Terms))
    
    // 遍历所有项
    for i, term := range config.Terms {
        fmt.Printf("项 %d: %s\n", i+1, term.String())
    }
}
```

### 预期输出

```
成功解析配置！
顶级项数量: 3
项 1: {erl_opts, [debug_info, warnings_as_errors]}
项 2: {deps, [{cowboy, "2.9.0"}, {jsx, "3.1.0"}, {lager, "3.9.2"}]}
项 3: {app_name, my_awesome_app}
```

## 从字符串解析

有时您需要从字符串解析配置内容：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 配置作为字符串
    configContent := `
    {erl_opts, [debug_info]}.
    {deps, [
        {cowboy, "2.9.0"}
    ]}.
    `
    
    // 从字符串解析
    config, err := parser.Parse(configContent)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Printf("从字符串解析了 %d 个项\n", len(config.Terms))
    
    // 访问原始内容
    fmt.Printf("原始内容长度: %d 个字符\n", len(config.Raw))
}
```

## 从 Reader 解析

对于更高级的场景，您可以从任何 `io.Reader` 解析：

```go
package main

import (
    "bytes"
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 示例 1: 从字符串读取器解析
    configStr := `{erl_opts, [debug_info]}.`
    reader := strings.NewReader(configStr)
    
    config, err := parser.ParseReader(reader)
    if err != nil {
        log.Fatalf("从读取器解析失败: %v", err)
    }
    
    fmt.Printf("从字符串读取器解析: %d 个项\n", len(config.Terms))
    
    // 示例 2: 从字节缓冲区解析
    var buf bytes.Buffer
    buf.WriteString(`{deps, [{jsx, "3.1.0"}]}.`)
    
    config2, err := parser.ParseReader(&buf)
    if err != nil {
        log.Fatalf("从缓冲区解析失败: %v", err)
    }
    
    fmt.Printf("从缓冲区解析: %d 个项\n", len(config2.Terms))
}
```

## 理解术语类型

让我们检查不同的 Erlang 术语类型及其表示方式：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 包含各种术语类型的配置
    configContent := `
    {atom_example, simple_atom}.
    {quoted_atom_example, 'quoted-atom'}.
    {string_example, "hello world"}.
    {integer_example, 42}.
    {float_example, 3.14}.
    {negative_number, -123}.
    {scientific_notation, 1.5e-3}.
    {tuple_example, {key, value, 123}}.
    {list_example, [item1, item2, item3]}.
    {nested_example, {deps, [{cowboy, "2.9.0"}]}}.
    `
    
    config, err := parser.Parse(configContent)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 检查每个项
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
            name := tuple.Elements[0]
            value := tuple.Elements[1]
            
            fmt.Printf("配置: %s\n", name.String())
            fmt.Printf("  类型: %T\n", value)
            fmt.Printf("  值: %s\n", value.String())
            
            // 详细类型分析
            analyzeTermType(value)
            fmt.Println()
        }
    }
}

func analyzeTermType(term parser.Term) {
    switch t := term.(type) {
    case parser.Atom:
        fmt.Printf("  原子详情: 值='%s', 引号=%t\n", t.Value, t.IsQuoted)
    case parser.String:
        fmt.Printf("  字符串详情: 值='%s', 长度=%d\n", t.Value, len(t.Value))
    case parser.Integer:
        fmt.Printf("  整数详情: 值=%d\n", t.Value)
    case parser.Float:
        fmt.Printf("  浮点数详情: 值=%f\n", t.Value)
    case parser.Tuple:
        fmt.Printf("  元组详情: %d 个元素\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("    [%d]: %s (%T)\n", i, elem.String(), elem)
        }
    case parser.List:
        fmt.Printf("  列表详情: %d 个元素\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("    [%d]: %s (%T)\n", i, elem.String(), elem)
        }
    }
}
```

## 错误处理示例

### 处理解析错误

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 无效配置示例
    invalidConfigs := []string{
        `{incomplete_tuple`,           // 缺少右大括号
        `{unterminated, "string}`,     // 未终止的字符串
        `{invalid_number, 12.34.56}`,  // 无效数字格式
        `{missing_comma [item1 item2]}`, // 列表中缺少逗号
    }
    
    for i, configStr := range invalidConfigs {
        fmt.Printf("测试无效配置 %d:\n", i+1)
        fmt.Printf("内容: %s\n", configStr)
        
        _, err := parser.Parse(configStr)
        if err != nil {
            fmt.Printf("错误（预期）: %v\n", err)
            
            // 分类错误类型
            if strings.Contains(err.Error(), "syntax error") {
                fmt.Println("  类别: 语法错误")
            } else if strings.Contains(err.Error(), "unexpected") {
                fmt.Println("  类别: 意外字符")
            } else {
                fmt.Println("  类别: 其他解析错误")
            }
        } else {
            fmt.Println("意外: 没有发生错误！")
        }
        fmt.Println()
    }
}
```

## 完整示例

这是一个演示所有基本解析概念的完整示例：

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
    // 检查配置文件是否存在
    configFile := "rebar.config"
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        // 如果不存在则创建示例配置
        createSampleConfig(configFile)
    }
    
    // 解析配置
    config, err := parser.ParseFile(configFile)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 显示结果
    fmt.Printf("✓ 成功解析 %s\n", configFile)
    fmt.Printf("  项: %d\n", len(config.Terms))
    fmt.Printf("  原始大小: %d 字节\n", len(config.Raw))
    
    // 分析每个项
    fmt.Println("\n项分析:")
    for i, term := range config.Terms {
        fmt.Printf("  [%d] %T: %s\n", i+1, term, term.String())
    }
    
    fmt.Println("\n格式化输出:")
    fmt.Println(config.Format(2))
}

func createSampleConfig(filename string) {
    content := `{erl_opts, [debug_info, warnings_as_errors]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
{app_name, sample_app}.
`
    
    err := os.WriteFile(filename, []byte(content), 0644)
    if err != nil {
        log.Fatalf("创建示例配置失败: %v", err)
    }
    
    fmt.Printf("创建了示例配置: %s\n", filename)
}
```

## 下一步

现在您了解了基本解析：

1. **[配置访问](./config-access)** - 学习使用辅助方法
2. **[美化输出](./pretty-printing)** - 美观地格式化配置
3. **[术语比较](./comparison)** - 比较配置
4. **[复杂分析](./complex-analysis)** - 高级解析场景
