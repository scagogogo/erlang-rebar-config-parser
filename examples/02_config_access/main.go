package main

import (
	"fmt"
	"log"

	"github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
	// 创建一个包含多种配置项的rebar.config字符串
	configStr := `
	{app_name, my_awesome_app}.
	{minimum_otp_vsn, "24.0"}.
	{erl_opts, [debug_info, {parse_transform, lager_transform}, warnings_as_errors]}.
	{deps, [
		{cowboy, "2.9.0"},
		{jsx, "3.1.0"},
		{lager, {git, "https://github.com/erlang-lager/lager.git", {tag, "3.9.2"}}}
	]}.
	{plugins, [rebar3_hex, rebar3_appup_plugin]}.
	{relx, [
		{release, {my_app, "0.1.0"}, [my_app, sasl]},
		{dev_mode, true},
		{include_erts, false},
		{extended_start_script, true}
	]}.
	{profiles, [
		{dev, [
			{deps, [{meck, "0.9.0"}]},
			{erl_opts, [debug_info, {d, 'DEBUG', true}]}
		]},
		{test, [
			{deps, [{proper, "1.3.0"}]},
			{erl_opts, [debug_info, nowarn_export_all]}
		]}
	]}.
	`

	// 解析配置
	config, err := parser.Parse(configStr)
	if err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	fmt.Println("=== 访问配置元素的示例 ===")

	// 1. 获取应用程序名称
	fmt.Println("\n1. 获取应用程序名称")
	if appName, ok := config.GetAppName(); ok {
		fmt.Printf("应用程序名称: %s\n", appName)
	} else {
		fmt.Println("未找到应用程序名称")
	}

	// 2. 获取编译选项
	fmt.Println("\n2. 获取编译选项")
	if erlOpts, ok := config.GetErlOpts(); ok {
		fmt.Printf("找到编译选项: %s\n", erlOpts[0])

		// 2.1 遍历编译选项列表
		if list, ok := erlOpts[0].(parser.List); ok {
			fmt.Println("编译选项列表中的元素:")
			for i, opt := range list.Elements {
				fmt.Printf("  %d. %s\n", i+1, opt)
			}
		}
	} else {
		fmt.Println("未找到编译选项")
	}

	// 3. 获取依赖项
	fmt.Println("\n3. 获取依赖项")
	if deps, ok := config.GetDeps(); ok {
		fmt.Printf("找到依赖项: %s\n", deps[0])

		// 3.1 遍历依赖项列表并提取依赖名称和版本
		if list, ok := deps[0].(parser.List); ok {
			fmt.Println("解析后的依赖项:")
			for _, dep := range list.Elements {
				if tuple, ok := dep.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
					if atom, ok := tuple.Elements[0].(parser.Atom); ok {
						fmt.Printf("  依赖名称: %s", atom.Value)

						// 检查版本信息类型 (字符串或Git信息)
						switch v := tuple.Elements[1].(type) {
						case parser.String:
							fmt.Printf(", 版本: %s\n", v.Value)
						case parser.Tuple:
							// Git依赖
							fmt.Printf(", 类型: git")
							if len(v.Elements) > 1 {
								if urlStr, ok := v.Elements[1].(parser.String); ok {
									fmt.Printf(", URL: %s", urlStr.Value)
								}
							}
							fmt.Println()
						default:
							fmt.Printf(", 版本信息类型: %T\n", v)
						}
					}
				}
			}
		}
	} else {
		fmt.Println("未找到依赖项")
	}

	// 4. 获取插件
	fmt.Println("\n4. 获取插件")
	if plugins, ok := config.GetPlugins(); ok {
		fmt.Printf("找到插件: %s\n", plugins[0])

		// 4.1 遍历插件列表
		if list, ok := plugins[0].(parser.List); ok {
			fmt.Println("插件列表:")
			for i, plugin := range list.Elements {
				if atom, ok := plugin.(parser.Atom); ok {
					fmt.Printf("  %d. %s\n", i+1, atom.Value)
				}
			}
		}
	} else {
		fmt.Println("未找到插件")
	}

	// 5. 获取Relx配置
	fmt.Println("\n5. 获取Relx配置")
	if relx, ok := config.GetRelxConfig(); ok {
		fmt.Printf("找到Relx配置: %s\n", relx[0])

		// 5.1 检查dev_mode设置
		if list, ok := relx[0].(parser.List); ok {
			for _, item := range list.Elements {
				if tuple, ok := item.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
					if atom, ok := tuple.Elements[0].(parser.Atom); ok && atom.Value == "dev_mode" {
						fmt.Printf("  dev_mode设置为: %s\n", tuple.Elements[1])
						break
					}
				}
			}
		}
	} else {
		fmt.Println("未找到Relx配置")
	}

	// 6. 获取配置文件
	fmt.Println("\n6. 获取配置文件")
	if profiles, ok := config.GetProfilesConfig(); ok {
		fmt.Printf("找到配置文件: %s\n", profiles[0])

		// 6.1 遍历配置文件并显示其名称
		if list, ok := profiles[0].(parser.List); ok {
			fmt.Println("配置文件:")
			for _, profile := range list.Elements {
				if tuple, ok := profile.(parser.Tuple); ok && len(tuple.Elements) >= 1 {
					if atom, ok := tuple.Elements[0].(parser.Atom); ok {
						fmt.Printf("  配置名称: %s\n", atom.Value)
					}
				}
			}
		}
	} else {
		fmt.Println("未找到配置文件")
	}

	// 7. 使用GetTerm获取任意命名的配置项
	fmt.Println("\n7. 使用GetTerm获取特定配置项")
	if term, ok := config.GetTerm("minimum_otp_vsn"); ok {
		if tuple, ok := term.(parser.Tuple); ok && len(tuple.Elements) >= 2 {
			if str, ok := tuple.Elements[1].(parser.String); ok {
				fmt.Printf("  最小OTP版本: %s\n", str.Value)
			}
		}
	} else {
		fmt.Println("未找到minimum_otp_vsn配置项")
	}

	// 8. 使用GetTupleElements获取元组元素
	fmt.Println("\n8. 使用GetTupleElements获取任意元组元素")
	if elements, ok := config.GetTupleElements("app_name"); ok {
		fmt.Printf("  app_name元组元素: %s\n", elements[0])
	} else {
		fmt.Println("未找到app_name元组元素")
	}
}

