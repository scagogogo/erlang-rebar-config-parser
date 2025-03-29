package main

import (
	"fmt"

	"github.com/scagogogo/erlang-rebar-config-parser/pkg/parser"
)

func main() {
	fmt.Println("=== Erlang术语比较示例 ===")

	// 1. 比较原子
	fmt.Println("\n1. 比较原子")
	atom1 := parser.Atom{Value: "test", IsQuoted: false}
	atom2 := parser.Atom{Value: "test", IsQuoted: true} // 相同值但一个带引号
	atom3 := parser.Atom{Value: "different", IsQuoted: false}

	fmt.Printf("atom1: %s\n", atom1)
	fmt.Printf("atom2: %s (带引号)\n", atom2)
	fmt.Printf("atom3: %s\n", atom3)

	fmt.Printf("atom1与atom2比较: %v (值相同但引号状态不同，应为true)\n", atom1.Compare(atom2))
	fmt.Printf("atom1与atom3比较: %v (值不同，应为false)\n", atom1.Compare(atom3))

	// 2. 比较字符串
	fmt.Println("\n2. 比较字符串")
	str1 := parser.String{Value: "hello"}
	str2 := parser.String{Value: "hello"}
	str3 := parser.String{Value: "world"}

	fmt.Printf("str1: %s\n", str1)
	fmt.Printf("str2: %s\n", str2)
	fmt.Printf("str3: %s\n", str3)

	fmt.Printf("str1与str2比较: %v (相同的字符串，应为true)\n", str1.Compare(str2))
	fmt.Printf("str1与str3比较: %v (不同的字符串，应为false)\n", str1.Compare(str3))

	// 3. 比较数字
	fmt.Println("\n3. 比较数字")
	int1 := parser.Integer{Value: 42}
	int2 := parser.Integer{Value: 42}
	int3 := parser.Integer{Value: 100}
	float1 := parser.Float{Value: 3.14}
	float2 := parser.Float{Value: 3.14}
	float3 := parser.Float{Value: 2.71}

	fmt.Printf("int1: %s\n", int1)
	fmt.Printf("int2: %s\n", int2)
	fmt.Printf("int3: %s\n", int3)
	fmt.Printf("float1: %s\n", float1)
	fmt.Printf("float2: %s\n", float2)
	fmt.Printf("float3: %s\n", float3)

	fmt.Printf("int1与int2比较: %v (相同的整数，应为true)\n", int1.Compare(int2))
	fmt.Printf("int1与int3比较: %v (不同的整数，应为false)\n", int1.Compare(int3))
	fmt.Printf("float1与float2比较: %v (相同的浮点数，应为true)\n", float1.Compare(float2))
	fmt.Printf("float1与float3比较: %v (不同的浮点数，应为false)\n", float1.Compare(float3))
	fmt.Printf("int1与float1比较: %v (整数与浮点数比较，应为false)\n", int1.Compare(float1))

	// 4. 比较列表
	fmt.Println("\n4. 比较列表")
	list1 := parser.List{Elements: []parser.Term{
		parser.Atom{Value: "a"},
		parser.Integer{Value: 1},
		parser.String{Value: "str"},
	}}
	list2 := parser.List{Elements: []parser.Term{
		parser.Atom{Value: "a"},
		parser.Integer{Value: 1},
		parser.String{Value: "str"},
	}}
	list3 := parser.List{Elements: []parser.Term{
		parser.Atom{Value: "a"},
		parser.Integer{Value: 2}, // 改变了一个元素
		parser.String{Value: "str"},
	}}
	list4 := parser.List{Elements: []parser.Term{
		parser.Atom{Value: "a"},
		parser.Integer{Value: 1},
	}} // 少一个元素

	emptyList1 := parser.List{Elements: []parser.Term{}}
	emptyList2 := parser.List{Elements: []parser.Term{}}

	fmt.Printf("list1: %s\n", list1)
	fmt.Printf("list2: %s\n", list2)
	fmt.Printf("list3: %s (第二个元素不同)\n", list3)
	fmt.Printf("list4: %s (少一个元素)\n", list4)
	fmt.Printf("emptyList1: %s\n", emptyList1)
	fmt.Printf("emptyList2: %s\n", emptyList2)

	fmt.Printf("list1与list2比较: %v (完全相同的列表，应为true)\n", list1.Compare(list2))
	fmt.Printf("list1与list3比较: %v (元素不同的列表，应为false)\n", list1.Compare(list3))
	fmt.Printf("list1与list4比较: %v (长度不同的列表，应为false)\n", list1.Compare(list4))
	fmt.Printf("emptyList1与emptyList2比较: %v (空列表比较，应为true)\n", emptyList1.Compare(emptyList2))

	// 5. 比较元组
	fmt.Println("\n5. 比较元组")
	tuple1 := parser.Tuple{Elements: []parser.Term{
		parser.Atom{Value: "key"},
		parser.String{Value: "value"},
	}}
	tuple2 := parser.Tuple{Elements: []parser.Term{
		parser.Atom{Value: "key"},
		parser.String{Value: "value"},
	}}
	tuple3 := parser.Tuple{Elements: []parser.Term{
		parser.Atom{Value: "key"},
		parser.String{Value: "different"},
	}}
	tuple4 := parser.Tuple{Elements: []parser.Term{
		parser.Atom{Value: "key"},
	}}

	emptyTuple1 := parser.Tuple{Elements: []parser.Term{}}
	emptyTuple2 := parser.Tuple{Elements: []parser.Term{}}

	fmt.Printf("tuple1: %s\n", tuple1)
	fmt.Printf("tuple2: %s\n", tuple2)
	fmt.Printf("tuple3: %s (第二个元素不同)\n", tuple3)
	fmt.Printf("tuple4: %s (少一个元素)\n", tuple4)
	fmt.Printf("emptyTuple1: %s\n", emptyTuple1)
	fmt.Printf("emptyTuple2: %s\n", emptyTuple2)

	fmt.Printf("tuple1与tuple2比较: %v (完全相同的元组，应为true)\n", tuple1.Compare(tuple2))
	fmt.Printf("tuple1与tuple3比较: %v (元素不同的元组，应为false)\n", tuple1.Compare(tuple3))
	fmt.Printf("tuple1与tuple4比较: %v (长度不同的元组，应为false)\n", tuple1.Compare(tuple4))
	fmt.Printf("emptyTuple1与emptyTuple2比较: %v (空元组比较，应为true)\n", emptyTuple1.Compare(emptyTuple2))

	// 6. 比较不同类型
	fmt.Println("\n6. 比较不同类型")
	fmt.Printf("atom1与str1比较: %v (原子与字符串比较，应为false)\n", atom1.Compare(str1))
	fmt.Printf("int1与str1比较: %v (整数与字符串比较，应为false)\n", int1.Compare(str1))
	fmt.Printf("list1与tuple1比较: %v (列表与元组比较，应为false)\n", list1.Compare(tuple1))
	fmt.Printf("emptyList1与emptyTuple1比较: %v (空列表与空元组比较，应为false)\n", emptyList1.Compare(emptyTuple1))

	// 7. 比较复杂的嵌套结构
	fmt.Println("\n7. 比较复杂的嵌套结构")
	// 创建两个相同的复杂结构
	complex1 := parser.Tuple{Elements: []parser.Term{
		parser.Atom{Value: "complex"},
		parser.List{Elements: []parser.Term{
			parser.Integer{Value: 1},
			parser.Tuple{Elements: []parser.Term{
				parser.Atom{Value: "nested"},
				parser.String{Value: "value"},
			}},
		}},
	}}

	complex2 := parser.Tuple{Elements: []parser.Term{
		parser.Atom{Value: "complex"},
		parser.List{Elements: []parser.Term{
			parser.Integer{Value: 1},
			parser.Tuple{Elements: []parser.Term{
				parser.Atom{Value: "nested"},
				parser.String{Value: "value"},
			}},
		}},
	}}

	// 创建一个略有不同的结构
	complex3 := parser.Tuple{Elements: []parser.Term{
		parser.Atom{Value: "complex"},
		parser.List{Elements: []parser.Term{
			parser.Integer{Value: 1},
			parser.Tuple{Elements: []parser.Term{
				parser.Atom{Value: "nested"},
				parser.String{Value: "different"}, // 改变了最深层的值
			}},
		}},
	}}

	fmt.Printf("complex1: %s\n", complex1)
	fmt.Printf("complex2: %s\n", complex2)
	fmt.Printf("complex3: %s (深层嵌套值不同)\n", complex3)

	fmt.Printf("complex1与complex2比较: %v (完全相同的复杂结构，应为true)\n", complex1.Compare(complex2))
	fmt.Printf("complex1与complex3比较: %v (深层嵌套值不同的复杂结构，应为false)\n", complex1.Compare(complex3))

	// 8. 使用Compare方法比较解析得到的配置
	fmt.Println("\n8. 比较两个解析的配置")
	config1, _ := parser.Parse(`{erl_opts, [debug_info]}. {deps, [{cowboy, "2.9.0"}]}.`)
	config2, _ := parser.Parse(`{erl_opts, [debug_info]}. {deps, [{cowboy, "2.9.0"}]}.`)
	config3, _ := parser.Parse(`{erl_opts, [debug_info]}. {deps, [{cowboy, "2.8.0"}]}.`) // 版本不同

	// 手动比较两个配置
	fmt.Printf("config1与config2比较: %v (完全相同的配置，应为true)\n", compareConfigs(config1, config2))
	fmt.Printf("config1与config3比较: %v (依赖版本不同的配置，应为false)\n", compareConfigs(config1, config3))
}

