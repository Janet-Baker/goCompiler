/**
 * ============================================================================
 *                                   (/^▽^)/
 *                                THE TOKENIZER!
 * ============================================================================
 */

package main

import (
	"errors"
	"fmt"
)

// token type iota
const (
	tInvalid = iota
	tNewLine // \n
	tString  // ".*"
	tInteger // [0-9]+
	tDot     // "."
	tComma   // ","
	tBreak   // ";"
	tLParen  // "("
	tRParen  // ")"
	tLBrace  // "{"
	tRBrace  // "}"
	// Operators and functions
	// binary operators
	tPlus     // "+"
	tMinus    // "-"
	tMultiply // "*"
	tDivide   // "/"
	// compare operators
	tCalcNotEqual     // "!="
	tCalcLessThan     // "<"
	tCalcLessEqual    // "<="
	tCalcGreaterThan  // ">"
	tCalcGreaterEqual // ">="
	tCalcEqual        // "=="
	tEqual            // "="
	// keywords Statement
	tReturn     // "return"
	tIf         // "if"
	tElse       // "else"
	tFor        // "for"
	tWhile      // "while"
	tPrint      // "print"
	tIdentifier // [a-zA-Z_][a-zA-Z0-9_]*
)

type token struct {
	kind  int
	value string
	line  int
	col   int
}

