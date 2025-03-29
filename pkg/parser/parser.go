// Package parser 提供解析 Erlang rebar 配置文件的功能。
// @pkg 该包用于解析 Erlang 的 rebar.config 配置文件，将其转换为 Go 的数据结构，方便 Go 程序操作和使用这些配置。
package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Parser 表示 Erlang 项解析器
// @pkg Parser 是一个用于解析 Erlang 项的解析器，跟踪输入字符串的位置、行号和列号
type Parser struct {
	input    string // 输入字符串
	position int    // 当前位置
	line     int    // 当前行号
	column   int    // 当前列号
}

// NewParser 创建一个新的 Parser 实例
// @pkg 根据输入字符串创建一个新的解析器实例
// 输入:
//   - input: 要解析的字符串
//
// 输出:
//   - *Parser: 新的解析器实例
//
// 示例:
//
//	parser := NewParser("{deps, [{cowboy, \"2.9.0\"}]}.")
func NewParser(input string) *Parser {
	return &Parser{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
	}
}

// ParseFile 解析指定路径的 rebar.config 文件
// @pkg 从文件系统读取并解析 rebar.config 文件
// 输入:
//   - path: 文件路径，如 "./rebar.config"
//
// 输出:
//   - *RebarConfig: 解析后的配置对象
//   - error: 解析过程中的错误，如文件不存在或解析失败
//
// 示例:
//
//	config, err := parser.ParseFile("./rebar.config")
//	if err != nil {
//	  log.Fatalf("解析失败: %v", err)
//	}
//	fmt.Printf("配置项数量: %d\n", len(config.Terms))
func ParseFile(path string) (*RebarConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return Parse(string(content))
}

// ParseReader 从给定的 reader 解析 rebar.config
// @pkg 从 io.Reader 接口（如文件、HTTP 响应等）读取并解析 rebar.config
// 输入:
//   - r: io.Reader 接口，提供配置内容
//
// 输出:
//   - *RebarConfig: 解析后的配置对象
//   - error: 解析过程中的错误
//
// 示例:
//
//	file, _ := os.Open("./rebar.config")
//	config, err := parser.ParseReader(file)
//	if err != nil {
//	  log.Fatalf("解析失败: %v", err)
//	}
func ParseReader(r io.Reader) (*RebarConfig, error) {
	var builder strings.Builder
	reader := bufio.NewReader(r)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				builder.WriteString(line)
				break
			}
			return nil, fmt.Errorf("error reading input: %w", err)
		}
		builder.WriteString(line)
	}

	return Parse(builder.String())
}

// Parse 将输入字符串解析为 rebar.config 文件
// @pkg 解析包含 Erlang 项的字符串为 RebarConfig 对象
// 输入:
//   - input: 包含 Erlang 配置的字符串
//
// 输出:
//   - *RebarConfig: 解析后的配置对象
//   - error: 解析过程中的错误
//
// 示例:
//
//	configStr := `{erl_opts, [debug_info]}.
//	             {deps, [{cowboy, "2.9.0"}]}.`
//	config, err := parser.Parse(configStr)
//	if err != nil {
//	  log.Fatalf("解析失败: %v", err)
//	}
//
//	// 访问解析后的配置
//	deps, ok := config.GetDeps()
//	if ok {
//	  fmt.Println("依赖项:", deps)
//	}
func Parse(input string) (*RebarConfig, error) {
	parser := NewParser(input)
	terms, err := parser.parseTerms()
	if err != nil {
		return nil, err
	}

	return &RebarConfig{
		Raw:   input,
		Terms: terms,
	}, nil
}

// parseTerms 解析输入中的所有项
// @pkg 解析输入字符串中的所有顶级 Erlang 项
// 每个项以点号(.)结尾，跳过注释和空白字符
// 输出:
//   - []Term: 解析出的所有项
//   - error: 解析过程中的错误
func (p *Parser) parseTerms() ([]Term, error) {
	terms := []Term{}

	for p.position < len(p.input) {
		p.skipWhitespace()

		if p.position >= len(p.input) {
			break
		}

		// 跳过注释
		if p.currentChar() == '%' {
			p.skipToEndOfLine()
			continue
		}

		term, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		terms = append(terms, term)

		// 跳过末尾的点号
		p.skipWhitespace()
		if p.position < len(p.input) && p.currentChar() == '.' {
			p.advance()
		} else {
			return nil, p.errorAt("expected '.' after term")
		}
	}

	return terms, nil
}

