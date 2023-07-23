package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Defines different types of tokens
const (
	TokenConsole  = "CONSOLE"
	TokenLog      = "LOG"
	TokenString   = "STRING"
	TokenInt      = "INT"
	TokenPlus     = "PLUS"
	TokenMinus    = "MINUS"
	TokenMultiply = "MULTIPLY"
	TokenDivide   = "DIVIDE"
	TokenModulo   = "MODULO"
	TokenPower    = "POWER"
)

// Token struct
type Token struct {
	Type    string
	Literal string
}

// Node interface
type Node interface {
	Execute() string
}

// Node type for console.log statements
type ConsoleLogNode struct {
	Arguments []Node
}

// Execute for ConsoleLogNode
func (n *ConsoleLogNode) Execute() string {
	args := make([]string, len(n.Arguments))
	for i, arg := range n.Arguments {
		args[i] = arg.Execute()
	}
	return strings.Join(args, " ")
}

// Node type for string literals
type StringNode struct {
	Value string
}

// Execute for StringNode
func (n *StringNode) Execute() string {
	return n.Value
}

// Node type for addition operation
type PlusNode struct {
	Left  Node
	Right Node
}

// Execute for PlusNode
func (n *PlusNode) Execute() string {
	left, _ := strconv.Atoi(n.Left.Execute())
	right, _ := strconv.Atoi(n.Right.Execute())
	return strconv.Itoa(left + right)
}

// Node type for subtraction operation
type MinusNode struct {
	Left  Node
	Right Node
}

// Execute for MinusNode
func (n *MinusNode) Execute() string {
	left, _ := strconv.Atoi(n.Left.Execute())
	right, _ := strconv.Atoi(n.Right.Execute())
	return strconv.Itoa(left - right)
}

// Node type for multiplication operation
type MultiplyNode struct {
	Left  Node
	Right Node
}

// Execute for MultiplyNode
func (n *MultiplyNode) Execute() string {
	left, _ := strconv.Atoi(n.Left.Execute())
	right, _ := strconv.Atoi(n.Right.Execute())
	return strconv.Itoa(left * right)
}

// Node type for division operation
type DivideNode struct {
	Left  Node
	Right Node
}

// Execute for DivideNode
func (n *DivideNode) Execute() string {
	left, _ := strconv.Atoi(n.Left.Execute())
	right, _ := strconv.Atoi(n.Right.Execute())
	return strconv.Itoa(left / right)
}

// Node type for modulo operation
type ModuloNode struct {
	Left  Node
	Right Node
}

// Execute for ModuloNode
func (n *ModuloNode) Execute() string {
	left, _ := strconv.Atoi(n.Left.Execute())
	right, _ := strconv.Atoi(n.Right.Execute())
	return strconv.Itoa(left % right)
}

// Node type for power operation
type PowerNode struct {
	Left  Node
	Right Node
}

// Execute for PowerNode
func (n *PowerNode) Execute() string {
	left, _ := strconv.Atoi(n.Left.Execute())
	right, _ := strconv.Atoi(n.Right.Execute())
	result := math.Pow(float64(left), float64(right))
	return strconv.Itoa(int(result))
}

// Node type for integer literals
type IntNode struct {
	Value string
}

// Execute for IntNode
func (n *IntNode) Execute() string {
	return n.Value
}

