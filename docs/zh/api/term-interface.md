# Term 接口

`Term` 接口是 Erlang Rebar 配置解析器类型系统的基础。所有 Erlang 数据类型都实现此接口，为字符串表示和比较提供一致的方法。

## 接口定义

```go
type Term interface {
    String() string
    Compare(other Term) bool
}
```

## 方法

### String()

返回术语的 Erlang 语法字符串表示。

#### 返回值

- `string`: 术语的 Erlang 语法表示

#### 示例

```go
atom := parser.Atom{Value: "debug_info", IsQuoted: false}
fmt.Println(atom.String()) // 输出: debug_info

quotedAtom := parser.Atom{Value: "complex-name", IsQuoted: true}
fmt.Println(quotedAtom.String()) // 输出: 'complex-name'

str := parser.String{Value: "hello world"}
fmt.Println(str.String()) // 输出: "hello world"

integer := parser.Integer{Value: 42}
fmt.Println(integer.String()) // 输出: 42

float := parser.Float{Value: 3.14}
fmt.Println(float.String()) // 输出: 3.14
```

### Compare(other Term)

比较此术语与另一个术语的相等性。

#### 参数

- `other` (Term): 要比较的术语

#### 返回值

- `bool`: 如果术语相等则为 `true`，否则为 `false`

#### 比较规则

1. **类型匹配**: 术语必须是相同类型才能相等
2. **值比较**: 实际值必须相同
3. **特殊情况**:
   - 对于 `Atom`: 只比较 `Value`，忽略 `IsQuoted`
   - 对于嵌套结构（`Tuple`、`List`）: 所有元素必须相等
   - 对于数字: 需要精确值匹配（无类型强制转换）

#### 示例

```go
// 原子比较（忽略 IsQuoted）
atom1 := parser.Atom{Value: "test", IsQuoted: false}
atom2 := parser.Atom{Value: "test", IsQuoted: true}
fmt.Println(atom1.Compare(atom2)) // 输出: true

// 字符串比较
str1 := parser.String{Value: "hello"}
str2 := parser.String{Value: "hello"}
str3 := parser.String{Value: "world"}
fmt.Println(str1.Compare(str2)) // 输出: true
fmt.Println(str1.Compare(str3)) // 输出: false

// 数字比较
int1 := parser.Integer{Value: 42}
int2 := parser.Integer{Value: 42}
float1 := parser.Float{Value: 42.0}
fmt.Println(int1.Compare(int2))   // 输出: true
fmt.Println(int1.Compare(float1)) // 输出: false（不同类型）
```

## 处理术语

### 类型断言

使用 Go 的类型断言来处理特定的术语类型：

```go
func processTerm(term parser.Term) {
    switch t := term.(type) {
    case parser.Atom:
        fmt.Printf("原子: %s\n", t.Value)
        if t.IsQuoted {
            fmt.Println("  (引号)")
        }
    case parser.String:
        fmt.Printf("字符串: %s\n", t.Value)
    case parser.Integer:
        fmt.Printf("整数: %d\n", t.Value)
    case parser.Float:
        fmt.Printf("浮点数: %f\n", t.Value)
    case parser.Tuple:
        fmt.Printf("元组，包含 %d 个元素:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    case parser.List:
        fmt.Printf("列表，包含 %d 个元素:\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("  [%d]: %s\n", i, elem.String())
        }
    default:
        fmt.Printf("未知术语类型: %T\n", t)
    }
}
```

### 安全类型检查

```go
func isAtom(term parser.Term) bool {
    _, ok := term.(parser.Atom)
    return ok
}

func getAtomValue(term parser.Term) (string, bool) {
    if atom, ok := term.(parser.Atom); ok {
        return atom.Value, true
    }
    return "", false
}

func getStringValue(term parser.Term) (string, bool) {
    if str, ok := term.(parser.String); ok {
        return str.Value, true
    }
    return "", false
}
```

### 处理集合

```go
func processCollection(term parser.Term) {
    switch t := term.(type) {
    case parser.Tuple:
        fmt.Printf("处理包含 %d 个元素的元组\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("元素 %d: %s\n", i, elem.String())
        }
    case parser.List:
        fmt.Printf("处理包含 %d 个元素的列表\n", len(t.Elements))
        for i, elem := range t.Elements {
            fmt.Printf("元素 %d: %s\n", i, elem.String())
        }
    default:
        fmt.Println("不是集合类型")
    }
}
```

## 比较示例

### 简单比较

```go
// 创建术语
atom1 := parser.Atom{Value: "debug_info"}
atom2 := parser.Atom{Value: "debug_info"}
atom3 := parser.Atom{Value: "warnings_as_errors"}

// 比较
fmt.Println(atom1.Compare(atom2)) // true
fmt.Println(atom1.Compare(atom3)) // false

// 跨类型比较
str := parser.String{Value: "debug_info"}
fmt.Println(atom1.Compare(str)) // false（不同类型）
```

### 复杂结构比较

```go
// 创建相同的元组
tuple1 := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "cowboy"},
        parser.String{Value: "2.9.0"},
    },
}

tuple2 := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "cowboy"},
        parser.String{Value: "2.9.0"},
    },
}

tuple3 := parser.Tuple{
    Elements: []parser.Term{
        parser.Atom{Value: "cowboy"},
        parser.String{Value: "2.8.0"}, // 不同版本
    },
}

fmt.Println(tuple1.Compare(tuple2)) // true
fmt.Println(tuple1.Compare(tuple3)) // false
```

### 列表比较

```go
list1 := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "warnings_as_errors"},
    },
}

list2 := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        parser.Atom{Value: "warnings_as_errors"},
    },
}

list3 := parser.List{
    Elements: []parser.Term{
        parser.Atom{Value: "debug_info"},
        // 缺少第二个元素
    },
}

fmt.Println(list1.Compare(list2)) // true
fmt.Println(list1.Compare(list3)) // false（不同长度）
```

## 实用函数

### 在集合中查找术语

```go
func findAtomInList(list parser.List, atomValue string) bool {
    for _, elem := range list.Elements {
        if atom, ok := elem.(parser.Atom); ok && atom.Value == atomValue {
            return true
        }
    }
    return false
}

func findTupleByFirstAtom(list parser.List, atomValue string) (parser.Tuple, bool) {
    for _, elem := range list.Elements {
        if tuple, ok := elem.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == atomValue {
                return tuple, true
            }
        }
    }
    return parser.Tuple{}, false
}
```

### 术语验证

```go
func validateDependency(term parser.Term) error {
    tuple, ok := term.(parser.Tuple)
    if !ok {
        return fmt.Errorf("依赖项必须是元组")
    }
    
    if len(tuple.Elements) < 2 {
        return fmt.Errorf("依赖项元组必须至少有 2 个元素")
    }
    
    if _, ok := tuple.Elements[0].(parser.Atom); !ok {
        return fmt.Errorf("依赖项名称必须是原子")
    }
    
    switch tuple.Elements[1].(type) {
    case parser.String, parser.Tuple:
        // 有效的版本规范
        return nil
    default:
        return fmt.Errorf("依赖项版本必须是字符串或元组")
    }
}
```
