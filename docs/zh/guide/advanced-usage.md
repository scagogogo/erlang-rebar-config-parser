# 高级用法

本指南涵盖了 Erlang Rebar 配置解析器的高级使用模式和最佳实践，包括性能优化、错误处理策略和复杂用例。

## 性能优化

### 内存管理

处理大型配置或处理多个文件时，考虑内存使用：

```go
func processLargeConfigs(configPaths []string) error {
    for _, path := range configPaths {
        // 一次处理一个配置以限制内存使用
        config, err := parser.ParseFile(path)
        if err != nil {
            log.Printf("解析 %s 失败: %v", path, err)
            continue
        }
        
        // 处理配置
        processConfig(config)
        
        // 允许垃圾回收
        config = nil
        runtime.GC()
    }
    
    return nil
}
```

### 并发处理

并发处理多个配置：

```go
func processConcurrently(configPaths []string, maxWorkers int) {
    jobs := make(chan string, len(configPaths))
    results := make(chan ConfigResult, len(configPaths))
    
    // 启动工作器
    for w := 0; w < maxWorkers; w++ {
        go configWorker(jobs, results)
    }
    
    // 发送作业
    for _, path := range configPaths {
        jobs <- path
    }
    close(jobs)
    
    // 收集结果
    for i := 0; i < len(configPaths); i++ {
        result := <-results
        if result.Error != nil {
            log.Printf("处理 %s 出错: %v", result.Path, result.Error)
        } else {
            log.Printf("成功处理 %s", result.Path)
        }
    }
}

type ConfigResult struct {
    Path   string
    Config *parser.RebarConfig
    Error  error
}

func configWorker(jobs <-chan string, results chan<- ConfigResult) {
    for path := range jobs {
        config, err := parser.ParseFile(path)
        results <- ConfigResult{
            Path:   path,
            Config: config,
            Error:  err,
        }
    }
}
```

## 高级错误处理

### 自定义错误类型

创建自定义错误类型以便更好地处理错误：

```go
type ConfigError struct {
    Type    string
    Path    string
    Line    int
    Column  int
    Message string
    Cause   error
}

func (e *ConfigError) Error() string {
    if e.Path != "" {
        return fmt.Sprintf("%s 错误在 %s: %s", e.Type, e.Path, e.Message)
    }
    return fmt.Sprintf("%s 错误: %s", e.Type, e.Message)
}

func (e *ConfigError) Unwrap() error {
    return e.Cause
}

func parseWithDetailedErrors(path string) (*parser.RebarConfig, error) {
    config, err := parser.ParseFile(path)
    if err != nil {
        // 使用上下文增强错误
        configErr := &ConfigError{
            Type:    "解析",
            Path:    path,
            Message: err.Error(),
            Cause:   err,
        }
        
        // 如果可用，尝试提取行/列信息
        if strings.Contains(err.Error(), "line") {
            // 从错误消息解析行信息
            // 这是特定于实现的
        }
        
        return nil, configErr
    }
    
    return config, nil
}
```

### 错误恢复策略

实现从部分解析失败中恢复的策略：

```go
func parseWithRecovery(content string) (*parser.RebarConfig, []error) {
    var errors []error
    
    // 首先尝试解析整个内容
    config, err := parser.Parse(content)
    if err == nil {
        return config, nil
    }
    
    errors = append(errors, err)
    
    // 如果解析失败，尝试解析单个术语
    lines := strings.Split(content, "\n")
    var validTerms []parser.Term
    
    for i, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "%") {
            continue // 跳过空行和注释
        }
        
        // 尝试将单行解析为术语
        if strings.HasSuffix(line, ".") {
            termConfig, err := parser.Parse(line)
            if err != nil {
                errors = append(errors, fmt.Errorf("第 %d 行: %w", i+1, err))
            } else if len(termConfig.Terms) > 0 {
                validTerms = append(validTerms, termConfig.Terms...)
            }
        }
    }
    
    if len(validTerms) > 0 {
        // 返回部分配置
        return &parser.RebarConfig{
            Terms: validTerms,
            Raw:   content,
        }, errors
    }
    
    return nil, errors
}
```

## 配置转换

### 配置迁移

在不同格式或版本之间迁移配置：

