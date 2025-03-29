package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
	// 创建一个具有复杂结构的rebar.config字符串
	// 故意使其格式不规范，为了展示格式化的效果
	configStr := `{erl_opts,[debug_info,{parse_transform,lager_transform}]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.0.0"},{lager,{git,"https://github.com/erlang-lager/lager.git",{tag,"3.9.2"}}}]}.{profiles,[{dev,[{deps,[{meck,"0.9.0"}]},{erl_opts,[debug_info,{d,'DEBUG',true}]}]},{test,[{deps,[{proper,"1.3.0"}]},{erl_opts,[debug_info,nowarn_export_all]}]}]}.{relx,[{release,{my_app,"0.1.0"},[my_app,sasl]},{dev_mode,true},{include_erts,false},{extended_start_script,true}]}.{shell, [{config, "config/sys.config"}, {apps, [my_app]}]}.`

	// 解析配置
	config, err := parser.Parse(configStr)
	if err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	fmt.Println("=== 格式化输出示例 ===")
	fmt.Printf("原始配置包含 %d 个顶层条目\n", len(config.Terms))

	// 展示原始未格式化的配置
	fmt.Println("\n原始配置字符串 (压缩格式):")
	fmt.Println(configStr)

	// 使用2空格缩进格式化
	fmt.Println("\n使用2空格缩进格式化:")
	formatted2 := config.Format(2)
	fmt.Println(formatted2)

	// 使用4空格缩进格式化
	fmt.Println("\n使用4空格缩进格式化:")
	formatted4 := config.Format(4)
	fmt.Println(formatted4)

	// 使用8空格缩进格式化
	fmt.Println("\n使用8空格缩进格式化:")
	formatted8 := config.Format(8)
	fmt.Println(formatted8)

	// 演示特定结构的格式化效果
	demonstrateFormattingOfSpecificStructures()
}

// demonstrateFormattingOfSpecificStructures 演示特定数据结构的格式化
func demonstrateFormattingOfSpecificStructures() {
	fmt.Println("\n=== 特定数据结构的格式化示例 ===")

	// 1. 嵌套列表的格式化
	nestedListStr := `{nested_list, [1, [2, 3, [4, 5]], 6]}.`
	fmt.Println("\n1. 嵌套列表的格式化")
	fmt.Println("原始字符串:")
	fmt.Println(nestedListStr)

	nestedListConfig, _ := parser.Parse(nestedListStr)
	fmt.Println("\n格式化后 (4空格缩进):")
	fmt.Println(nestedListConfig.Format(4))

	// 2. 嵌套元组的格式化
	nestedTupleStr := `{nested_tuple, {level1, {level2, {level3, value}}}}.`
	fmt.Println("\n2. 嵌套元组的格式化")
	fmt.Println("原始字符串:")
	fmt.Println(nestedTupleStr)

	nestedTupleConfig, _ := parser.Parse(nestedTupleStr)
	fmt.Println("\n格式化后 (4空格缩进):")
	fmt.Println(nestedTupleConfig.Format(4))

	// 3. 包含不同类型值的列表格式化
	mixedTypesStr := `{mixed_types, [atom, "string", 123, 3.14, true, {tuple, val}, [list, items]]}.`
	fmt.Println("\n3. 混合类型列表的格式化")
	fmt.Println("原始字符串:")
	fmt.Println(mixedTypesStr)

	mixedTypesConfig, _ := parser.Parse(mixedTypesStr)
	fmt.Println("\n格式化后 (4空格缩进):")
	fmt.Println(mixedTypesConfig.Format(4))

	// 4. 带有引号的原子和特殊字符的格式化
	quotedAtomsStr := `{'quoted-atom', {'special@atom', "string with \"quotes\""}}.`
	fmt.Println("\n4. 带引号原子的格式化")
	fmt.Println("原始字符串:")
	fmt.Println(quotedAtomsStr)

	quotedAtomsConfig, _ := parser.Parse(quotedAtomsStr)
	fmt.Println("\n格式化后 (4空格缩进):")
	fmt.Println(quotedAtomsConfig.Format(4))

	// 5. 空元组和列表的格式化
	emptyStructuresStr := `{empty_tuple, {}}. {empty_list, []}.`
	fmt.Println("\n5. 空元组和列表的格式化")
	fmt.Println("原始字符串:")
	fmt.Println(emptyStructuresStr)

	emptyStructuresConfig, _ := parser.Parse(emptyStructuresStr)
	fmt.Println("\n格式化后 (4空格缩进):")
	fmt.Println(emptyStructuresConfig.Format(4))
}

// 注意：输出示例非常长，以下是部分关键输出
// 运行此示例会产生格式化前后的对比效果，展示不同缩进级别的输出:
//
// === 格式化输出示例 ===
// 原始配置包含 5 个顶层条目
//
// 原始配置字符串 (压缩格式):
// {erl_opts,[debug_info,{parse_transform,lager_transform}]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.0.0"},{lager,{git,"https://github.com/erlang-lager/lager.git",{tag,"3.9.2"}}}]}.{profiles,[{dev,[{deps,[{meck,"0.9.0"}]},{erl_opts,[debug_info,{d,'DEBUG',true}]}]},{test,[{deps,[{proper,"1.3.0"}]},{erl_opts,[debug_info,nowarn_export_all]}]}]}.{relx,[{release,{my_app,"0.1.0"},[my_app,sasl]},{dev_mode,true},{include_erts,false},{extended_start_script,true}]}.{shell, [{config, "config/sys.config"}, {apps, [my_app]}]}.
//
// 使用2空格缩进格式化:
// {erl_opts, [debug_info, {parse_transform, lager_transform}]}.
//
// {deps, [
//   {cowboy, "2.9.0"},
//   {jsx, "3.0.0"},
//   {lager, {git, "https://github.com/erlang-lager/lager.git", {tag, "3.9.2"}}}
// ]}.
//
// ...
//
// === 特定数据结构的格式化示例 ===
//
// 1. 嵌套列表的格式化
// 原始字符串:
// {nested_list, [1, [2, 3, [4, 5]], 6]}.
//
// 格式化后 (4空格缩进):
// {nested_list, [1, [2, 3, [4, 5]], 6]}.
//
// ...其他示例输出...
