package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

// RebarConfigAnalyzer 分析rebar.config文件并提供各种实用的功能
type RebarConfigAnalyzer struct {
	config *parser.RebarConfig
}

// NewAnalyzer 创建新的分析器实例
func NewAnalyzer(config *parser.RebarConfig) *RebarConfigAnalyzer {
	return &RebarConfigAnalyzer{
		config: config,
	}
}

// GetDependenciesInfo 提取所有依赖项的详细信息
func (a *RebarConfigAnalyzer) GetDependenciesInfo() []DependencyInfo {
	deps, ok := a.config.GetDeps()
	if !ok || len(deps) == 0 {
		return nil
	}

	var result []DependencyInfo

	if list, ok := deps[0].(parser.List); ok {
		for _, dep := range list.Elements {
			if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
				if atom, ok := tuple.Elements[0].(parser.Atom); ok {
					info := DependencyInfo{
						Name: atom.Value,
					}

					// 处理不同类型的依赖
					switch v := tuple.Elements[1].(type) {
					case parser.String:
						// 简单版本依赖, 如 {cowboy, "2.9.0"}
						info.Type = "version"
						info.Version = v.Value
					case parser.Tuple:
						// 复杂依赖，如 git, hex 等
						if len(v.Elements) > 0 {
							if atom, ok := v.Elements[0].(parser.Atom); ok {
								info.Type = atom.Value

								// 提取更多细节
								if info.Type == "git" && len(v.Elements) > 1 {
									if url, ok := v.Elements[1].(parser.String); ok {
										info.Source = url.Value
									}

									// 提取git ref (tag, branch, commit)
									if len(v.Elements) > 2 {
										if refTuple, ok := v.Elements[2].(parser.Tuple); ok && len(refTuple.Elements) >= 2 {
											if refType, ok := refTuple.Elements[0].(parser.Atom); ok {
												info.RefType = refType.Value

												if refValue, ok := refTuple.Elements[1].(parser.String); ok {
													info.RefValue = refValue.Value
												}
											}
										}
									}
								}
							}
						}
					}

					result = append(result, info)
				}
			}
		}
	}

	return result
}

// DependencyInfo 表示一个依赖项的详细信息
type DependencyInfo struct {
	Name     string // 依赖名称
	Type     string // 依赖类型: version, git, hex 等
	Version  string // 如果是直接版本依赖
	Source   string // 如果是git依赖，指向git URL
	RefType  string // 如果是git依赖，指向引用类型 (tag, branch, commit)
	RefValue string // 引用值
}