// lex and tokenize the content
func tokenize(content []byte) (tokens []token, err error) {
	tokens = make([]token, len(content))
	i := 0
	line := 1
	col := 1
	for currPos := 0; currPos < len(content); currPos++ {
		switch {
		// "\n"                    SAVE_TOKEN; return tNewLine;
		case content[currPos] == '\n':
			if content[currPos-1] == '\r' {
				tokens[i] = token{tNewLine, "\\n", line, col - 1}
			} else {
				tokens[i] = token{tNewLine, "\\n", line, col}
				i++
			}
			line++
			col = 1
			break

		// '\"'                     SAVE_TOKEN; return tString;
		case content[currPos] == '"':
			t := token{tString, "", line, col}
			// read to the next "
			targetPos := currPos + 1
			for targetPos < len(content) && content[targetPos] != '"' {
				if content[targetPos] == '\n' {
					fmt.Printf("invalid token: not paired \" at pos %d, line %d, col %d", currPos, line, col)
					line++
					col = 1
				}
				targetPos++
			}
			t.value = string(content[currPos:targetPos])
			col = col + targetPos - currPos
			currPos = targetPos
			tokens[i] = t
			i++
			break

		// blank				   Skip;
		case content[currPos] == ' ' || content[currPos] == '\t' || content[currPos] == '\r':
			col++
			break

		// [0-9]+                  SAVE_TOKEN; return tInteger;
		case (content[currPos] >= '0') && (content[currPos] <= '9'):
			// store the number string to value
			// tokens[i] = token{TINTEGER, value, line, col}
			tokens[i] = token{tInteger, "", line, col}
			targetPos := currPos + 1
			for targetPos < len(content) && content[targetPos] >= '0' && content[targetPos] <= '9' {
				targetPos++
			}
			tokens[i].value = string(content[currPos:targetPos])
			i++
			col = col + targetPos - currPos
			currPos = targetPos - 1
			break

		// "."                     return TOKEN(tDot);
		case content[currPos] == '.':
			tokens[i] = token{tDot, ".", line, col}
			i++
			col++
			break

		// ","                     return TOKEN(tComma);
		case content[currPos] == ',':
			tokens[i] = token{tComma, ",", line, col}
			i++
			col++
			break

		// ";"                     return TOKEN(tBreak);
		case content[currPos] == ';':
			tokens[i] = token{tBreak, ";", line, col}
			i++
			col++
			break

		// "+"                     return TOKEN(tPlus);
		case content[currPos] == '+':
			tokens[i] = token{tPlus, "+", line, col}
			i++
			col++
			break

		// "-"                     return TOKEN(tMinus);
		case content[currPos] == '-':
			tokens[i] = token{tMinus, "-", line, col}
			i++
			col++
			break

		// "*"                     return TOKEN(tMultiple);
		case content[currPos] == '*':
			tokens[i] = token{tMultiply, "*", line, col}
			i++
			col++
			break

		// "/"                     return TOKEN(tDivide);
		case content[currPos] == '/':
			tokens[i] = token{tDivide, "/", line, col}
			i++
			col++
			break

		// "("                     return TOKEN(tLParen);
		case content[currPos] == '(':
			tokens[i] = token{tLParen, "(", line, col}
			i++
			col++
			break

		// ")"                     return TOKEN(tRParen);
		case content[currPos] == ')':
			tokens[i] = token{tRParen, ")", line, col}
			i++
			col++
			break

		// "{"                     return TOKEN(tLBrace);
		case content[currPos] == '{':
			tokens[i] = token{tLBrace, "{", line, col}
			i++
			col++
			break

		// "}"                     return TOKEN(tRBrace);
		case content[currPos] == '}':
			tokens[i] = token{tRBrace, "}", line, col}
			i++
			col++
			break

		//"!="                    return TOKEN(tCalcNotEqual);
		case content[currPos] == '!':
			if content[currPos+1] == '=' {
				tokens[i] = token{tCalcNotEqual, "!=", line, col}
				i++
				col = col + 2
				currPos++
			}
			fmt.Printf("invalid token at pos %d, line %d, col %d", currPos, line, col)
			break

		//"<"                     return TOKEN(tCalcLessThan);
		//"<="                    return TOKEN(tCalcLessEqual);
		case content[currPos] == '<':
			if content[currPos+1] == '=' {
				tokens[i] = token{tCalcLessEqual, "<=", line, col}
				i++
				col = col + 2
				currPos++
			} else {
				tokens[i] = token{tCalcLessThan, "<", line, col}
				i++
				col++
			}
			break

		//">"                     return TOKEN(tCalcGreaterThan);
		//">="                    return TOKEN(tCalcGreaterEqual);
		case content[currPos] == '>':
			if content[currPos+1] == '=' {
				tokens[i] = token{tCalcGreaterEqual, ">=", line, col}
				i++
				col = col + 2
				currPos++
			} else {
				tokens[i] = token{tCalcGreaterThan, ">", line, col}
				i++
				col++
			}
			break

		//"=="                    return TOKEN(tCalcEqual);
		//"="                     return TOKEN(tEqual);
		case content[currPos] == '=':
			if content[currPos+1] == '=' {
				tokens[i] = token{tCalcEqual, "==", line, col}
				i++
				col = col + 2
				currPos = currPos + 1
			} else {
				tokens[i] = token{tEqual, "=", line, col}
				i++
				col++
			}
			break

		// "return"                return TOKEN(tReturn);
		case content[currPos] == 'r' && content[currPos+1] == 'e' && content[currPos+2] == 't' && content[currPos+3] == 'u' && content[currPos+4] == 'r' && content[currPos+5] == 'n':
			tokens[i] = token{tReturn, "return", line, col}
			i++
			col = col + 6
			currPos = currPos + 5
			break

		// "if"                    return TOKEN(tIf);
		case content[currPos] == 'i' && content[currPos+1] == 'f':
			tokens[i] = token{tIf, "if", line, col}
			i++
			col = col + 2
			currPos = currPos + 1
			break

		// "else"                  return TOKEN(tElse);
		case content[currPos] == 'e' && content[currPos+1] == 'l' && content[currPos+2] == 's' && content[currPos+3] == 'e':
			tokens[i] = token{tElse, "else", line, col}
			i++
			col = col + 4
			currPos = currPos + 3
			break

		// "for"                   return TOKEN(tFor);
		case content[currPos] == 'f' && content[currPos+1] == 'o' && content[currPos+2] == 'r':
			tokens[i] = token{tFor, "for", line, col}
			i++
			col = col + 3
			currPos = currPos + 2
			break

		// "while"                 return TOKEN(tWhile);
		case content[currPos] == 'w' && content[currPos+1] == 'h' && content[currPos+2] == 'i' && content[currPos+3] == 'l' && content[currPos+4] == 'e':
			tokens[i] = token{tWhile, "while", line, col}
			i++
			col = col + 5
			currPos = currPos + 4
			break

		//	[a-zA-Z_][a-zA-Z0-9_]*  SAVE_TOKEN; return tIdentifier;
		case content[currPos] >= 'a' && content[currPos] <= 'z' || content[currPos] >= 'A' && content[currPos] <= 'Z' || content[currPos] == '_':
			targetPos := currPos + 1
			for targetPos < len(content) && ((content[targetPos] >= 'a' && content[targetPos] <= 'z') || (content[targetPos] >= 'A' && content[targetPos] <= 'Z') || (content[targetPos] == '_') || (content[targetPos] >= '0' && content[targetPos] <= '9')) {
				targetPos++
			}
			value := string(content[currPos:targetPos])
			tokens[i] = token{tIdentifier, value, line, col}
			i++
			col = col + targetPos - currPos
			currPos = targetPos - 1
			break

		default:
			fmt.Printf("invalid token at pos %d, line %d, col %d", currPos, line, col)
			return nil, errors.New("invalid token")
		}
	}
	return tokens[0:i], nil
}
