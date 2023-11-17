package main

import (
	"fmt"
	"strconv"
)

// variables map[identifier]value
var variables = make(map[string]string)

func (n *node) runExpression() (expressionResult node) {
	switch n.name {
	case "print":
		fmt.Println(n.params[0].run().value)
		return node{
			kind:  aNumberLiteral,
			value: "1",
		}
	case "*":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		return node{
			kind:  aNumberLiteral,
			value: strconv.Itoa(left * right),
		}
	case "/":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		return node{
			kind:  aNumberLiteral,
			value: strconv.Itoa(left / right),
		}
	case "+":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		return node{
			kind:  aNumberLiteral,
			value: strconv.Itoa(left + right),
		}
	case "-":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		return node{
			kind:  aNumberLiteral,
			value: strconv.Itoa(left - right),
		}
	case ">":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		if left > right {
			return node{
				kind:  aNumberLiteral,
				value: "1",
			}
		} else {
			return node{
				kind:  aNumberLiteral,
				value: "0",
			}
		}
	case ">=":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		if left >= right {
			return node{
				kind:  aNumberLiteral,
				value: "1",
			}
		} else {
			return node{
				kind:  aNumberLiteral,
				value: "0",
			}
		}
	case "<":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		if left < right {
			return node{
				kind:  aNumberLiteral,
				value: "1",
			}
		} else {
			return node{
				kind:  aNumberLiteral,
				value: "0",
			}
		}
	case "<=":
		left, _ := strconv.Atoi(n.params[0].run().value)
		right, _ := strconv.Atoi(n.params[1].run().value)
		if left <= right {
			return node{
				kind:  aNumberLiteral,
				value: "1",
			}
		} else {
			return node{
				kind:  aNumberLiteral,
				value: "0",
			}
		}
	case "==":
		left := n.params[0].run().value
		right := n.params[1].run().value
		if left == right {
			return node{
				kind:  aNumberLiteral,
				value: "1",
			}
		} else {
			return node{
				kind:  aNumberLiteral,
				value: "0",
			}
		}
	case "!=":
		left := n.params[0].run().value
		right := n.params[1].run().value
		if left != right {
			return node{
				kind:  aNumberLiteral,
				value: "1",
			}
		} else {
			return node{
				kind:  aNumberLiteral,
				value: "0",
			}
		}
	default:
		return node{
			kind:  aNumberLiteral,
			value: variables[n.name],
		}
	}
}

func (n *node) runStatement() (expressionResult node) {
	l := len(n.body)
	for i := 0; i < l; i++ {
		n.body[i].run()
	}
	return node{
		kind:  aNumberLiteral,
		value: "1",
	}
}

func (n *node) runAssignmentStatement() (expressionResult node) {
	variables[n.params[0].name] = n.params[1].run().value
	return node{
		kind:  aNumberLiteral,
		value: "1",
	}
}

func (n *node) runStatementIf() (expressionResult node) {
	if n.params[0].run().value != "0" {
		n.body[0].run()
	} else {
		if len(n.body) > 1 {
			n.body[1].run()
		}
	}
	return node{
		kind:  aNumberLiteral,
		value: "1",
	}
}

func (n *node) runStatementWhile() (expressionResult node) {
	for n.params[0].run().value != "0" {
		n.body[0].run()
	}

	return node{
		kind:  aNumberLiteral,
		value: "1",
	}
}

func (n *node) runStatementFor() (expressionResult node) {
	n.params[0].run()
	for n.params[2].run().value != "0" {
		n.body[0].run()
		n.params[4].run()
	}

	return node{
		kind:  aNumberLiteral,
		value: "1",
	}
}

// run that ast
func (n *node) run() node {
	var expressionResult node
	switch n.kind {
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
