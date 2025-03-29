// Package parser 提供解析 Erlang rebar 配置文件的功能。
// @pkg 该包用于解析 Erlang 的 rebar.config 配置文件，将其转换为 Go 的数据结构，方便 Go 程序操作和使用这些配置。
package parser

import (
	"fmt"
	"strings"
)

// RebarConfig 表示解析后的 rebar.config 文件
// @pkg RebarConfig 包含完整的配置文件内容和所有顶级配置项
// 数据样例：
// ```erlang
// {erl_opts, [debug_info]}.
// {deps, [{cowboy, "2.9.0"}]}.
// ```
// 解析后:
//
//	RebarConfig{
//	  Terms: [
//	    Tuple{Elements: [Atom{Value: "erl_opts"}, List{Elements: [Atom{Value: "debug_info"}]}]},
//	    Tuple{Elements: [Atom{Value: "deps"}, List{Elements: [Tuple{...}]}]}
//	  ]
//	}
type RebarConfig struct {
	// Raw 存储原始内容，以备参考
	Raw string
	// Terms 是配置文件中的顶级配置项列表
	Terms []Term
}

// Term 表示配置文件中的一个 Erlang 项
// @pkg Term 是一个接口，代表任何 Erlang 数据类型（原子、字符串、数字、元组、列表等）
// 所有具体类型都实现了这个接口
type Term interface {
	String() string
	// Compare 比较此 Term 与另一个 Term
	Compare(other Term) bool
}

// Tuple 表示 Erlang 元组 {elem1, elem2, ...}
// @pkg Tuple 对应 Erlang 的元组类型，由大括号包围的多个元素组成
// 数据样例: {deps, [{cowboy, "2.9.0"}]} 被解析为
// Tuple{Elements: [
//
//	Atom{Value: "deps"},
//	List{Elements: [
//	  Tuple{Elements: [Atom{Value: "cowboy"}, String{Value: "2.9.0"}]}
//	]}
//
// ]}
type Tuple struct {
	Elements []Term
}

