# RebarConfig 方法

`RebarConfig` 类型提供了便捷的方法来访问常见的配置元素和格式化配置。

## 配置访问方法

### GetTerm

```go
func (c *RebarConfig) GetTerm(name string) (Term, bool)
```

通过名称从配置中检索特定的顶级项。

#### 参数

- `name` (string): 要查找的项的名称

#### 返回值

- `Term`: 找到的项
- `bool`: 是否找到该项

#### 示例

```go
config, _ := parser.Parse(`{erl_opts, [debug_info]}.`)

term, ok := config.GetTerm("erl_opts")
if ok {
    fmt.Printf("找到 erl_opts: %s\n", term.String())
}
```

---

### GetTupleElements

```go
func (c *RebarConfig) GetTupleElements(name string) ([]Term, bool)
```

获取命名元组的元素（不包括名称本身）。

#### 参数

- `name` (string): 元组的名称

#### 返回值

- `[]Term`: 元组元素列表（不包括名称）
- `bool`: 是否找到元组并且有元素

#### 示例

```go
config, _ := parser.Parse(`{deps, [{cowboy, "2.9.0"}]}.`)

elements, ok := config.GetTupleElements("deps")
if ok {
    fmt.Printf("deps 有 %d 个元素\n", len(elements))
    // elements[0] 将是列表 [{cowboy, "2.9.0"}]
}
```

---

### GetDeps

```go
func (c *RebarConfig) GetDeps() ([]Term, bool)
```

检索依赖项配置。

#### 返回值

- `[]Term`: 依赖项列表
- `bool`: 是否找到依赖项

#### 示例

```go
config, _ := parser.Parse(`
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
`)

deps, ok := config.GetDeps()
if ok && len(deps) > 0 {
    if depsList, ok := deps[0].(parser.List); ok {
        fmt.Printf("找到 %d 个依赖项\n", len(depsList.Elements))
        
        for _, dep := range depsList.Elements {
            if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
                if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("- %s\n", atom.Value)
                }
            }
        }
    }
}
```

---

### GetErlOpts

```go
func (c *RebarConfig) GetErlOpts() ([]Term, bool)
```

检索 Erlang 编译选项。

#### 返回值

- `[]Term`: Erlang 编译选项列表
- `bool`: 是否找到 erl_opts

#### 示例

```go
config, _ := parser.Parse(`{erl_opts, [debug_info, warnings_as_errors]}.`)

opts, ok := config.GetErlOpts()
if ok && len(opts) > 0 {
    if optsList, ok := opts[0].(parser.List); ok {
        fmt.Printf("Erlang 选项:\n")
        for _, opt := range optsList.Elements {
            fmt.Printf("- %s\n", opt.String())
        }
    }
}
```

---

### GetAppName

```go
func (c *RebarConfig) GetAppName() (string, bool)
```

检索应用程序名称。

#### 返回值

- `string`: 应用程序名称
- `bool`: 是否找到应用程序名称

#### 示例

```go
config, _ := parser.Parse(`{app_name, "my_awesome_app"}.`)

appName, ok := config.GetAppName()
if ok {
    fmt.Printf("应用程序名称: %s\n", appName)
}

// 也适用于原子值
config2, _ := parser.Parse(`{app_name, my_app}.`)
appName2, ok := config2.GetAppName()
if ok {
    fmt.Printf("应用程序名称: %s\n", appName2)
}
```

---

### GetPlugins

```go
func (c *RebarConfig) GetPlugins() ([]Term, bool)
```

检索插件配置。

#### 返回值

- `[]Term`: 插件项列表
- `bool`: 是否找到插件

#### 示例

```go
config, _ := parser.Parse(`{plugins, [rebar3_hex, rebar3_auto]}.`)

plugins, ok := config.GetPlugins()
if ok && len(plugins) > 0 {
    if pluginsList, ok := plugins[0].(parser.List); ok {
        fmt.Printf("插件:\n")
        for _, plugin := range pluginsList.Elements {
            if atom, ok := plugin.(parser.Atom); ok {
                fmt.Printf("- %s\n", atom.Value)
            }
        }
    }
}
```

