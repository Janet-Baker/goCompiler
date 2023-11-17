package main

import (
	"fmt"
	"strconv"
)

// variables map[identifier]value
var variables = make(map[string]string)

func (n *Node) runExpression() (expressionResult Node) {
	//if len(n.Params) == 1 && n.Name == "(" {
	//	return n.Params[0].run()
	//}
	switch n.Name {
	case "print":
		fmt.Println(n.Params[0].run().Value)
		return Node{
			Kind:  aNumberLiteral,
			Value: "1",
		}
	case "*":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		return Node{
			Kind:  aNumberLiteral,
			Value: strconv.Itoa(left * right),
		}
	case "/":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		return Node{
			Kind:  aNumberLiteral,
			Value: strconv.Itoa(left / right),
		}
	case "+":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		return Node{
			Kind:  aNumberLiteral,
			Value: strconv.Itoa(left + right),
		}
	case "-":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		return Node{
			Kind:  aNumberLiteral,
			Value: strconv.Itoa(left - right),
		}
	case ">":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		if left > right {
			return Node{
				Kind:  aNumberLiteral,
				Value: "1",
			}
		} else {
			return Node{
				Kind:  aNumberLiteral,
				Value: "0",
			}
		}
	case ">=":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		if left >= right {
			return Node{
				Kind:  aNumberLiteral,
				Value: "1",
			}
		} else {
			return Node{
				Kind:  aNumberLiteral,
				Value: "0",
			}
		}
	case "<":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		if left < right {
			return Node{
				Kind:  aNumberLiteral,
				Value: "1",
			}
		} else {
			return Node{
				Kind:  aNumberLiteral,
				Value: "0",
			}
		}
	case "<=":
		left, _ := strconv.Atoi(n.Params[0].run().Value)
		right, _ := strconv.Atoi(n.Params[1].run().Value)
		if left <= right {
			return Node{
				Kind:  aNumberLiteral,
				Value: "1",
			}
		} else {
			return Node{
				Kind:  aNumberLiteral,
				Value: "0",
			}
		}
	case "==":
		left := n.Params[0].run().Value
		right := n.Params[1].run().Value
		if left == right {
			return Node{
				Kind:  aNumberLiteral,
				Value: "1",
			}
		} else {
			return Node{
				Kind:  aNumberLiteral,
				Value: "0",
			}
		}
	case "!=":
		left := n.Params[0].run().Value
		right := n.Params[1].run().Value
		if left != right {
			return Node{
				Kind:  aNumberLiteral,
				Value: "1",
			}
		} else {
			return Node{
				Kind:  aNumberLiteral,
				Value: "0",
			}
		}
	case "(":
		return n.Params[0].run()
	default:
		return Node{
			Kind:  aNumberLiteral,
			Value: variables[n.Name],
		}
	}
}

func (n *Node) runStatement() (expressionResult Node) {
	l := len(n.Body)
	for i := 0; i < l; i++ {
		n.Body[i].run()
	}
	return Node{
		Kind:  aNumberLiteral,
		Value: "1",
	}
}

func (n *Node) runAssignmentStatement() (expressionResult Node) {
	variables[n.Params[0].Name] = n.Params[1].run().Value
	return Node{
		Kind:  aNumberLiteral,
		Value: "1",
	}
}

func (n *Node) runStatementIf() (expressionResult Node) {
	if n.Params[0].run().Value != "0" {
		n.Body[0].run()
	} else {
		if len(n.Body) > 1 {
			n.Body[1].run()
		}
	}
	return Node{
		Kind:  aNumberLiteral,
		Value: "1",
	}
}

func (n *Node) runStatementWhile() (expressionResult Node) {
	for n.Params[0].run().Value != "0" {
		n.runStatement()
	}

	return Node{
		Kind:  aNumberLiteral,
		Value: "1",
	}
}

func (n *Node) runStatementFor() (expressionResult Node) {
	if len(n.Params) != 5 {
		fmt.Printf("for statement error, skipping.\ninvalid params at line %d, col %d\n", n.token.line, n.token.col)
		return Node{
			Kind:  aNumberLiteral,
			Value: "1",
		}
	}
	n.Params[0].run()
	for n.Params[2].run().Value != "0" {
		n.runStatement()
		n.Params[4].run()
	}

	return Node{
		Kind:  aNumberLiteral,
		Value: "1",
	}
}

// run that ast
func (n *Node) run() Node {
	var expressionResult Node
	switch n.Kind {
	case aExpression:
		expressionResult = n.runExpression()
		break
	case aProgram:
		expressionResult = n.runStatement()
		break
	case aStatement:
		expressionResult = n.runStatement()
		break
	case aAssignmentStatement:
		expressionResult = n.runAssignmentStatement()
		break
	case aStatementIf:
		expressionResult = n.runStatementIf()
		break
	case aStatementWhile:
		expressionResult = n.runStatementWhile()
		break
	case aStatementFor:
		expressionResult = n.runStatementFor()
		break
	case aNumberLiteral:
		expressionResult = *n
		break
	case aStringLiteral:
		expressionResult = *n
		break
	}
	return expressionResult
}
