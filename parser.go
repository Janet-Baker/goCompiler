/**
 * ============================================================================
 *                                 ヽ/❀o ل͜ o\ﾉ
 *                                THE PARSER!!!
 * ============================================================================
 */

package main

import (
	"errors"
	"fmt"
)

/*
We will define our type, `node` here. Within node are pointers types to what
would otherwise be recursive types in Go. e.g.

callee    node

Would cause the Go compiler to complain about a recursive type. When we want
to use one of these types to pass through to a function, for example, we'd
use `&` as it'd be a reference. But we'll come to that a bit later on.
*/
type node struct {
	kind   int
	value  string
	name   string
	token  token // 打印错误信息用
	body   []node
	params []node
	//callee     *node
	//expression *node
	//arguments  *[]node
	//context    *[]node
}

// kind of ast
// aProgram -> aStatement | aExpression
// aExpression -> + | - | * | / | Function
// aStatement -> (If | Else | For) [aExpression] {aStatement}
// aNumberLiteral -> tInteger
const (
	aBlank = 100 + iota
	//aProgram 一个子程序
	aProgram

	//aExpression 语法表达式，即语法。
	//是指一个计算值的代码片段，它可以由变量引用、数值计算、函数调用等组成。
	//表达式通常会产生一个值，并可以用于组成更复杂的表达式或用于赋值给某个变量。
	// 在params中解析
	aExpression

	//aStatement 语法树节点，即语句。
	//是指一条执行操作或者完成某个动作的代码指令。
	//语句不一定产生一个值，它可以是赋值、条件判断、循环等。
	//语句用于组织代码的执行顺序，使程序按照预期的逻辑执行。
	// 通常在body中解析
	aStatement

	//有参数的语句块，例如if(){}else{}
	/*var ifNode = node{
		kind: aStatementWithBodyAndParams,
		name: "if",
		params: []node{
			{
				kind: aExpression,
			},
		},
		body: make([]node, 2),
	}*/
	aStatementIf
	aStatementFor
	aStatementWhile

	//aAssignmentStatement 赋值语句
	aAssignmentStatement

	//aLiteral 一个字面量（常量）
	//aNumberLiteral 一个数字字面量
	aNumberLiteral
	//aStringLiteral 一个字符串字面量
	aStringLiteral
)

/*This is the counter variable that we'll use for parsing.*/
var pc int

/*This variable will store our slice of `token`s inside of it.*/
var pt []token

/*Okay, so we define a `parser` function that accepts our slice of `tokens`.*/
/*
var astNode = ast{
	kind: aProgram,
	body: []node{
		{
			kind: aExpression,
			name: "+",
			params: []node{
				{
					kind:  aNumberLiteral,
					value: "2",
				},{
					kind:  aNumberLiteral,
					value: "4",
				},
			},
		},
	},
}*/
func parser(tokens *[]token) (node, error) {
	/*Here, we initially give both the parser counter and the parser tokens a
	value.*/
	pc = 0
	pt = *tokens

	/*Now, we're going to create our AST which will have a root which is a
	`Program` node.*/
	astRoot := node{
		kind: aProgram,
		body: []node{},
	}
	ns.push(&astRoot)
	/*And we're going to kickstart our `walk` function, which you can find just
	below this, we'll be pushing nodes to our `ast.body` slice.

	The reason we are doing this inside a loop is because our program can have
	`CallExpressions` after one another instead of being nested.

	  a = 100 + 200
	  print(a)
	*/
	for pc < len(*tokens) {
		astBodyNode, err := walk()
		if err == nil {
			astRoot.body = append(astRoot.body, astBodyNode)
		} else if err.Error() == "skip" {
			continue
		} else {
			return node{}, err
		}
	}

	/*At the end of our parser we'll return the AST.*/
	return astRoot, nil
}

/*
But this time we're going to use recursion instead of a `while` loop. So we
define a `walk` function.
*/
// walk through the tokens and generate ast

// Priority:
// 1. ()
// 2. * /
// 3. + -
// 4. > >= < <= == !=
// 5. =