// compareConfigs 比较两个RebarConfig的内容是否相同
// 这是一个辅助函数，与parser库内部的同名函数功能类似
func compareConfigs(c1, c2 *parser.RebarConfig) bool {
	if len(c1.Terms) != len(c2.Terms) {
		return false
	}
	for i := range c1.Terms {
		if !c1.Terms[i].Compare(c2.Terms[i]) {
			return false
		}
	}
	return true
}

// 运行此示例的输出将非常长。以下是关键部分示例：
//
// === Erlang术语比较示例 ===
//
// 1. 比较原子
// atom1: test
// atom2: 'test' (带引号)
// atom3: different
// atom1与atom2比较: true (值相同但引号状态不同，应为true)
// atom1与atom3比较: false (值不同，应为false)
//
// 2. 比较字符串
// str1: "hello"
// str2: "hello"
// str3: "world"
// str1与str2比较: true (相同的字符串，应为true)
// str1与str3比较: false (不同的字符串，应为false)
//
// ...更多比较结果...
//
// 7. 比较复杂的嵌套结构
// complex1: {complex, [1, {nested, "value"}]}
// complex2: {complex, [1, {nested, "value"}]}
// complex3: {complex, [1, {nested, "different"}]} (深层嵌套值不同)
// complex1与complex2比较: true (完全相同的复杂结构，应为true)
// complex1与complex3比较: false (深层嵌套值不同的复杂结构，应为false)
//
// 8. 比较两个解析的配置
// config1与config2比较: true (完全相同的配置，应为true)
// config1与config3比较: false (依赖版本不同的配置，应为false)