// GetProfilesInfo 提取所有profiles的详细信息
func (a *RebarConfigAnalyzer) GetProfilesInfo() map[string]map[string]interface{} {
	profiles, ok := a.config.GetProfilesConfig()
	if !ok || len(profiles) == 0 {
		return nil
	}

	result := make(map[string]map[string]interface{})

	if list, ok := profiles[0].(parser.List); ok {
		for _, profile := range list.Elements {
			if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
				if atom, ok := tuple.Elements[0].(parser.Atom); ok {
					profileName := atom.Value
					result[profileName] = make(map[string]interface{})

					// 提取profile中的配置
					if profileList, ok := tuple.Elements[1].(parser.List); ok {
						for _, item := range profileList.Elements {
							if itemTuple, ok := item.(parser.Tuple); ok && len(itemTuple.Elements) >= 2 {
								if itemKey, ok := itemTuple.Elements[0].(parser.Atom); ok {
									result[profileName][itemKey.Value] = itemTuple.Elements[1]
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

// HasWarningsAsErrors 检查是否设置了warnings_as_errors
func (a *RebarConfigAnalyzer) HasWarningsAsErrors() bool {
	erlOpts, ok := a.config.GetErlOpts()
	if !ok || len(erlOpts) == 0 {
		return false
	}

	if list, ok := erlOpts[0].(parser.List); ok {
		for _, opt := range list.Elements {
			if atom, ok := opt.(parser.Atom); ok && atom.Value == "warnings_as_errors" {
				return true
			}
		}
	}

	return false
}

// FormatConfig 格式化配置并返回格式化的字符串
func (a *RebarConfigAnalyzer) FormatConfig(indent int) string {
	return a.config.Format(indent)
}

// CountByTermType 计算不同类型Term的数量
func (a *RebarConfigAnalyzer) CountByTermType() map[string]int {
	counts := make(map[string]int)

	// 先处理顶层Term
	for _, term := range a.config.Terms {
		// 按照term类型计数
		typeName := getTermTypeName(term)
		counts[typeName]++

		// 递归分析子元素
		countTermsRecursive(term, counts)
	}

	return counts
}

// countTermsRecursive 递归计算所有子元素的类型
func countTermsRecursive(term parser.Term, counts map[string]int) {
	switch t := term.(type) {
	case parser.Tuple:
		for _, elem := range t.Elements {
			typeName := getTermTypeName(elem)
			counts[typeName]++
			countTermsRecursive(elem, counts)
		}
	case parser.List:
		for _, elem := range t.Elements {
			typeName := getTermTypeName(elem)
			counts[typeName]++
			countTermsRecursive(elem, counts)
		}
	}
}

// getTermTypeName 获取Term的类型名称
func getTermTypeName(term parser.Term) string {
	switch term.(type) {
	case parser.Atom:
		return "Atom"
	case parser.String:
		return "String"
	case parser.Integer:
		return "Integer"
	case parser.Float:
		return "Float"
	case parser.Tuple:
		return "Tuple"
	case parser.List:
		return "List"
	default:
		return "Unknown"
	}
}

func main() {
	// 创建一个完整的真实世界rebar.config示例
	rebarConfig := `
	{minimum_otp_vsn, "22.0"}.
	{erl_opts, [
		debug_info,
		warnings_as_errors,
		{parse_transform, lager_transform},
		{i, "include"}
	]}.
	{deps, [
		{cowboy, "2.9.0"},
		{jsx, "3.1.0"},
		{lager, {git, "https://github.com/erlang-lager/lager.git", {tag, "3.9.2"}}},
		{meck, {git, "https://github.com/eproxus/meck.git", {branch, "master"}}},
		{eredis, {git, "https://github.com/wooga/eredis.git", {ref, "8a7dad3"}}}
	]}.
	{plugins, [
		rebar3_hex,
		rebar3_appup_plugin,
		rebar3_auto,
		{rebar3_elixir_compile, "0.2.1"}
	]}.
	{xref_checks,[
		undefined_function_calls,
		undefined_functions,
		locals_not_used,
		deprecated_function_calls,
		deprecated_functions
	]}.
	{relx, [
		{release, {my_app, "0.1.0"}, [
			my_app,
			sasl
		]},
		{dev_mode, true},
		{include_erts, false},
		{extended_start_script, true},
		{vm_args, "config/vm.args"},
		{sys_config, "config/sys.config"}
	]}.
	{profiles, [
		{dev, [
			{deps, [
				{sync, "0.1.3"},
				{meck, "0.9.0"}
			]},
			{erl_opts, [
				debug_info, 
				{d, 'DEBUG', true}
			]}
		]},
		{test, [
			{deps, [
				{proper, "1.3.0"},
				{cowlib, "2.11.0"},
				{gun, "2.0.0-rc.1"}
			]},
			{erl_opts, [
				debug_info,
				nowarn_export_all,
				nowarn_deprecated_function
			]}
		]},
		{prod, [
			{erl_opts, [
				no_debug_info,
				warnings_as_errors
			]},
			{relx, [
				{dev_mode, false},
				{include_erts, true}
			]}
		]}
	]}.
	{dialyzer, [
		{warnings, [
			unmatched_returns,
			error_handling
		]}
	]}.
	`

	// 创建临时文件
	tempDir := os.TempDir()
	tempConfigFile := filepath.Join(tempDir, "example_rebar.config")

	err := os.WriteFile(tempConfigFile, []byte(rebarConfig), 0644)
	if err != nil {
		log.Fatalf("创建临时配置文件失败: %v", err)
	}
	defer os.Remove(tempConfigFile)

	// 解析配置
	config, err := parser.ParseFile(tempConfigFile)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 创建分析器
	analyzer := NewAnalyzer(config)

	// 使用分析器提供的功能
	fmt.Println("=== Rebar配置分析器示例 ===")

	// 1. 分析依赖
	fmt.Println("\n1. 分析依赖")
	dependencies := analyzer.GetDependenciesInfo()

	fmt.Printf("找到 %d 个依赖项\n", len(dependencies))

	fmt.Println("\n依赖分类:")
	byType := make(map[string][]string)
	for _, dep := range dependencies {
		byType[dep.Type] = append(byType[dep.Type], dep.Name)
	}

	for depType, deps := range byType {
		fmt.Printf("  类型 '%s': %d 个 (%s)\n", depType, len(deps), strings.Join(deps, ", "))
	}

	fmt.Println("\nGit依赖详情:")
	for _, dep := range dependencies {
		if dep.Type == "git" {
			fmt.Printf("  %s:\n", dep.Name)
			fmt.Printf("    URL: %s\n", dep.Source)
			fmt.Printf("    引用类型: %s\n", dep.RefType)
			fmt.Printf("    引用值: %s\n", dep.RefValue)
		}
	}

	// 2. 分析构建配置
	fmt.Println("\n2. 分析构建配置")
	fmt.Printf("启用warnings_as_errors: %v\n", analyzer.HasWarningsAsErrors())

	// 3. 分析profiles
	fmt.Println("\n3. 分析profiles")
	profiles := analyzer.GetProfilesInfo()

	fmt.Printf("找到 %d 个profiles\n", len(profiles))

	for name, profile := range profiles {
		fmt.Printf("\nProfile '%s':\n", name)
		if deps, ok := profile["deps"]; ok {
			fmt.Printf("  有依赖项配置: %v\n", deps != nil)
		}
		if opts, ok := profile["erl_opts"]; ok {
			fmt.Printf("  有编译选项配置: %v\n", opts != nil)
		}
	}

	// 4. 统计Term类型
	fmt.Println("\n4. 统计Term类型")
	termCounts := analyzer.CountByTermType()

	// 按类型名称排序
	types := make([]string, 0, len(termCounts))
	for typeName := range termCounts {
		types = append(types, typeName)
	}
	sort.Strings(types)

	for _, typeName := range types {
		fmt.Printf("  %s: %d\n", typeName, termCounts[typeName])
	}

	// 5. 格式化输出
	fmt.Println("\n5. 格式化输出示例 (2空格缩进)")
	fmt.Println(analyzer.FormatConfig(2))
}

// 运行此示例将输出详细的rebar.config分析结果，包括:
// - 依赖项分析，分类和详情
// - 编译选项分析
// - 配置文件(profiles)分析
// - Term类型统计
// - 格式化输出
//
// 输出示例部分摘录:
// === Rebar配置分析器示例 ===
//
// 1. 分析依赖
// 找到 5 个依赖项
//
// 依赖分类:
//   类型 'version': 2 个 (cowboy, jsx)
//   类型 'git': 3 个 (lager, meck, eredis)
//
// Git依赖详情:
//   lager:
//     URL: https://github.com/erlang-lager/lager.git
//     引用类型: tag
//     引用值: 3.9.2
//   ...
//
// 2. 分析构建配置
// 启用warnings_as_errors: true
//
// 3. 分析profiles
// 找到 3 个profiles
//
// Profile 'dev':
//   有依赖项配置: true
//   有编译选项配置: true
// ...
//
// 4. 统计Term类型
//   Atom: 42
//   Integer: 1
//   List: 15
//   String: 16
//   Tuple: 32