```go
type ConfigMigrator struct {
    migrations []Migration
}

type Migration interface {
    Name() string
    Description() string
    Apply(config *parser.RebarConfig) (*parser.RebarConfig, error)
    CanApply(config *parser.RebarConfig) bool
}

// 示例迁移：将旧式依赖项转换为新格式
type DependencyFormatMigration struct{}

func (m DependencyFormatMigration) Name() string {
    return "dependency-format-v2"
}

func (m DependencyFormatMigration) Description() string {
    return "将旧式依赖项格式转换为新格式"
}

func (m DependencyFormatMigration) CanApply(config *parser.RebarConfig) bool {
    // 检查配置是否使用旧格式
    if deps, ok := config.GetDeps(); ok && len(deps) > 0 {
        if depsList, ok := deps[0].(parser.List); ok {
            for _, dep := range depsList.Elements {
                // 检查旧格式模式
                if atom, ok := dep.(parser.Atom); ok {
                    // 旧格式：只有原子名称
                    _ = atom
                    return true
                }
            }
        }
    }
    return false
}

func (m DependencyFormatMigration) Apply(config *parser.RebarConfig) (*parser.RebarConfig, error) {
    newTerms := make([]parser.Term, len(config.Terms))
    copy(newTerms, config.Terms)
    
    // 查找并转换依赖项术语
    for i, term := range newTerms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == "deps" {
                newDeps, err := m.transformDependencies(tuple.Elements[1])
                if err != nil {
                    return nil, err
                }
                
                newTerms[i] = parser.Tuple{
                    Elements: []parser.Term{
                        tuple.Elements[0],
                        newDeps,
                    },
                }
            }
        }
    }
    
    return &parser.RebarConfig{
        Terms: newTerms,
        Raw:   "", // 格式化时将重新生成原始内容
    }, nil
}

func (m DependencyFormatMigration) transformDependencies(deps parser.Term) (parser.Term, error) {
    if depsList, ok := deps.(parser.List); ok {
        var newElements []parser.Term
        
        for _, dep := range depsList.Elements {
            if atom, ok := dep.(parser.Atom); ok {
                // 将原子转换为 {atom, "latest"} 元组
                newDep := parser.Tuple{
                    Elements: []parser.Term{
                        atom,
                        parser.String{Value: "latest"},
                    },
                }
                newElements = append(newElements, newDep)
            } else {
                // 保持现有格式
                newElements = append(newElements, dep)
            }
        }
        
        return parser.List{Elements: newElements}, nil
    }
    
    return deps, nil
}

func (cm *ConfigMigrator) Migrate(config *parser.RebarConfig) (*parser.RebarConfig, error) {
    current := config
    
    for _, migration := range cm.migrations {
        if migration.CanApply(current) {
            log.Printf("应用迁移: %s", migration.Name())
            
            migrated, err := migration.Apply(current)
            if err != nil {
                return nil, fmt.Errorf("迁移 %s 失败: %w", migration.Name(), err)
            }
            
            current = migrated
        }
    }
    
    return current, nil
}
```

### 配置合并

智能合并多个配置：

```go
type ConfigMerger struct {
    strategy MergeStrategy
}

type MergeStrategy interface {
    Merge(base, overlay *parser.RebarConfig) (*parser.RebarConfig, error)
}

// 策略：覆盖 - 覆盖层完全替换基础部分
type OverrideStrategy struct{}

func (s OverrideStrategy) Merge(base, overlay *parser.RebarConfig) (*parser.RebarConfig, error) {
    baseTerms := createTermMap(base)
    overlayTerms := createTermMap(overlay)
    
    // 从基础术语开始
    result := make(map[string]parser.Term)
    for key, term := range baseTerms {
        result[key] = term
    }
    
    // 用覆盖层术语覆盖
    for key, term := range overlayTerms {
        result[key] = term
    }
    
    // 转换回切片
    var terms []parser.Term
    for _, term := range result {
        terms = append(terms, term)
    }
    
    return &parser.RebarConfig{Terms: terms}, nil
}

func createTermMap(config *parser.RebarConfig) map[string]parser.Term {
    termMap := make(map[string]parser.Term)
    
    for _, term := range config.Terms {
        if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) > 0 {
            if atom, ok := tuple.Elements[0].(parser.Atom); ok {
                termMap[atom.Value] = term
            }
        }
    }
    
    return termMap
}
```

## 最佳实践总结

### 1. 内存管理
- 一次处理大文件
- 使用 goroutine 进行并发处理
- 在操作之间允许垃圾回收

### 2. 错误处理
- 为更好的上下文创建自定义错误类型
- 实现部分失败的恢复策略
- 提供带有位置信息的详细错误消息

### 3. 配置管理
- 对格式更改使用迁移模式
- 实现智能合并策略
- 在处理前验证配置

### 4. 可扩展性
- 为自定义分析设计插件系统
- 使用接口进行灵活实现
- 提供全面的报告机制

### 5. 性能
- 可能时缓存解析的配置
- 对非常大的文件使用流式处理
- 在生产中分析内存使用

本高级指南为使用 Erlang Rebar 配置解析器构建强大的、生产就绪的应用程序提供了基础。
