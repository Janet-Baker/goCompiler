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
We will define our type, `node` here. Within Node are pointers types to what
would otherwise be recursive types in Go. e.g.

callee    Node

Would cause the Go compiler to complain about a recursive type. When we want
to use one of these types to pass through to a function, for example, we'd
use `&` as it'd be a reference. But we'll come to that a bit later on.
*/

type Node struct {
	Kind   int
	Value  string
	Name   string
	token  token // 打印错误信息用
	Body   []Node
	Params []Node
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
	//aBlank 空白标识符
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
	/*var ifNode = Node{
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
func parser(tokens *[]token) (Node, error) {
	/*Here, we initially give both the parser counter and the parser tokens a
	value.*/
	pc = 0
	pt = *tokens

	/*Now, we're going to create our AST which will have a root which is a
	`Program` node.*/
	astRoot := Node{
		Kind: aProgram,
		Body: []Node{},
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
			astRoot.Body = append(astRoot.Body, astBodyNode)
		} else if err.Error() == "skip" {
			continue
		} else {
			return Node{}, err
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
func walk() (Node, error) {
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

		currentNode := Node{
			Kind:   aExpression,
			Name:   currentToken.value,
			token:  currentToken,
			Params: []Node{},
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
				currentNode.Params = append(currentNode.Params, tempNode)
				currentToken = pt[pc]
			} else if err.Error() == "skip" {
				currentToken = pt[pc]
			} else {
				return Node{}, err
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
		currentNode := Node{
			Kind:  aStatement,
			Name:  currentToken.value,
			token: currentToken,
			Body:  []Node{},
		}
		ns.push(&currentNode)
		for currentToken.kind != tRBrace {
			// we'll call the `walk` function which will return a `node` and we'll
			// push it into our `node.params`.
			tempNode, err := walk()
			if err == nil {
				currentNode.Body = append(currentNode.Body, tempNode)
				currentToken = pt[pc]
			} else if err.Error() == "skip" {
				currentToken = pt[pc]
			} else {
				return Node{}, err
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
		if lastSubNode.Kind == aProgram || lastSubNode.Kind == aStatement {
			// working at the top level
			l := len(lastSubNode.Body)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.Kind == aBlank {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					// insert to the ground of the tree
					parentOfLastNode := &lastSubNode.Body[l-1]
					if len(parentOfLastNode.Params) > 0 {
						lastNode := &parentOfLastNode.Params[len(parentOfLastNode.Params)-1]
						for lastNode.Name == "+" || lastNode.Name == "-" || lastNode.Name == ">" || lastNode.Name == ">=" ||
							lastNode.Name == "<" || lastNode.Name == "<=" || lastNode.Name == "==" || lastNode.Name == "!=" ||
							lastNode.Kind == aAssignmentStatement {
							parentOfLastNode = lastNode
							lastNode = &lastNode.Params[len(lastNode.Params)-1]
						}
						newNode := Node{
							Kind:  aExpression,
							Name:  currentToken.value,
							token: currentToken,
							Params: []Node{
								*lastNode,
								rightNode,
							},
						}
						parentOfLastNode.Params[len(parentOfLastNode.Params)-1] = newNode
						return newNode, errors.New("skip")
					} else {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
					}

				}
			} else {
				return Node{}, fmt.Errorf("unexpected end of tokens")
			}
		}
		// modify lastSubNode in () or {}
		if lastSubNode.Kind == aExpression {
			newNode := Node{
				Kind:  aExpression,
				token: currentToken,
				Name:  currentToken.value,
			}
			if pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.Kind == aBlank {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					parentOfLastNode := lastSubNode
					if len(parentOfLastNode.Params) > 0 {
						lastNode := &lastSubNode.Params[len(lastSubNode.Params)-1]
						for lastNode.Name == "+" || lastNode.Name == "-" || lastNode.Name == ">" || lastNode.Name == ">=" ||
							lastNode.Name == "<" || lastNode.Name == "<=" || lastNode.Name == "==" || lastNode.Name == "!=" ||
							lastNode.Kind == aAssignmentStatement {
							parentOfLastNode = lastNode
							lastNode = &lastNode.Params[len(lastNode.Params)-1]
						}
						newNode.Params = []Node{
							*lastNode,
							rightNode,
						}
						parentOfLastNode.Params[len(parentOfLastNode.Params)-1] = newNode
						if len(lastSubNode.Params) == 1 {
							_ = ns.pop()
							ns.push(&newNode)
						}
						return newNode, errors.New("skip")
					} else {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
					}

				}
			} else {
				return Node{}, fmt.Errorf("unexpected end of tokens")
			}
			return newNode, errors.New("skip")
		}
	}

	// + -
	if currentToken.kind == tPlus || currentToken.kind == tMinus {
		lastSubNode := ns.peek()
		if lastSubNode.Kind == aProgram || lastSubNode.Kind == aStatement {
			// working at the top level
			l := len(lastSubNode.Body)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.Kind == aBlank {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					// insert to the ground of the tree
					parentOfLastNode := &lastSubNode.Body[l-1]
					lastNode := &parentOfLastNode.Params[len(parentOfLastNode.Params)-1]
					for lastNode.Name == ">" || lastNode.Name == ">=" || lastNode.Name == "<" || lastNode.Name == "<=" ||
						lastNode.Name == "==" || lastNode.Name == "!=" || lastNode.Kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.Params[len(lastNode.Params)-1]
					}
					newNode := Node{
						Kind:  aExpression,
						Name:  currentToken.value,
						token: currentToken,
						Params: []Node{
							*lastNode,
							rightNode,
						},
					}
					parentOfLastNode.Params[len(parentOfLastNode.Params)-1] = newNode
					return newNode, errors.New("skip")

				}
			} else {
				return Node{}, fmt.Errorf("unexpected end of tokens")
			}
		}
		// modify lastSubNode in () or {}
		if lastSubNode.Kind == aExpression {
			newNode := Node{
				Kind:  aExpression,
				token: currentToken,
				Name:  currentToken.value,
			}
			if pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.Kind == aBlank {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					parentOfLastNode := lastSubNode
					lastNode := &lastSubNode.Params[len(lastSubNode.Params)-1]
					for lastNode.Name == ">" || lastNode.Name == ">=" || lastNode.Name == "<" || lastNode.Name == "<=" ||
						lastNode.Name == "==" || lastNode.Name == "!=" || lastNode.Kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.Params[len(lastNode.Params)-1]
					}
					newNode.Params = []Node{
						*lastNode,
						rightNode,
					}
					parentOfLastNode.Params[len(parentOfLastNode.Params)-1] = newNode
					if len(lastSubNode.Params) == 1 {
						_ = ns.pop()
						ns.push(&newNode)
					}
					return newNode, errors.New("skip")
				}
			} else {
				return Node{}, fmt.Errorf("unexpected end of tokens")
			}
			return newNode, errors.New("skip")
		}
	}

	// > >= < <= == !=
	if currentToken.kind >= tCalcNotEqual && currentToken.kind <= tCalcEqual {
		lastSubNode := ns.peek()
		if lastSubNode.Kind == aProgram || lastSubNode.Kind == aStatement {
			// working at the top level
			l := len(lastSubNode.Body)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.Kind == aBlank {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					// insert to the ground of the tree
					parentOfLastNode := &lastSubNode.Body[l-1]
					lastNode := &parentOfLastNode.Params[len(parentOfLastNode.Params)-1]
					for lastNode.Kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.Params[len(lastNode.Params)-1]
					}

					newNode := Node{
						Kind:  aExpression,
						Name:  currentToken.value,
						token: currentToken,
						Params: []Node{
							*lastNode,
							rightNode,
						},
					}
					parentOfLastNode.Params[len(parentOfLastNode.Params)-1] = newNode
					return newNode, errors.New("skip")

				}
			} else {
				return Node{}, fmt.Errorf("unexpected end of tokens")
			}
		}
		// modify lastSubNode in () or {}
		if lastSubNode.Kind == aExpression {
			newNode := Node{
				Kind:  aExpression,
				token: currentToken,
				Name:  currentToken.value,
			}
			if pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.Kind == aBlank {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					parentOfLastNode := lastSubNode
					lastNode := &lastSubNode.Params[len(lastSubNode.Params)-1]
					for lastNode.Kind == aAssignmentStatement {
						parentOfLastNode = lastNode
						lastNode = &lastNode.Params[len(lastNode.Params)-1]
					}

					newNode.Params = []Node{
						*lastNode,
						rightNode,
					}
					parentOfLastNode.Params[len(parentOfLastNode.Params)-1] = newNode
					return newNode, errors.New("skip")
				}
			} else {
				return Node{}, fmt.Errorf("unexpected end of tokens")
			}
			return newNode, errors.New("skip")
		}
	}

	if currentToken.kind == tEqual {
		lastSubNode := ns.peek()
		if lastSubNode.Kind == aProgram || lastSubNode.Kind == aStatement {
			// working at the top level
			l := len(lastSubNode.Body)
			if l > 0 {
				if pc < len(pt)-1 {
					pc++
					rightNode, err := walk()
					if err == nil {
						if rightNode.Kind == aBlank {
							return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
						}
						// a = 1 + 2 * 3

						if len(lastSubNode.Body) > 0 {
							lastNode := &lastSubNode.Body[l-1]
							if lastNode.Kind == aAssignmentStatement {
								return Node{}, fmt.Errorf("assigning to an assigning statement at line%d, column%d", currentToken.line, currentToken.col)
							}
							if lastNode.token.kind != tIdentifier {
								return Node{}, fmt.Errorf("trying to assign to a non-id target at line%d, column%d", currentToken.line, currentToken.col)
							}

							newNode := Node{
								Kind:  aAssignmentStatement,
								Name:  currentToken.value,
								token: currentToken,
								Params: []Node{
									*lastNode,
									rightNode,
								},
							}
							lastSubNode.Body[l-1] = newNode
							return newNode, errors.New("skip")
						} else {
							return Node{}, fmt.Errorf("can not find assigning target at line%d, column%d", currentToken.line, currentToken.col)
						}
					}
				} else {
					return Node{}, fmt.Errorf("unexpected end of tokens")
				}
			} else {
				return Node{}, fmt.Errorf("can not find assigning target at line%d, column%d", currentToken.line, currentToken.col)
			}
		} else if lastSubNode.Kind == aExpression {
			// for(i=0;i<n;i++)
			l := len(lastSubNode.Params)
			if l > 0 && pc < len(pt)-1 {
				pc++
				rightNode, err := walk()
				if err == nil {
					if rightNode.Kind == aBlank {
						return Node{}, fmt.Errorf("unexpected token at line%d, column%d", rightNode.token.line, rightNode.token.col)
					}
					// a = 1 + 2 * 3

					if len(lastSubNode.Params) > 0 {
						lastNode := &lastSubNode.Params[l-1]
						if lastNode.Kind == aAssignmentStatement {
							return Node{}, fmt.Errorf("assigning to an assigning statement at line%d, column%d", currentToken.line, currentToken.col)
						}
						if lastNode.token.kind != tIdentifier {
							return Node{}, fmt.Errorf("trying to assign to a non-id target at line%d, column%d", currentToken.line, currentToken.col)
						}
						newNode := Node{
							Kind:  aAssignmentStatement,
							Name:  currentToken.value,
							token: currentToken,
							Params: []Node{
								*lastNode,
								rightNode,
							},
						}
						lastSubNode.Params[l-1] = newNode
						return newNode, errors.New("skip")
					} else {
						return Node{}, fmt.Errorf("can not find assigning target at line%d, column%d", currentToken.line, currentToken.col)
					}
				}
			} else {
				return Node{}, fmt.Errorf("unexpected end of tokens")
			}
		} else {
			return Node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
		}
	}

	if currentToken.kind == tInteger {
		/*If we have one, we'll increment `current`.*/
		pc++
		/*And we'll return a new AST node called `NumberLiteral` and setting its
		value to the value of our token.*/
		newNode := Node{
			Kind:  aNumberLiteral,
			Name:  currentToken.value,
			Value: currentToken.value,
			token: currentToken,
		}
		return newNode, nil
	}

	if currentToken.kind == tString {
		pc++
		newNode := Node{
			Kind:  aStringLiteral,
			Name:  currentToken.value,
			Value: currentToken.value,
			token: currentToken,
		}
		return newNode, nil
	}

	// tPrint calls fmt.Println() directly
	if currentToken.kind == tPrint {
		currentNode := Node{
			Kind:  aExpression,
			Name:  currentToken.value,
			token: currentToken,
		}
		pc++
		p1, err := walk()
		if err != nil {
			return Node{}, err
		}
		currentNode.Params = p1.Params
		return p1, nil
	}

	// tIf -> aStatement{body={aStatementIf{param=(aExpression); body={aStatement}}}}
	if currentToken.kind == tIf {
		currentNode := Node{
			Kind:  aStatementIf,
			Name:  currentToken.value,
			token: currentToken,
		}

		// if param
		ifExpression := Node{
			Kind: aExpression,
		}
		ns.push(&ifExpression)

		pc++
		p1, err := walk()
		if err != nil {
			return Node{}, err
		}
		_ = ns.pop()
		currentNode.Params = p1.Params

		// if body
		if pt[pc].kind == tLBrace {
			pc++
			ifTrue, err := walk()
			if err != nil {
				return Node{}, err
			}
			currentNode.Body = []Node{ifTrue}
		} else {
			return Node{}, fmt.Errorf("unexpected token at line%d, column%d", currentToken.line, currentToken.col)
		}

		// else? append to body
		if pt[pc+1].kind == tElse {
			pc = pc + 2
			ifFalseElse, err := walk()
			if err != nil {
				return Node{}, err
			}
			currentNode.Body = append(currentNode.Body, ifFalseElse)
		}

		return currentNode, nil
	}

	if currentToken.kind == tFor {
		currentNode := Node{
			Kind:  aStatementFor,
			Name:  currentToken.value,
			token: currentToken,
		}
		// for param
		forExpression := Node{
			Kind: aExpression,
		}
		ns.push(&forExpression)

		pc++
		p1, err := walk()
		if err != nil {
			return Node{}, err
		}
		_ = ns.pop()
		currentNode.Params = p1.Params

		// for body
		if pt[pc].kind == tLBrace {
			pc++
			forBody, err := walk()
			if err != nil {
				return Node{}, err
			}
			currentNode.Body = []Node{forBody}
		} else {
			return Node{}, fmt.Errorf("unexpected token at line%d, column%d, should be for body", currentToken.line, currentToken.col)
		}

		return currentNode, nil
	}

	if currentToken.kind == tWhile {
		currentNode := Node{
			Kind:  aStatementWhile,
			Name:  currentToken.value,
			token: currentToken,
		}
		// while param
		whileExpression := Node{
			Kind: aExpression,
		}
		ns.push(&whileExpression)

		pc++
		p1, err := walk()
		if err != nil {
			return Node{}, err
		}
		_ = ns.pop()
		currentNode.Params = p1.Params

		// while body
		if pt[pc].kind == tLBrace {
			pc++
			whileBody, err := walk()
			if err != nil {
				return Node{}, err
			}
			currentNode.Body = []Node{whileBody}
		} else {
			return Node{}, fmt.Errorf("unexpected token at line%d, column%d, should be while body", currentToken.line, currentToken.col)
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
				currentNode := Node{
					Kind:  aExpression,
					Name:  currentToken.value,
					token: currentToken,
				}
				pc++
				p1, err := walk()
				if err == nil {
					currentNode.Params = []Node{p1}
				}
				return currentNode, nil
			} else {
				// looks like someone is calling us like `1 + a`
				pc++
				currentNode := Node{
					Kind:  aExpression,
					Name:  currentToken.value,
					token: currentToken,
				}
				return currentNode, nil
			}
		}
	}

	// we skip it when we don't know what is it
	pc++
	return Node{
		Kind:  aBlank,
		Name:  currentToken.value,
		token: currentToken,
	}, nil
}

// nodeStack is a simple LIFO
// that stores for LRs objects like `()` and `{}`
type nodeStack struct {
	nodes []*Node
}

var ns = nodeStack{make([]*Node, 0)}

func (s *nodeStack) push(n *Node) {
	s.nodes = append(s.nodes, n)
}

func (s *nodeStack) pop() *Node {
	if len(s.nodes) == 0 {
		return nil
	}
	n := s.nodes[len(s.nodes)-1]
	s.nodes = s.nodes[:len(s.nodes)-1]
	return n
}

func (s *nodeStack) peek() *Node {
	if len(s.nodes) == 0 {
		return nil
	}
	return s.nodes[len(s.nodes)-1]
}

func (s *nodeStack) isEmpty() bool {
	return len(s.nodes) == 0
}