// 运行此示例的输出:
// === 访问配置元素的示例 ===
//
// 1. 获取应用程序名称
// 应用程序名称: my_awesome_app
//
// 2. 获取编译选项
// 找到编译选项: [debug_info, {parse_transform, lager_transform}, warnings_as_errors]
// 编译选项列表中的元素:
//   1. debug_info
//   2. {parse_transform, lager_transform}
//   3. warnings_as_errors
//
// 3. 获取依赖项
// 找到依赖项: [{cowboy, "2.9.0"}, {jsx, "3.1.0"}, {lager, {git, "https://github.com/erlang-lager/lager.git", {tag, "3.9.2"}}}]
// 解析后的依赖项:
//   依赖名称: cowboy, 版本: 2.9.0
//   依赖名称: jsx, 版本: 3.1.0
//   依赖名称: lager, 类型: git, URL: https://github.com/erlang-lager/lager.git
//
// 4. 获取插件
// 找到插件: [rebar3_hex, rebar3_appup_plugin]
// 插件列表:
//   1. rebar3_hex
//   2. rebar3_appup_plugin
//
// 5. 获取Relx配置
// 找到Relx配置: [{release, {my_app, "0.1.0"}, [my_app, sasl]}, {dev_mode, true}, {include_erts, false}, {extended_start_script, true}]
//   dev_mode设置为: true
//
// 6. 获取配置文件
// 找到配置文件: [{dev, [{deps, [{meck, "0.9.0"}]}, {erl_opts, [debug_info, {d, 'DEBUG', true}]}]}, {test, [{deps, [{proper, "1.3.0"}]}, {erl_opts, [debug_info, nowarn_export_all]}]}]
// 配置文件:
//   配置名称: dev
//   配置名称: test
//
// 7. 使用GetTerm获取特定配置项
//   最小OTP版本: 24.0
//
// 8. 使用GetTupleElements获取任意元组元素
//   app_name元组元素: my_awesome_app
