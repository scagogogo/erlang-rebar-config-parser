// Package parser 提供解析 Erlang rebar 配置文件的功能。
// @pkg 该包用于解析 Erlang 的 rebar.config 配置文件，将其转换为 Go 的数据结构，方便 Go 程序操作和使用这些配置。
package parser

import (
	"fmt"
	"strings"
)

// Format 返回配置的格式化字符串表示
// @pkg 根据指定的缩进返回格式化的配置字符串，美化配置文件的输出
// 输入:
//   - indent: 缩进空格数量，如 2 或 4
//
// 输出:
//   - string: 格式化后的配置字符串
//
// 示例:
//
//	config, _ := parser.ParseFile("./rebar.config")
//	formatted := config.Format(4) // 使用4个空格缩进
//	fmt.Println(formatted)
//
// 数据样例:
// 输入配置:
//
//	{deps,[{cowboy,"2.9.0"},{jsx,"3.1.0"}]}.
//
// 格式化输出(indent=2):
//
//	{deps, [
//	  {cowboy, "2.9.0"},
//	  {jsx, "3.1.0"}
//	]}.
func (c *RebarConfig) Format(indent int) string {
	var result strings.Builder

	for i, term := range c.Terms {
		result.WriteString(formatTerm(term, 0, indent))
		result.WriteString(".")

		if i < len(c.Terms)-1 {
			result.WriteString("\n\n")
		} else {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// formatTerm 格式化单个 Term，加上适当的缩进
// @pkg 根据缩进级别格式化单个 Term
// 输入:
//   - term: 要格式化的项
//   - level: 当前缩进级别
//   - spaces: 每级缩进的空格数
//
// 输出:
//   - string: 格式化后的字符串
//
// 递归处理复杂的嵌套结构，对不同类型的 Term 应用不同的格式化规则
func formatTerm(term Term, level, spaces int) string {
	indent := strings.Repeat(" ", level*spaces)

	switch t := term.(type) {
	case Atom:
		if t.IsQuoted {
			return "'" + t.Value + "'"
		}
		return t.Value

	case String:
		return fmt.Sprintf("%q", t.Value)

	case Integer:
		return fmt.Sprintf("%d", t.Value)

	case Float:
		return fmt.Sprintf("%g", t.Value)

	case Tuple:
		if len(t.Elements) == 0 {
			return "{}"
		}

		// 针对 rebar.config 中常见模式的特殊处理
		if len(t.Elements) >= 2 {
			if atom, ok := t.Elements[0].(Atom); ok {
				// 对于 {key, value} 形式的简单元组
				if isSimpleTerm(t.Elements[1]) {
					elems := make([]string, len(t.Elements))
					for i, e := range t.Elements {
						elems[i] = formatTerm(e, 0, spaces)
					}
					return "{" + strings.Join(elems, ", ") + "}"
				}

				// 对于 {key, [list_items]} 或 {key, {nested_tuple}} 形式的元组
				var result strings.Builder
				result.WriteString("{")
				result.WriteString(atom.String())
				result.WriteString(", ")

				for i := 1; i < len(t.Elements); i++ {
					if i > 1 {
						result.WriteString(", ")
					}
					// 对其余元素使用增加的缩进级别
					result.WriteString(formatTerm(t.Elements[i], level+1, spaces))
				}

				result.WriteString("}")
				return result.String()
			}
		}

		// 元组的默认处理方式
		var result strings.Builder
		result.WriteString("{\n")

		innerIndent := strings.Repeat(" ", (level+1)*spaces)
		for i, elem := range t.Elements {
			result.WriteString(innerIndent)
			result.WriteString(formatTerm(elem, level+1, spaces))

			if i < len(t.Elements)-1 {
				result.WriteString(",\n")
			} else {
				result.WriteString("\n")
			}
		}

		result.WriteString(indent)
		result.WriteString("}")
		return result.String()

	case List:
		if len(t.Elements) == 0 {
			return "[]"
		}

		// 对于只包含简单项的短列表，保持在一行
		if len(t.Elements) <= 3 && allSimpleTerms(t.Elements) {
			elems := make([]string, len(t.Elements))
			for i, e := range t.Elements {
				elems[i] = formatTerm(e, 0, spaces)
			}
			return "[" + strings.Join(elems, ", ") + "]"
		}

		// 其他情况使用合适的缩进格式化
		var result strings.Builder
		result.WriteString("[\n")

		innerIndent := strings.Repeat(" ", (level+1)*spaces)
		for i, elem := range t.Elements {
			result.WriteString(innerIndent)
			result.WriteString(formatTerm(elem, level+1, spaces))

			if i < len(t.Elements)-1 {
				result.WriteString(",\n")
			} else {
				result.WriteString("\n")
			}
		}

		result.WriteString(indent)
		result.WriteString("]")
		return result.String()

	default:
		return "UNKNOWN_TERM"
	}
}

// isSimpleTerm 检查一个 Term 是否是"简单的"（可以格式化在单行上）
// @pkg 判断一个 Term 是否足够简单可以在一行内显示
// 简单 Term 包括：
// - 原子、字符串、整数、浮点数
// - 元素数量少且所有元素都是简单 Term 的列表
// - 元素数量少且所有元素都是简单 Term 的元组
// 输入:
//   - term: 要检查的 Term
//
// 输出:
//   - bool: 如果是简单 Term 返回 true，否则返回 false
func isSimpleTerm(term Term) bool {
	switch t := term.(type) {
	case Atom, String, Integer, Float:
		return true
	case List:
		return len(t.Elements) <= 3 && allSimpleTerms(t.Elements)
	case Tuple:
		return len(t.Elements) <= 2 && allSimpleTerms(t.Elements)
	default:
		return false
	}
}

// allSimpleTerms 检查切片中所有 Term 是否都是"简单的"
// @pkg 检查一个 Term 列表中是否所有元素都是简单 Term
// 输入:
//   - terms: 要检查的 Term 列表
//
// 输出:
//   - bool: 如果所有元素都是简单 Term 返回 true，否则返回 false
func allSimpleTerms(terms []Term) bool {
	for _, term := range terms {
		if !isSimpleTerm(term) {
			return false
		}
	}
	return true
}
