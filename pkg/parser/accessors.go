// Package parser 提供解析 Erlang rebar 配置文件的功能。
// @pkg 该包用于解析 Erlang 的 rebar.config 配置文件，将其转换为 Go 的数据结构，方便 Go 程序操作和使用这些配置。
package parser

// GetTerm 根据名称获取配置中的特定项
// @pkg 通过名称检索配置中的特定顶级项
// 输入:
//   - name: 要查找的项名称
//
// 输出:
//   - Term: 找到的项
//   - bool: 是否找到该项
//
// 示例:
//
//	term, ok := config.GetTerm("deps")
//	if ok {
//	  fmt.Println("找到 deps 配置项:", term)
//	}
func (c *RebarConfig) GetTerm(name string) (Term, bool) {
	for _, term := range c.Terms {
		if tuple, ok := term.(Tuple); ok && len(tuple.Elements) >= 1 {
			if atom, ok := tuple.Elements[0].(Atom); ok && atom.Value == name {
				return term, true
			}
		}
	}
	return nil, false
}

// GetTupleElements 获取命名元组的元素（在 rebar 配置中很常见）
// @pkg 获取指定命名元组中的元素列表，不包括名称本身
// 输入:
//   - name: 元组的名称
//
// 输出:
//   - []Term: 元组元素列表（不包括名称）
//   - bool: 是否找到并成功获取元素
//
// 示例:
//
//	elements, ok := config.GetTupleElements("deps")
//	if ok {
//	  fmt.Printf("deps 有 %d 个依赖项\n", len(elements))
//	}
func (c *RebarConfig) GetTupleElements(name string) ([]Term, bool) {
	term, ok := c.GetTerm(name)
	if !ok {
		return nil, false
	}

	if tuple, ok := term.(Tuple); ok && len(tuple.Elements) > 1 {
		return tuple.Elements[1:], true
	}

	return nil, false
}

// GetDeps 获取 deps 配置（如果存在）
// @pkg 获取项目依赖配置列表
// 输出:
//   - []Term: 依赖项列表
//   - bool: 是否找到 deps 配置
//
// 示例:
//
//	deps, ok := config.GetDeps()
//	if ok {
//	  for _, dep := range deps {
//	    if depTuple, ok := dep.(Tuple); ok {
//	      fmt.Println("依赖项:", depTuple)
//	    }
//	  }
//	}
//
// 数据样例:
// 原始配置: {deps, [{cowboy, "2.9.0"}, {jsx, "3.1.0"}]}.
// 返回: []Term{Tuple{...cowboy...}, Tuple{...jsx...}}, true
func (c *RebarConfig) GetDeps() ([]Term, bool) {
	return c.GetTupleElements("deps")
}

// GetErlOpts 获取 erl_opts 配置（如果存在）
// @pkg 获取 Erlang 编译选项列表
// 输出:
//   - []Term: 编译选项列表
//   - bool: 是否找到 erl_opts 配置
//
// 示例:
//
//	opts, ok := config.GetErlOpts()
//	if ok {
//	  for _, opt := range opts {
//	    fmt.Println("编译选项:", opt)
//	  }
//	}
//
// 数据样例:
// 原始配置: {erl_opts, [debug_info, {parse_transform, lager_transform}]}.
// 返回: []Term{Atom{Value: "debug_info"}, Tuple{...}}, true
func (c *RebarConfig) GetErlOpts() ([]Term, bool) {
	return c.GetTupleElements("erl_opts")
}

// GetAppName 获取应用名称（如果存在）
// @pkg 获取应用程序名称
// 输出:
//   - string: 应用程序名称
//   - bool: 是否找到应用程序名称
//
// 示例:
//
//	name, ok := config.GetAppName()
//	if ok {
//	  fmt.Println("应用名称:", name)
//	}
//
// 数据样例:
// 原始配置: {app_name, "my_app"}.
// 返回: "my_app", true
func (c *RebarConfig) GetAppName() (string, bool) {
	elements, ok := c.GetTupleElements("app_name")
	if !ok || len(elements) == 0 {
		return "", false
	}

	if str, ok := elements[0].(String); ok {
		return str.Value, true
	}

	if atom, ok := elements[0].(Atom); ok {
		return atom.Value, true
	}

	return "", false
}

// GetPlugins 获取 plugins 配置（如果存在）
// @pkg 获取插件列表
// 输出:
//   - []Term: 插件列表
//   - bool: 是否找到 plugins 配置
//
// 示例:
//
//	plugins, ok := config.GetPlugins()
//	if ok {
//	  for _, plugin := range plugins {
//	    if atom, ok := plugin.(Atom); ok {
//	      fmt.Println("插件:", atom.Value)
//	    }
//	  }
//	}
//
// 数据样例:
// 原始配置: {plugins, [rebar3_hex, rebar3_auto]}.
// 返回: []Term{Atom{Value: "rebar3_hex"}, Atom{Value: "rebar3_auto"}}, true
func (c *RebarConfig) GetPlugins() ([]Term, bool) {
	return c.GetTupleElements("plugins")
}

// GetRelxConfig 获取 relx 配置（如果存在）
// @pkg 获取 relx 发布配置
// 输出:
//   - []Term: relx 配置项列表
//   - bool: 是否找到 relx 配置
//
// 示例:
//
//	relx, ok := config.GetRelxConfig()
//	if ok {
//	  fmt.Println("Relx 配置:", relx)
//	}
func (c *RebarConfig) GetRelxConfig() ([]Term, bool) {
	return c.GetTupleElements("relx")
}

// GetProfilesConfig 获取 profiles 配置（如果存在）
// @pkg 获取项目的不同环境配置
// 输出:
//   - []Term: profiles 配置项列表
//   - bool: 是否找到 profiles 配置
//
// 示例:
//
//	profiles, ok := config.GetProfilesConfig()
//	if ok {
//	  fmt.Println("环境配置:", profiles)
//	}
//
// 数据样例:
// 原始配置: {profiles, [{dev, [...]}, {prod, [...]}]}.
// 返回: []Term{Tuple{...dev...}, Tuple{...prod...}}, true
func (c *RebarConfig) GetProfilesConfig() ([]Term, bool) {
	return c.GetTupleElements("profiles")
}