---

### GetRelxConfig

```go
func (c *RebarConfig) GetRelxConfig() ([]Term, bool)
```

检索 relx（发布）配置。

#### 返回值

- `[]Term`: relx 配置项列表
- `bool`: 是否找到 relx 配置

#### 示例

```go
config, _ := parser.Parse(`
{relx, [
    {release, {my_app, "0.1.0"}, [my_app, sasl]},
    {dev_mode, true},
    {include_erts, false}
]}.
`)

relx, ok := config.GetRelxConfig()
if ok && len(relx) > 0 {
    if relxList, ok := relx[0].(parser.List); ok {
        fmt.Printf("Relx 配置有 %d 个选项\n", len(relxList.Elements))
    }
}
```

---

### GetProfilesConfig

```go
func (c *RebarConfig) GetProfilesConfig() ([]Term, bool)
```

检索不同环境的配置文件配置。

#### 返回值

- `[]Term`: 配置文件配置项列表
- `bool`: 是否找到配置文件配置

#### 示例

```go
config, _ := parser.Parse(`
{profiles, [
    {dev, [{deps, [{sync, "0.1.3"}]}]},
    {test, [{deps, [{proper, "1.3.0"}]}]}
]}.
`)

profiles, ok := config.GetProfilesConfig()
if ok && len(profiles) > 0 {
    if profilesList, ok := profiles[0].(parser.List); ok {
        fmt.Printf("找到 %d 个配置文件\n", len(profilesList.Elements))
        
        for _, profile := range profilesList.Elements {
            if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                    fmt.Printf("- 配置文件: %s\n", atom.Value)
                }
            }
        }
    }
}
```

---

## 格式化方法

### Format

```go
func (c *RebarConfig) Format(indent int) string
```

返回配置的格式化字符串表示，使用指定的缩进。

#### 参数

- `indent` (int): 缩进的空格数（通常是 2 或 4）

#### 返回值

- `string`: 格式化的配置字符串

#### 示例

```go
config, _ := parser.Parse(`{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"}]}.`)

// 使用 2 个空格缩进格式化
formatted := config.Format(2)
fmt.Println(formatted)
// 输出:
// {deps, [
//   {cowboy, "2.9.0"},
//   {jsx, "3.1.0"}
// ]}.

// 使用 4 个空格缩进格式化
formatted4 := config.Format(4)
fmt.Println(formatted4)
// 输出:
// {deps, [
//     {cowboy, "2.9.0"},
//     {jsx, "3.1.0"}
// ]}.
```

---

## 使用模式

### 检查必需配置

```go
func validateConfig(config *parser.RebarConfig) error {
    if _, ok := config.GetDeps(); !ok {
        return fmt.Errorf("缺少必需的 'deps' 配置")
    }
    
    if _, ok := config.GetErlOpts(); !ok {
        return fmt.Errorf("缺少必需的 'erl_opts' 配置")
    }
    
    return nil
}
```

### 提取特定信息

```go
func analyzeConfig(config *parser.RebarConfig) {
    // 检查应用程序名称
    if appName, ok := config.GetAppName(); ok {
        fmt.Printf("应用程序: %s\n", appName)
    }
    
    // 计算依赖项数量
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            fmt.Printf("依赖项: %d\n", len(depsList.Elements))
        }
    }
    
    // 检查开发配置文件
    if profiles, ok := config.GetProfilesConfig(); ok && len(profiles) > 0 {
        if profilesList, ok := profiles[0].(parser.List); ok {
            for _, profile := range profilesList.Elements {
                if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
                    if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == "dev" {
                        fmt.Println("找到开发配置文件")
                        break
                    }
                }
            }
        }
    }
}
```

### 配置修改和格式化

```go
func modifyAndFormat(config *parser.RebarConfig) string {
    // 修改配置（根据需要创建新项）
    // ... 修改逻辑 ...
    
    // 使用一致的缩进格式化
    return config.Format(4)
}
```
