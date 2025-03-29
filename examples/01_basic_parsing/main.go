package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
	// 示例1：从字符串解析rebar.config配置
	fmt.Println("=== 示例1：从字符串解析 ===")
	configStr := `
	{erl_opts, [debug_info, warnings_as_errors]}.
	{deps, [
		{cowboy, "2.9.0"},
		{jsx, "3.1.0"}
	]}.
	`

	// 使用Parse函数解析字符串
	config, err := parser.Parse(configStr)
	if err != nil {
		log.Fatalf("解析配置字符串失败: %v", err)
	}

	// 显示解析到的顶层条目数量
	fmt.Printf("从字符串解析到 %d 个顶层条目\n", len(config.Terms))

	// 输出每个顶层条目
	for i, term := range config.Terms {
		fmt.Printf("条目 %d: %s\n", i+1, term.String())
	}

	// 示例2：创建临时文件并解析
	fmt.Println("\n=== 示例2：从文件解析 ===")

	// 创建临时rebar.config文件
	tempFile := "temp_rebar.config"
	err = os.WriteFile(tempFile, []byte(`{erl_opts, [debug_info]}.
{deps, [{cowboy, "2.9.0"}]}.
{profiles, [{test, [{deps, [{meck, "0.9.0"}]}]}]}.`), 0644)

	if err != nil {
		log.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tempFile) // 程序结束时删除临时文件

	// 使用ParseFile函数解析文件
	fileConfig, err := parser.ParseFile(tempFile)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 显示解析到的顶层条目数量
	fmt.Printf("从文件解析到 %d 个顶层条目\n", len(fileConfig.Terms))

	// 输出每个顶层条目
	for i, term := range fileConfig.Terms {
		fmt.Printf("条目 %d: %s\n", i+1, term.String())
	}

	// 示例3: 从io.Reader解析
	fmt.Println("\n=== 示例3：从io.Reader解析 ===")

	// 打开之前创建的临时文件
	file, err := os.Open(tempFile)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 使用ParseReader函数从reader解析
	readerConfig, err := parser.ParseReader(file)
	if err != nil {
		log.Fatalf("从Reader解析失败: %v", err)
	}

	// 显示解析到的顶层条目数量
	fmt.Printf("从Reader解析到 %d 个顶层条目\n", len(readerConfig.Terms))
}

// 运行此示例的输出:
// === 示例1：从字符串解析 ===
// 从字符串解析到 2 个顶层条目
// 条目 1: {erl_opts, [debug_info, warnings_as_errors]}
// 条目 2: {deps, [{cowboy, "2.9.0"}, {jsx, "3.1.0"}]}
//
// === 示例2：从文件解析 ===
// 从文件解析到 3 个顶层条目
// 条目 1: {erl_opts, [debug_info]}
// 条目 2: {deps, [{cowboy, "2.9.0"}]}
// 条目 3: {profiles, [{test, [{deps, [{meck, "0.9.0"}]}]}]}
//
// === 示例3：从io.Reader解析 ===
// 从Reader解析到 3 个顶层条目