// parseTerm 解析单个 Erlang 项
// @pkg 根据当前字符解析不同类型的 Erlang 项
// 根据起始字符决定解析方式:
// - '{' 解析为元组
// - '[' 解析为列表
// - '"' 解析为字符串
// - '\” 解析为带引号的原子
// - '-' 或数字解析为数字
// - 其他字母开头解析为原子
// 输出:
//   - Term: 解析出的项
//   - error: 解析过程中的错误
func (p *Parser) parseTerm() (Term, error) {
	p.skipWhitespace()

	if p.position >= len(p.input) {
		return nil, p.errorAt("unexpected end of input")
	}

	switch p.currentChar() {
	case '{':
		return p.parseTuple()
	case '[':
		return p.parseList()
	case '"':
		return p.parseString()
	case '\'':
		return p.parseQuotedAtom()
	case '-':
		// 可能是负数
		return p.parseNumber()
	default:
		if isDigit(p.currentChar()) {
			return p.parseNumber()
		} else if isAtomStart(p.currentChar()) {
			return p.parseAtom()
		}
		return nil, p.errorAt(fmt.Sprintf("unexpected character: %c", p.currentChar()))
	}
}

// parseTuple 解析 Erlang 元组: {elem1, elem2, ...}
// @pkg 解析以 '{' 开始的 Erlang 元组
// 元组格式为 {元素1, 元素2, ...}，元素间用逗号分隔
// 输出:
//   - Term: 解析出的元组
//   - error: 解析过程中的错误
//
// 数据样例:
// "{deps, [{cowboy, \"2.9.0\"}]}" 被解析为
// Tuple{Elements: [Atom{Value: "deps"}, List{...}]}
func (p *Parser) parseTuple() (Term, error) {
	// 跳过 '{'
	p.advance()

	elements := []Term{}

	p.skipWhitespace()
	if p.currentChar() == '}' {
		p.advance()
		return Tuple{Elements: elements}, nil
	}

	for {
		element, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)

		p.skipWhitespace()
		if p.currentChar() == '}' {
			p.advance()
			return Tuple{Elements: elements}, nil
		}

		if p.currentChar() != ',' {
			return nil, p.errorAt("expected ',' or '}' in tuple")
		}

		// 跳过 ','
		p.advance()
		p.skipWhitespace()
	}
}

// parseList 解析 Erlang 列表: [elem1, elem2, ...]
// @pkg 解析以 '[' 开始的 Erlang 列表
// 列表格式为 [元素1, 元素2, ...]，元素间用逗号分隔
// 输出:
//   - Term: 解析出的列表
//   - error: 解析过程中的错误
//
// 数据样例:
// "[debug_info, {parse_transform, lager_transform}]" 被解析为
// List{Elements: [Atom{Value: "debug_info"}, Tuple{...}]}
func (p *Parser) parseList() (Term, error) {
	// 跳过 '['
	p.advance()

	elements := []Term{}

	p.skipWhitespace()
	if p.currentChar() == ']' {
		p.advance()
		return List{Elements: elements}, nil
	}

	for {
		element, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)

		p.skipWhitespace()
		if p.currentChar() == ']' {
			p.advance()
			return List{Elements: elements}, nil
		}

		if p.currentChar() != ',' {
			return nil, p.errorAt("expected ',' or ']' in list")
		}

		// 跳过 ','
		p.advance()
		p.skipWhitespace()
	}
}

// parseString 解析 Erlang 字符串（双引号包围）
// @pkg 解析以 '"' 开始的 Erlang 字符串，处理转义序列
// 输出:
//   - Term: 解析出的字符串
//   - error: 解析过程中的错误
//
// 数据样例:
// "\"hello world\"" 被解析为 String{Value: "hello world"}
func (p *Parser) parseString() (Term, error) {
	// 跳过开头的引号
	p.advance()

	startPos := p.position
	for p.position < len(p.input) && p.currentChar() != '"' {
		// 处理转义序列
		if p.currentChar() == '\\' {
			p.advance()
			if p.position >= len(p.input) {
				return nil, p.errorAt("unterminated string literal")
			}
		}
		p.advance()
	}

	if p.position >= len(p.input) {
		return nil, p.errorAt("unterminated string literal")
	}

	value := p.input[startPos:p.position]
	// 处理转义序列
	value = processEscapes(value)

	// 跳过结束引号
	p.advance()

	return String{Value: value}, nil
}

// parseQuotedAtom 解析带引号的原子 ('atom')
// @pkg 解析以 '\” 开始的 Erlang 带引号原子，处理转义序列
// 输出:
//   - Term: 解析出的原子
//   - error: 解析过程中的错误
//
// 数据样例:
// "'quoted-atom'" 被解析为 Atom{Value: "quoted-atom", IsQuoted: true}
func (p *Parser) parseQuotedAtom() (Term, error) {
	// 跳过开头的引号
	p.advance()

	startPos := p.position
	for p.position < len(p.input) && p.currentChar() != '\'' {
		// 处理转义序列
		if p.currentChar() == '\\' {
			p.advance()
			if p.position >= len(p.input) {
				return nil, p.errorAt("unterminated atom literal")
			}
		}
		p.advance()
	}

	if p.position >= len(p.input) {
		return nil, p.errorAt("unterminated atom literal")
	}

	value := p.input[startPos:p.position]
	// 处理转义序列
	value = processEscapes(value)

	// 跳过结束引号
	p.advance()

	return Atom{Value: value, IsQuoted: true}, nil
}