// String 返回元组的字符串表示
// @pkg 将 Tuple 转换为字符串形式，例如 "{atom, 123}"
func (t Tuple) String() string {
	elems := make([]string, len(t.Elements))
	for i, e := range t.Elements {
		elems[i] = e.String()
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

// Compare 比较两个 Tuple 是否相等
// @pkg 比较当前 Tuple 与另一个 Term 是否相等
// 如果另一个 Term 不是 Tuple 或元素数量不同，返回 false
// 如果所有元素都相等（通过递归比较），返回 true
// 示例:
// tuple1 := Tuple{Elements: []Term{Atom{Value: "test"}, Integer{Value: 123}}}
// tuple2 := Tuple{Elements: []Term{Atom{Value: "test"}, Integer{Value: 123}}}
// tuple1.Compare(tuple2) // 返回 true
func (t Tuple) Compare(other Term) bool {
	otherTuple, ok := other.(Tuple)
	if !ok || len(t.Elements) != len(otherTuple.Elements) {
		return false
	}

	for i, element := range t.Elements {
		if !element.Compare(otherTuple.Elements[i]) {
			return false
		}
	}

	return true
}

// List 表示 Erlang 列表 [elem1, elem2, ...]
// @pkg List 对应 Erlang 的列表类型，由方括号包围的多个元素组成
// 数据样例: [debug_info, {debug_info}] 被解析为
// List{Elements: [
//
//	Atom{Value: "debug_info"},
//	Tuple{Elements: [Atom{Value: "debug_info"}]}
//
// ]}
type List struct {
	Elements []Term
}

// String 返回列表的字符串表示
// @pkg 将 List 转换为字符串形式，例如 "[atom, 123]"
func (l List) String() string {
	elems := make([]string, len(l.Elements))
	for i, e := range l.Elements {
		elems[i] = e.String()
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

// Compare 比较两个 List 是否相等
// @pkg 比较当前 List 与另一个 Term 是否相等
// 如果另一个 Term 不是 List 或元素数量不同，返回 false
// 如果所有元素都相等（通过递归比较），返回 true
// 示例:
// list1 := List{Elements: []Term{Atom{Value: "test"}, Integer{Value: 123}}}
// list2 := List{Elements: []Term{Atom{Value: "test"}, Integer{Value: 123}}}
// list1.Compare(list2) // 返回 true
func (l List) Compare(other Term) bool {
	otherList, ok := other.(List)
	if !ok || len(l.Elements) != len(otherList.Elements) {
		return false
	}

	for i, element := range l.Elements {
		if !element.Compare(otherList.Elements[i]) {
			return false
		}
	}

	return true
}

// Atom 表示 Erlang 原子（一个符号）
// @pkg Atom 对应 Erlang 的原子类型，可以是普通原子或引号包围的原子
// 数据样例:
// - 普通原子: debug_info 被解析为 Atom{Value: "debug_info", IsQuoted: false}
// - 引号原子: 'quoted-atom' 被解析为 Atom{Value: "quoted-atom", IsQuoted: true}
type Atom struct {
	Value string
	// IsQuoted 表示这个原子在原始语法中是否被引号包围
	IsQuoted bool
}

// String 返回原子的字符串表示
// @pkg 将 Atom 转换为字符串形式
// 如果原子是引号包围的，返回如 'atom-name'
// 否则直接返回原子名称，如 atom_name
func (a Atom) String() string {
	if a.IsQuoted {
		return "'" + a.Value + "'"
	}
	return a.Value
}

// Compare 比较两个 Atom 是否相等
// @pkg 比较当前 Atom 与另一个 Term 是否相等
// 如果另一个 Term 不是 Atom，返回 false
// 只比较 Value 值，忽略 IsQuoted 标志
// 示例:
// atom1 := Atom{Value: "test", IsQuoted: false}
// atom2 := Atom{Value: "test", IsQuoted: true}
// atom1.Compare(atom2) // 返回 true，因为只比较 Value
func (a Atom) Compare(other Term) bool {
	otherAtom, ok := other.(Atom)
	if !ok {
		return false
	}
	return a.Value == otherAtom.Value
}

// String 表示 Erlang 字符串
// @pkg String 对应 Erlang 的字符串类型，双引号包围的文本
// 数据样例: "hello world" 被解析为 String{Value: "hello world"}
type String struct {
	Value string
}

// String 返回字符串的字符串表示（带引号）
// @pkg 将 String 转换为字符串形式（带双引号），如 "hello world"
func (s String) String() string {
	return "\"" + s.Value + "\""
}

// Compare 比较两个 String 是否相等
// @pkg 比较当前 String 与另一个 Term 是否相等
// 如果另一个 Term 不是 String，返回 false
// 比较两个字符串的 Value 是否相等
// 示例:
// str1 := String{Value: "test"}
// str2 := String{Value: "test"}
// str1.Compare(str2) // 返回 true
func (s String) Compare(other Term) bool {
	otherString, ok := other.(String)
	if !ok {
		return false
	}
	return s.Value == otherString.Value
}

// Integer 表示 Erlang 整数
// @pkg Integer 对应 Erlang 的整数类型
// 数据样例: 123 被解析为 Integer{Value: 123}
type Integer struct {
	Value int64
}

// String 返回整数的字符串表示
// @pkg 将 Integer 转换为字符串形式，如 "123"
func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

// Compare 比较两个 Integer 是否相等
// @pkg 比较当前 Integer 与另一个 Term 是否相等
// 如果另一个 Term 不是 Integer，返回 false
// 比较两个整数的 Value 是否相等
// 示例:
// int1 := Integer{Value: 123}
// int2 := Integer{Value: 123}
// int1.Compare(int2) // 返回 true
func (i Integer) Compare(other Term) bool {
	otherInt, ok := other.(Integer)
	if !ok {
		return false
	}
	return i.Value == otherInt.Value
}

// Float 表示 Erlang 浮点数
// @pkg Float 对应 Erlang 的浮点数类型
// 数据样例: 3.14 被解析为 Float{Value: 3.14}
type Float struct {
	Value float64
}

// String 返回浮点数的字符串表示
// @pkg 将 Float 转换为字符串形式，如 "3.14"
func (f Float) String() string {
	return fmt.Sprintf("%g", f.Value)
}

// Compare 比较两个 Float 是否相等
// @pkg 比较当前 Float 与另一个 Term 是否相等
// 如果另一个 Term 不是 Float，返回 false
// 比较两个浮点数的 Value 是否相等
// 示例:
// float1 := Float{Value: 3.14}
// float2 := Float{Value: 3.14}
// float1.Compare(float2) // 返回 true
func (f Float) Compare(other Term) bool {
	otherFloat, ok := other.(Float)
	if !ok {
		return false
	}
	return f.Value == otherFloat.Value
}
