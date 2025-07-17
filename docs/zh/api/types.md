# 类型定义

本节记录了 Erlang Rebar 配置解析器库中使用的所有数据类型。

## RebarConfig

```go
type RebarConfig struct {
    Raw   string
    Terms []Term
}
```

表示已解析的 rebar.config 文件，包含所有配置项。

### 字段

- `Raw` (string): 配置文件的原始内容
- `Terms` ([]Term): 所有顶级配置项的列表

### 示例

```go
config, _ := parser.Parse(`{erl_opts, [debug_info]}.`)
fmt.Printf("原始内容: %s\n", config.Raw)
fmt.Printf("项数量: %d\n", len(config.Terms))
```

---

## Term 接口

```go
type Term interface {
    String() string
    Compare(other Term) bool
}
```

所有 Erlang 术语类型实现的基础接口。

### 方法

- `String()`: 返回术语的字符串表示
- `Compare(other Term)`: 比较此术语与另一个术语的相等性

---

## Atom

```go
type Atom struct {
    Value    string
    IsQuoted bool
}
```

表示 Erlang 原子（符号）。

### 字段

- `Value` (string): 原子的值
- `IsQuoted` (bool): 原子在原始语法中是否被引号包围

### 方法

```go
func (a Atom) String() string
func (a Atom) Compare(other Term) bool
```

### 示例

```go
// 普通原子: debug_info
atom1 := parser.Atom{Value: "debug_info", IsQuoted: false}
fmt.Println(atom1.String()) // 输出: debug_info

// 引号原子: 'complex-name'
atom2 := parser.Atom{Value: "complex-name", IsQuoted: true}
fmt.Println(atom2.String()) // 输出: 'complex-name'

// 比较（忽略 IsQuoted）
atom3 := parser.Atom{Value: "test", IsQuoted: false}
atom4 := parser.Atom{Value: "test", IsQuoted: true}
fmt.Println(atom3.Compare(atom4)) // 输出: true
```

---

## String

```go
type String struct {
    Value string
}
```

表示 Erlang 字符串（双引号文本）。

### 字段

- `Value` (string): 字符串内容

### 方法

```go
func (s String) String() string
func (s String) Compare(other Term) bool
```

### 示例

```go
str := parser.String{Value: "hello world"}
fmt.Println(str.String()) // 输出: "hello world"

str1 := parser.String{Value: "test"}
str2 := parser.String{Value: "test"}
fmt.Println(str1.Compare(str2)) // 输出: true
```

---

## Integer

```go
type Integer struct {
    Value int64
}
```

表示 Erlang 整数。

### 字段

- `Value` (int64): 整数值

### 方法

```go
func (i Integer) String() string
func (i Integer) Compare(other Term) bool
```

### 示例

```go
num := parser.Integer{Value: 42}
fmt.Println(num.String()) // 输出: 42

num1 := parser.Integer{Value: 123}
num2 := parser.Integer{Value: 123}
fmt.Println(num1.Compare(num2)) // 输出: true
```

---

## Float

```go
type Float struct {
    Value float64
}
```

表示 Erlang 浮点数。

### 字段

- `Value` (float64): 浮点值

### 方法

```go
func (f Float) String() string
func (f Float) Compare(other Term) bool
```

### 示例

```go
num := parser.Float{Value: 3.14}
fmt.Println(num.String()) // 输出: 3.14

// 科学记数法
num2 := parser.Float{Value: 1.5e-3}
fmt.Println(num2.String()) // 输出: 0.0015
```

---

## Tuple

```go
type Tuple struct {
    Elements []Term
}
```

表示 Erlang 元组 `{elem1, elem2, ...}`。

### 字段

- `Elements` ([]Term): 元组中的元素列表

### 方法

```go
func (t Tuple) String() string
func (t Tuple) Compare(other Term) bool
```

### 示例

```go
// 简单元组: {key, value}
tuple := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "key"},
        parser.String{Value: "value"},
    },
}
fmt.Println(tuple.String()) // 输出: {key, "value"}

// 嵌套元组: {deps, [{cowboy, "2.9.0"}]}
nestedTuple := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "deps"},
        parser.List{
            Elements: []parser.Term{
                parser.Tuple{
                    Elements: []parser.Term{
                        parser.Atom{Value: "cowboy"},
                        parser.String{Value: "2.9.0"},
                    },
                },
            },
        },
    },
}
```

---

## List

```go
type List struct {
    Elements []Term
}
```

表示 Erlang 列表 `[elem1, elem2, ...]`。

### 字段

- `Elements` ([]Term): 列表中的元素列表

### 方法

```go
func (l List) String() string
func (l List) Compare(other Term) bool
```

### 示例

```go
// 简单列表: [debug_info, warnings_as_errors]
list := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "warnings_as_errors"},
    },
}
fmt.Println(list.String()) // 输出: [debug_info, warnings_as_errors]

// 混合类型列表: [atom, "string", 123]
mixedList := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "atom"},
        parser.String{Value: "string"},
        parser.Integer{Value: 123},
    },
}
```

---

## 类型检查和转换

### 安全类型断言

```go
func processTerms(terms []parser.Term) {
    for _, term := range terms {
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
            fmt.Printf("元组，包含 %d 个元素\n", len(t.Elements))
        case parser.List:
            fmt.Printf("列表，包含 %d 个元素\n", len(t.Elements))
        default:
            fmt.Printf("未知术语类型: %T\n", t)
        }
    }
}
```

### 处理嵌套结构

```go
func extractDependencies(config *parser.RebarConfig) []string {
    var deps []string
    
    if depsTerms, ok := config.GetDeps(); ok && len(depsTerms) > 0 {
        if depsList, ok := depsTerms[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                        deps = append(deps, atom.Value)
                    }
                }
            }
        }
    }
    
    return deps
}
```
