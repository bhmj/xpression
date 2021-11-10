package xpression

import (
	"fmt"
)

func Parse(path []byte, i int) ([]*Token, error) {
	path = path[:trimSpaces(path)]
	l := len(path)

	// lexer
	tokens := make([]*Token, 0)
	var tok *Token
	var err error
	prevOperator := opPlus
	for i < l {
		i, tok, err = readNextToken(path, i, prevOperator)
		if err != nil {
			return nil, fmt.Errorf("%w at %d: %s", err, i, getLastWord(path[i:]))
		}
		if tok != nil {
			prevOperator = tok.Operator
			tokens = append(tokens, tok)
		}
	}

	//parser
	opStack := new(tokenStack)
	result := new(tokenStack)
	for _, token := range reverse(tokens) {
		switch token.Category {
		case tcLiteral, tcVariable:
			result.push(token)
		case tcOperator:
			for {
				top := opStack.peek()
				if top == nil {
					break
				}
				if top.Category == tcRightParenthesis {
					break
				}
				tokenPrecedence := operatorDetails[token.Operator].Precedence
				topPrecedence := operatorDetails[top.Operator].Precedence
				tokenAssociativity := operatorDetails[token.Operator].Associativity
				if tokenPrecedence < topPrecedence || (tokenPrecedence == topPrecedence && tokenAssociativity == aRight) {
					result.push(opStack.pop())
					continue
				}
				break
			}
			opStack.push(token)
		case tcRightParenthesis:
			opStack.push(token)
		case tcLeftParenthesis:
			for {
				top := opStack.peek()
				if top == nil {
					return nil, errMismatchedParentheses
				}
				if top.Category == tcRightParenthesis {
					opStack.pop()
					break
				}
				result.push(opStack.pop())
			}
		}
	}

	// pop the rest of the poerators
	for result.push(opStack.pop()) != nil {
	}

	return reverse(result.get()), nil
}

func readNextToken(path []byte, i int, prevOperator Operator) (int, *Token, error) {
	var err error
	i, err = skipSpaces(path, i)
	if err != nil {
		return i, nil, err
	}
	// parentheses
	if path[i] == '(' {
		return i + 1, &Token{Category: tcLeftParenthesis, Operator: opLeftParenthesis}, nil
	}
	if path[i] == ')' {
		return i + 1, &Token{Category: tcRightParenthesis, Operator: opRightParenthesis}, nil
	}
	// operator
	for _, op := range operatorSpelling {
		if matchSubslice(path[i:], op.Spelling) {
			if op.Code == opDivide && prevOperator != opNone {
				return readRegexp(path, i)
			}
			if op.Code == opMinus && prevOperator != opNone {
				op.Code = opUnaryMinus
			}
			return i + len(op.Spelling), &Token{Category: tcOperator, Operator: op.Code}, nil
		}
	}
	// literal:
	// number
	if path[i] >= '0' && path[i] <= '9' {
		return readNumber(path, i)
	}
	// string
	if path[i] == '"' || path[i] == '\'' {
		return readString(path, i)
	}
	// bool
	if path[i] == 't' || path[i] == 'f' {
		return readBool(path, i)
	}
	// null
	if path[i] == 'n' {
		return readNull(path, i)
	}
	if path[i] == '$' || path[i] == '@' {
		return readJsonpath(path, i)
	}

	return i, nil, errUnknownToken
}

func skipSpaces(input []byte, i int) (int, error) {
	l := len(input)
	for ; i < l; i++ {
		if !bytein(input[i], []byte{' ', ',', '\t', '\r', '\n'}) {
			break
		}
	}
	return i, nil
}

func trimSpaces(input []byte) int {
	l := len(input)
	i := l - 1
	for ; i >= 0; i-- {
		if !bytein(input[i], []byte{' ', ',', '\t', '\r', '\n'}) {
			break
		}
	}
	return i + 1
}

// returns true if b matches one of the elements of seq
func bytein(b byte, seq []byte) bool {
	for i := 0; i < len(seq); i++ {
		if b == seq[i] {
			return true
		}
	}
	return false
}