// parseAtom 解析 Erlang 原子（未带引号的符号）
// @pkg 解析 Erlang 未带引号的原子，原子必须以小写字母或下划线开头
// 输出:
//   - Term: 解析出的原子
//   - error: 解析过程中的错误
//
// 数据样例:
// "debug_info" 被解析为 Atom{Value: "debug_info", IsQuoted: false}
func (p *Parser) parseAtom() (Term, error) {
	startPos := p.position

	// 首字符已经检查为有效的原子起始字符
	p.advance()

	// 读取剩余的原子字符
	for p.position < len(p.input) && isAtomChar(p.currentChar()) {
		p.advance()
	}

	if p.position > startPos {
		value := p.input[startPos:p.position]
		return Atom{Value: value, IsQuoted: false}, nil
	}

	return nil, p.errorAt("invalid atom")
}

// parseNumber 解析 Erlang 数字（整数或浮点数）
// @pkg 解析 Erlang 数字，支持整数和浮点数（包括科学计数法）
// 输出:
//   - Term: 解析出的数字（Integer 或 Float）
//   - error: 解析过程中的错误
//
// 数据样例:
// - "123" 被解析为 Integer{Value: 123}
// - "-42" 被解析为 Integer{Value: -42}
// - "3.14" 被解析为 Float{Value: 3.14}
// - "-2.5e-3" 被解析为 Float{Value: -0.0025}
func (p *Parser) parseNumber() (Term, error) {
	startPos := p.position

	// 处理负号
	if p.currentChar() == '-' {
		p.advance()
	}

	// 读取小数点前的数字
	hasDigits := false
	for p.position < len(p.input) && isDigit(p.currentChar()) {
		hasDigits = true
		p.advance()
	}

	// 检查是否是浮点数
	isFloat := false
	if p.position < len(p.input) && p.currentChar() == '.' {
		isFloat = true
		p.advance()

		// 读取小数点后的数字
		hasDecimalDigits := false
		for p.position < len(p.input) && isDigit(p.currentChar()) {
			hasDecimalDigits = true
			p.advance()
		}

		if !hasDecimalDigits {
			return nil, p.errorAt("expected digits after decimal point")
		}
	}

	// 处理科学计数法
	if p.position < len(p.input) && (p.currentChar() == 'e' || p.currentChar() == 'E') {
		isFloat = true
		p.advance()

		// 处理指数中的符号
		if p.position < len(p.input) && (p.currentChar() == '+' || p.currentChar() == '-') {
			p.advance()
		}

		// 读取指数数字
		hasExpDigits := false
		for p.position < len(p.input) && isDigit(p.currentChar()) {
			hasExpDigits = true
			p.advance()
		}

		if !hasExpDigits {
			return nil, p.errorAt("expected digits in exponent")
		}
	}

	if !hasDigits {
		return nil, p.errorAt("expected digits in number")
	}

	value := p.input[startPos:p.position]

	if isFloat {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, p.errorAt(fmt.Sprintf("invalid float: %s", value))
		}
		return Float{Value: f}, nil
	} else {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, p.errorAt(fmt.Sprintf("invalid integer: %s", value))
		}
		return Integer{Value: i}, nil
	}
}

// Helper methods for the parser
// 解析器的辅助方法

// currentChar 返回当前字符
// @pkg 返回当前位置的字符，如果已到达输入末尾则返回 0
func (p *Parser) currentChar() byte {
	if p.position >= len(p.input) {
		return 0
	}
	return p.input[p.position]
}

// advance 前进到下一个字符
// @pkg 将位置前进一个字符，并更新行号和列号
func (p *Parser) advance() {
	if p.position < len(p.input) {
		if p.input[p.position] == '\n' {
			p.line++
			p.column = 1
		} else {
			p.column++
		}
		p.position++
	}
}

// skipWhitespace 跳过空白字符
// @pkg 跳过所有空格、制表符、换行符和回车符
func (p *Parser) skipWhitespace() {
	for p.position < len(p.input) && (p.currentChar() == ' ' || p.currentChar() == '\t' || p.currentChar() == '\n' || p.currentChar() == '\r') {
		p.advance()
	}
}

// skipToEndOfLine 跳到行尾
// @pkg 跳过当前行的剩余部分，用于处理注释
func (p *Parser) skipToEndOfLine() {
	for p.position < len(p.input) && p.currentChar() != '\n' {
		p.advance()
	}
	if p.position < len(p.input) {
		p.advance() // 跳过换行符
	}
}

// errorAt 生成带位置信息的错误
// @pkg 生成包含行号和列号的语法错误信息
// 输入:
//   - message: 错误消息
//
// 输出:
//   - error: 带位置信息的格式化错误
func (p *Parser) errorAt(message string) error {
	return fmt.Errorf("syntax error at line %d, column %d: %s", p.line, p.column, message)
}
