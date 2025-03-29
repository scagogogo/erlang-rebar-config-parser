# Erlang Rebar Config Parser 示例

本目录包含了使用 `erlang-rebar-config-parser` 库的示例代码。这些示例演示了库的各种功能，从基本的解析到高级的配置分析。

## 示例概述

所有示例均包含充分的注释，并在代码中提供了预期的输出结果。示例按复杂度递增的顺序排列：

1. **01_basic_parsing**: 演示如何从字符串、文件和io.Reader解析rebar.config
2. **02_config_access**: 演示如何访问配置中的各种元素
3. **03_pretty_printing**: 演示如何格式化rebar.config以提高可读性
4. **04_comparison**: 演示如何比较Erlang术语和配置
5. **05_complex_examples**: 包含更复杂的示例，演示如何构建实用工具

## 运行示例

每个示例都可以独立运行。进入相应的目录并使用Go命令执行：

```bash
cd examples/01_basic_parsing
go run main.go
```

## 示例详情

### 01_basic_parsing

演示了如何使用`Parse`、`ParseFile`和`ParseReader`函数来解析rebar.config文件。该示例覆盖了所有可能的解析方式。

### 02_config_access

演示了如何使用不同的访问方法（如`GetDeps`、`GetErlOpts`等）来获取配置中的元素，以及如何处理不同类型的Term（Atom、String、Integer、List、Tuple等）。

### 03_pretty_printing

演示了如何使用`Format`方法格式化rebar.config，并展示不同缩进大小的效果。它还展示了如何处理特定类型的结构，如嵌套列表、嵌套元组、混合类型列表等。

### 04_comparison

演示了如何使用`Compare`方法比较不同类型的Erlang术语，包括原子、字符串、数字、列表和元组。它还展示了如何比较复杂的嵌套结构。

### 05_complex_examples

包含一个更复杂的示例 - `rebar_config_analyzer`，它展示了如何构建一个完整的工具来分析rebar.config文件，提取依赖信息、检查编译选项、分析profiles等。

## 注意事项

- 所有示例均包含预期的输出结果作为注释，这样即使不运行代码也能理解其行为
- 示例代码中使用的所有参数都已在代码中硬编码，无需额外输入
- 部分示例会创建临时文件，但会在程序结束时自动删除 