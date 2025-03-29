// Package parser 提供解析 Erlang rebar 配置文件的功能。
// @pkg 该包用于解析 Erlang 的 rebar.config 配置文件，将其转换为 Go 的数据结构，方便 Go 程序操作和使用这些配置。
package parser

import (
	"strings"
)

// processEscapes 处理字符串字面量中的转义序列
// @pkg 处理字符串和原子中的转义字符，将转义序列转换为实际字符
//
// 注意: 此实现处理以下转义序列:
// - \\\" 变成 " (双引号)
// - \\\\ 变成 \\ (反斜杠)
// - \\n 变成 \n (换行符)
// - \\r 变成 \r (回车符)
// - \\t 变成 \t (制表符)
//
// Erlang 样式的用双引号转义(")不支持。
// 对于原子名称中的单引号，请使用反斜杠转义(\\')。
//
// 输入:
//   - s: 包含转义序列的字符串
//
// 输出:
//   - string: 处理转义序列后的字符串
//
// 示例:
//
//	processEscapes("hello\\nworld") // 返回 "hello\nworld"
//	processEscapes("\\\"quoted\\\"") // 返回 "\"quoted\""
func processEscapes(s string) string {
	// 处理常见的转义序列
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\r", "\r")
	s = strings.ReplaceAll(s, "\\t", "\t")

	return s
}

// 字符分类的辅助函数

// isDigit 检查字符是否是数字
// @pkg 判断一个字符是否是数字字符 (0-9)
// 输入:
//   - ch: 要检查的字符
//
// 输出:
//   - bool: 如果是数字返回 true，否则返回 false
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isAtomStart 检查字符是否可以作为原子的起始字符
// @pkg 判断一个字符是否可以作为 Erlang 原子的首字符
// Erlang 原子必须以小写字母或下划线开头
// 输入:
//   - ch: 要检查的字符
//
// 输出:
//   - bool: 如果可以作为原子起始字符返回 true，否则返回 false
func isAtomStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || ch == '_'
}

// isAtomChar 检查字符是否可以作为原子的组成字符
// @pkg 判断一个字符是否可以作为 Erlang 原子的组成部分
// Erlang 原子可以包含字母、数字、下划线和@符号
// 输入:
//   - ch: 要检查的字符
//
// 输出:
//   - bool: 如果可以作为原子组成字符返回 true，否则返回 false
func isAtomChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '@'
}