// Lex function to convert the input string into tokens
func Lex(input string) []Token {
	tokens := []Token{}
	statements := strings.Split(input, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		startIndex := strings.Index(stmt, "(")
		endIndex := strings.LastIndex(stmt, ")")

		consoleLog := strings.FieldsFunc(stmt[:startIndex], func(r rune) bool {
			return r == ' ' || r == '.'
		})
		arguments := strings.Split(stmt[startIndex+1:endIndex], ",")

		for _, word := range consoleLog {
			if word == "console" {
				tokens = append(tokens, Token{Type: TokenConsole, Literal: word})
			} else if word == "log" {
				tokens = append(tokens, Token{Type: TokenLog, Literal: word})
			}
		}

		for _, arg := range arguments {
			arg = strings.TrimSpace(arg)
			if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
				tokens = append(tokens, Token{Type: TokenString, Literal: arg[1 : len(arg)-1]})
			} else if strings.ContainsAny(arg, "+-*%/^") {
				operatorIndex := strings.IndexAny(arg, "+-*%/^")
				num1 := strings.TrimSpace(arg[:operatorIndex])
				operator := strings.TrimSpace(arg[operatorIndex : operatorIndex+1])
				num2 := strings.TrimSpace(arg[operatorIndex+1:])
				tokens = append(tokens, Token{Type: TokenInt, Literal: num1})
				switch operator {
				case "+":
					tokens = append(tokens, Token{Type: TokenPlus, Literal: operator})
				case "-":
					tokens = append(tokens, Token{Type: TokenMinus, Literal: operator})
				case "*":
					tokens = append(tokens, Token{Type: TokenMultiply, Literal: operator})
				case "/":
					tokens = append(tokens, Token{Type: TokenDivide, Literal: operator})
				case "%":
					tokens = append(tokens, Token{Type: TokenModulo, Literal: operator})
				case "^":
					tokens = append(tokens, Token{Type: TokenPower, Literal: operator})
				}
				tokens = append(tokens, Token{Type: TokenInt, Literal: num2})
			} else {
				tokens = append(tokens, Token{Type: TokenInt, Literal: arg})
			}
		}
	}

	return tokens
}

// Parse function to convert the tokens into AST nodes
func Parse(tokens []Token) []Node {
	nodes := []Node{}

	i := 0
	for i < len(tokens) {
		if tokens[i].Type == TokenConsole && tokens[i+1].Type == TokenLog {
			i += 2

			args := []Node{}
			for i < len(tokens) && tokens[i].Type != TokenConsole {
				if tokens[i].Type == TokenString {
					args = append(args, &StringNode{Value: tokens[i].Literal})
				} else if tokens[i].Type == TokenInt {
					if i+2 < len(tokens) && tokens[i+2].Type == TokenInt {
						switch tokens[i+1].Type {
						case TokenPlus:
							args = append(args, &PlusNode{Left: &IntNode{Value: tokens[i].Literal}, Right: &IntNode{Value: tokens[i+2].Literal}})
						case TokenMinus:
							args = append(args, &MinusNode{Left: &IntNode{Value: tokens[i].Literal}, Right: &IntNode{Value: tokens[i+2].Literal}})
						case TokenMultiply:
							args = append(args, &MultiplyNode{Left: &IntNode{Value: tokens[i].Literal}, Right: &IntNode{Value: tokens[i+2].Literal}})
						case TokenDivide:
							args = append(args, &DivideNode{Left: &IntNode{Value: tokens[i].Literal}, Right: &IntNode{Value: tokens[i+2].Literal}})
						case TokenModulo:
							args = append(args, &ModuloNode{Left: &IntNode{Value: tokens[i].Literal}, Right: &IntNode{Value: tokens[i+2].Literal}})
						case TokenPower:
							args = append(args, &PowerNode{Left: &IntNode{Value: tokens[i].Literal}, Right: &IntNode{Value: tokens[i+2].Literal}})
						}
						i += 2
					} else {
						args = append(args, &IntNode{Value: tokens[i].Literal})
					}
				}
				i++
			}

			nodes = append(nodes, &ConsoleLogNode{Arguments: args})
		} else {
			panic("Invalid syntax")
		}
	}

	return nodes
}

// Eval function to take a slice of nodes (AST) and evaluate them
func Eval(nodes []Node) {
	for _, node := range nodes {
		fmt.Println(node.Execute())
	}
}

// Main function to read the content of a .es file and pass it to the lexer, parser, and finally to the evaluator
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file to execute")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	tokens := Lex(string(data))
	fmt.Println("Tokens:")
	for _, token := range tokens {
		fmt.Printf("Type: %s, Literal: %s\n", token.Type, token.Literal)
	}

	ast := Parse(tokens)
	fmt.Println("\nAbstract Syntax Tree:")
	for _, node := range ast {
		fmt.Printf("%T: %s\n", node, node.Execute())
	}

	fmt.Println("\nOutput:")
	Eval(ast)
}