// a = 1 + 2
// a = 1 + 2 * 3
// print(a)
// int f(int a, int b) { return a + b; }
// f(1, 2)
// int main(){}
func walk() (node, error) {
	/*Inside the walk function we start by grabbing the `current` token.*/
	currentToken := pt[pc]

	// LParen as the start of a aExpression
	/*We start this off when we	encounter an open parenthesis.*/
	// ()
	if currentToken.kind == tLParen {

		/*We'll increment `current` to skip the parenthesis since we don't care
		about it in our AST.*/
		pc++
		currentToken = pt[pc]

		/*We create a base node with the type `aExpression`, and we're going
		to set the name as the current token's value since the next token after
		the open parenthesis is the name of the function.*/

		currentNode := node{
			kind:   aExpression,
			name:   currentToken.value,
			token:  currentToken,
			params: []node{},
		}
		ns.push(&currentNode)

		// So we create a `for` loop that will continue until it encounters a
		// token with a `type` of `'paren'` and a `value` of a closing
		// parenthesis.

		//for currentToken.kind != "paren" || (currentToken.kind == "paren" && currentToken.value != ")") {
		for currentToken.kind != tRParen {
			// we'll call the `walk` function which will return a `node` and we'll
			// push it into our `node.params`.
			tempNode, err := walk()
			if err == nil {
				currentNode.params = append(currentNode.params, tempNode)
				currentToken = pt[pc]
			} else if err.Error() == "skip" {
				currentToken = pt[pc]
			} else {
				return node{}, err
			}
		}

		// Finally we will increment `current` one last time to skip the closing
		// parenthesis.
		pc++
		currentNode = *ns.pop()
		// And return the node.
		return currentNode, nil
	}

	// {}
	if currentToken.kind == tLBrace {
		/*We'll increment `current` to skip the parenthesis since we don't care
		about it in our AST.*/
		pc++
		currentToken = pt[pc]

		/*We create a base node with the type `aExpression`, and we're going
		to set the name as the current token's value since the next token after
		the open parenthesis is the name of the function.*/
		currentNode := node{
			kind:  aStatement,
			name:  currentToken.value,
			token: currentToken,
			body:  []node{},
		}
		ns.push(&currentNode)
		for currentToken.kind != tRBrace {
			// we'll call the `walk` function which will return a `node` and we'll
			// push it into our `node.params`.
			tempNode, err := walk()
			if err == nil {
				currentNode.body = append(currentNode.body, tempNode)
				currentToken = pt[pc]
			} else if err.Error() == "skip" {
				currentToken = pt[pc]
			} else {
				return node{}, err
			}
		}

		// Finally we will increment `current` one last time to skip the closing
		// parenthesis.
		pc++

		currentNode = *ns.pop()
		// And return the node.
		return currentNode, nil
	}

	// * /
	if currentToken.kind == tMultiply || currentToken.kind == tDivide {
		lastSubNode := ns.peek()
		if lastSubNode.kind == aProgram || lastSubNode.kind == aStatement {
			// working at the top level
			l := len(lastSubNode.body)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.kind == aBlank {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					// insert to the ground of the tree
					parentOfLastNode := &lastSubNode.body[l-1]
					if len(parentOfLastNode.params) > 0 {
						lastNode := &parentOfLastNode.params[len(parentOfLastNode.params)-1]
						for lastNode.name == "+" || lastNode.name == "-" || lastNode.name == ">" || lastNode.name == ">=" ||
							lastNode.name == "<" || lastNode.name == "<=" || lastNode.name == "==" || lastNode.name == "!=" ||
							lastNode.kind == aAssignmentStatement {
							parentOfLastNode = lastNode
							lastNode = &lastNode.params[len(lastNode.params)-1]
						}
						newNode := node{
							kind:  aExpression,
							name:  currentToken.value,
							token: currentToken,
							params: []node{
								*lastNode,
								rightNode,
							},
						}
						parentOfLastNode.params[len(parentOfLastNode.params)-1] = newNode
						return newNode, errors.New("skip")
					} else {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
					}

				}
			} else {
				return node{}, fmt.Errorf("unexpected end of tokens")
			}
		}
		// modify lastSubNode in () or {}
		if lastSubNode.kind == aExpression {
			newNode := node{
				kind:  aExpression,
				token: currentToken,
				name:  currentToken.value,
			}
			if pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.kind == aBlank {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					parentOfLastNode := lastSubNode
					if len(parentOfLastNode.params) > 0 {
						lastNode := &lastSubNode.params[len(lastSubNode.params)-1]
						for lastNode.name == "+" || lastNode.name == "-" || lastNode.name == ">" || lastNode.name == ">=" ||
							lastNode.name == "<" || lastNode.name == "<=" || lastNode.name == "==" || lastNode.name == "!=" ||
							lastNode.kind == aAssignmentStatement {
							parentOfLastNode = lastNode
							lastNode = &lastNode.params[len(lastNode.params)-1]
						}
						newNode.params = []node{
							*lastNode,
							rightNode,
						}
						parentOfLastNode.params[len(parentOfLastNode.params)-1] = newNode
						return newNode, errors.New("skip")
					} else {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
					}

				}
			} else {
				return node{}, fmt.Errorf("unexpected end of tokens")
			}
			return newNode, errors.New("skip")
		}
	}

	// + -
	if currentToken.kind == tPlus || currentToken.kind == tMinus {
		lastSubNode := ns.peek()
		if lastSubNode.kind == aProgram || lastSubNode.kind == aStatement {
			// working at the top level
			l := len(lastSubNode.body)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.kind == aBlank {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					// insert to the ground of the tree
					parentOfLastNode := &lastSubNode.body[l-1]
					lastNode := &parentOfLastNode.params[len(parentOfLastNode.params)-1]
					for lastNode.name == ">" || lastNode.name == ">=" || lastNode.name == "<" || lastNode.name == "<=" ||
						lastNode.name == "==" || lastNode.name == "!=" || lastNode.kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.params[len(lastNode.params)-1]
					}
					newNode := node{
						kind:  aExpression,
						name:  currentToken.value,
						token: currentToken,
						params: []node{
							*lastNode,
							rightNode,
						},
					}
					parentOfLastNode.params[len(parentOfLastNode.params)-1] = newNode
					return newNode, errors.New("skip")

				}
			} else {
				return node{}, fmt.Errorf("unexpected end of tokens")
			}
		}
		// modify lastSubNode in () or {}
		if lastSubNode.kind == aExpression {
			newNode := node{
				kind:  aExpression,
				token: currentToken,
				name:  currentToken.value,
			}
			if pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.kind == aBlank {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					parentOfLastNode := lastSubNode
					lastNode := &lastSubNode.params[len(lastSubNode.params)-1]
					for lastNode.name == ">" || lastNode.name == ">=" || lastNode.name == "<" || lastNode.name == "<=" ||
						lastNode.name == "==" || lastNode.name == "!=" || lastNode.kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.params[len(lastNode.params)-1]
					}
					newNode.params = []node{
						*lastNode,
						rightNode,
					}
					parentOfLastNode.params[len(parentOfLastNode.params)-1] = newNode
					return newNode, errors.New("skip")
				}
			} else {
				return node{}, fmt.Errorf("unexpected end of tokens")
			}
			return newNode, errors.New("skip")
		}
	}

	// > >= < <= == !=
	if currentToken.kind >= tCalcNotEqual && currentToken.kind <= tCalcEqual {
		lastSubNode := ns.peek()
		if lastSubNode.kind == aProgram || lastSubNode.kind == aStatement {
			// working at the top level
			l := len(lastSubNode.body)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.kind == aBlank {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					// insert to the ground of the tree
					parentOfLastNode := &lastSubNode.body[l-1]
					lastNode := &parentOfLastNode.params[len(parentOfLastNode.params)-1]
					for lastNode.kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.params[len(lastNode.params)-1]
					}

					newNode := node{
						kind:  aExpression,
						name:  currentToken.value,
						token: currentToken,
						params: []node{
							*lastNode,
							rightNode,
						},
					}
					parentOfLastNode.params[len(parentOfLastNode.params)-1] = newNode
					return newNode, errors.New("skip")

				}
			} else {
				return node{}, fmt.Errorf("unexpected end of tokens")
			}
		}
		// modify lastSubNode in () or {}
		if lastSubNode.kind == aExpression {
			newNode := node{
				kind:  aExpression,
				token: currentToken,
				name:  currentToken.value,
			}
			if pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.kind == aBlank {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					parentOfLastNode := lastSubNode
					lastNode := &lastSubNode.params[len(lastSubNode.params)-1]
					for lastNode.kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.params[len(lastNode.params)-1]
					}

					newNode.params = []node{
						*lastNode,
						rightNode,
					}
					parentOfLastNode.params[len(parentOfLastNode.params)-1] = newNode
					return newNode, errors.New("skip")
				}
			} else {
				return node{}, fmt.Errorf("unexpected end of tokens")
			}
			return newNode, errors.New("skip")
		}
	}

	if currentToken.kind == tEqual {
		lastSubNode := ns.peek()
		if lastSubNode.kind == aProgram || lastSubNode.kind == aStatement {
			// working at the top level
			l := len(lastSubNode.body)
			if l > 0 {
				if pc < len(pt)-1 {
					pc++
					rightNode, err := walk()
					if err == nil {
						if rightNode.kind == aBlank {
							return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
						}
						// a = 1 + 2 * 3

						if len(lastSubNode.body) > 0 {
							lastNode := &lastSubNode.body[l-1]
							if lastNode.kind == aAssignmentStatement {
								return node{}, fmt.Errorf("assigning to an assigning statement at line%d, column%d", currentToken.line, currentToken.col)
							}
							if lastNode.token.kind != tIdentifier {
								return node{}, fmt.Errorf("trying to assign to a non-id target at line%d, column%d", currentToken.line, currentToken.col)
							}

							newNode := node{
								kind:  aAssignmentStatement,
								name:  currentToken.value,
								token: currentToken,
								params: []node{
									*lastNode,
									rightNode,
								},
							}
							lastSubNode.body[l-1] = newNode
							return newNode, errors.New("skip")
						} else {
							return node{}, fmt.Errorf("can not find assigning target at line%d, column%d", currentToken.line, currentToken.col)
						}
					}
				} else {
					return node{}, fmt.Errorf("unexpected end of tokens")
				}
			} else {
				return node{}, fmt.Errorf("can not find assigning target at line%d, column%d", currentToken.line, currentToken.col)
			}
		} else if lastSubNode.kind == aExpression {
			// for(i=0;i<n;i++)
			l := len(lastSubNode.params)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.kind == aBlank {
						return node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					if len(lastSubNode.params) > 0 {
						lastNode := &lastSubNode.params[l-1]
						if lastNode.kind == aAssignmentStatement {
							return node{}, fmt.Errorf("assigning to an assigning statement at line%d, column%d", currentToken.line, currentToken.col)
						}
						if lastNode.token.kind != tIdentifier {
							return node{}, fmt.Errorf("trying to assign to a non-id target at line%d, column%d", currentToken.line, currentToken.col)
						}
						newNode := node{
							kind:  aAssignmentStatement,
							name:  currentToken.value,
							token: currentToken,
							params: []node{
								*lastNode,
								rightNode,
							},
						}
						lastSubNode.params[l-1] = newNode
						return newNode, errors.New("skip")
					} else {
						return node{}, fmt.Errorf("can not find assigning target at line%d, column%d", currentToken.line, currentToken.col)
					}
				}
			} else {
				return node{}, fmt.Errorf("unexpected end of tokens")
			}
		} else {
			return node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
		}
	}

	if currentToken.kind == tInteger {
		/*If we have one, we'll increment `current`.*/
		pc++
		/*And we'll return a new AST node called `NumberLiteral` and setting its
		value to the value of our token.*/
		newNode := node{
			kind:  aNumberLiteral,
			name:  currentToken.value,
			value: currentToken.value,
			token: currentToken,
		}
		return newNode, nil
	}

	if currentToken.kind == tString {
		pc++
		newNode := node{
			kind:  aStringLiteral,
			name:  currentToken.value,
			value: currentToken.value,
			token: currentToken,
		}
		return newNode, nil
	}

	// tPrint calls fmt.Println() directly
	if currentToken.kind == tPrint {
		currentNode := node{
			kind:  aStatement,
			name:  currentToken.value,
			token: currentToken,
		}
		pc++
		p1, err := walk()
		if err != nil {
			return node{}, err
		}
		currentNode.params = []node{p1}
		return p1, nil
	}

	// tIf -> aStatement{body={aStatementIf{param=(aExpression); body={aStatement}}}}
	if currentToken.kind == tIf {
		currentNode := node{
			kind:  aStatementIf,
			name:  currentToken.value,
			token: currentToken,
		}

		// if param
		ifExpression := node{
			kind: aExpression,
		}
		ns.push(&ifExpression)

		pc++
		p1, err := walk()
		if err != nil {
			return node{}, err
		}
		ifExpression = *ns.pop()
		currentNode.params = []node{p1}

		// if body
		if pt[pc].kind == tLBrace {
			pc++
			ifTrue, err := walk()
			if err != nil {
				return node{}, err
			}
			currentNode.body = []node{ifTrue}
		} else {
			return node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
		}

		// else? append to body
		if pt[pc].kind == tElse {
			pc++
			ifFalseElse, err := walk()
			if err != nil {
				return node{}, err
			}
			currentNode.body = append(currentNode.body, ifFalseElse)
		}

		return currentNode, nil
	}

	if currentToken.kind == tFor {
		currentNode := node{
			kind:  aStatementFor,
			name:  currentToken.value,
			token: currentToken,
		}
		// for param
		forExpression := node{
			kind: aExpression,
		}
		ns.push(&forExpression)

		pc++
		p1, err := walk()
		if err != nil {
			return node{}, err
		}
		forExpression = *ns.pop()
		currentNode.params = []node{p1}

		// for body
		if pt[pc].kind == tLBrace {
			pc++
			forBody, err := walk()
			if err != nil {
				return node{}, err
			}
			currentNode.body = []node{forBody}
		} else {
			return node{}, fmt.Errorf("unexpected token at line%d, column%d, should be for body", currentToken.line, currentToken.col)
		}

		return currentNode, nil
	}

	if currentToken.kind == tWhile {
		currentNode := node{
			kind:  aStatementWhile,
			name:  currentToken.value,
			token: currentToken,
		}
		// while param
		whileExpression := node{
			kind: aExpression,
		}
		ns.push(&whileExpression)

		pc++
		p1, err := walk()
		if err != nil {
			return node{}, err
		}
		whileExpression = *ns.pop()
		currentNode.params = p1.params

		// while body
		if pt[pc].kind == tLBrace {
			pc++
			whileBody, err := walk()
			if err != nil {
				return node{}, err
			}
			currentNode.body = []node{whileBody}
		} else {
			return node{}, fmt.Errorf("unexpected token at line%d, column%d, should be while body", currentToken.line, currentToken.col)
		}

		return currentNode, nil
	}

	// tIdentifier Function Call
	if currentToken.kind == tIdentifier {
		if pc < len(pt)-1 {
			// look if it is Assignment Statement like `a = `
			/*			if pt[pc+1].kind == tEqual {
						currentNode := node{
							kind:  aAssignmentStatement,
							name:  currentToken.value,
							token: pt[pc+1],
						}
						pc = pc + 2
						p1, err := walk()
						if err == nil {
							currentNode.params = []node{p1}
						}
						return currentNode, nil
					} else if pt[pc+1].kind == tLParen {*/
			if pt[pc+1].kind == tLParen {
				// looks like it is Calling a function like `a()`
				currentNode := node{
					kind:  aExpression,
					name:  currentToken.value,
					token: currentToken,
				}
				pc++
				p1, err := walk()
				if err == nil {
					currentNode.params = []node{p1}
				}
				return currentNode, nil
			} else {
				// looks like someone is calling us like `1 + a`
				pc++
				currentNode := node{
					kind:  aExpression,
					name:  currentToken.value,
					token: currentToken,
				}
				return currentNode, nil
			}
		}
	}

	// we skip it when we don't know what is it
	pc++
	return node{
		kind:  aBlank,
		name:  currentToken.value,
		token: currentToken,
	}, nil
}

// nodeStack is a simple LIFO
// that stores for LRs objects like `()` and `{}`
type nodeStack struct {
	nodes []*node
}

var ns = nodeStack{make([]*node, 0)}

func (s *nodeStack) push(n *node) {
	s.nodes = append(s.nodes, n)
}

func (s *nodeStack) pop() *node {
	if len(s.nodes) == 0 {
		return nil
	}
	n := s.nodes[len(s.nodes)-1]
	s.nodes = s.nodes[:len(s.nodes)-1]
	return n
}

func (s *nodeStack) peek() *node {
	if len(s.nodes) == 0 {
		return nil
	}
	return s.nodes[len(s.nodes)-1]
}

func (s *nodeStack) isEmpty() bool {
	return len(s.nodes) == 0
}
