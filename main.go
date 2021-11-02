package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Operator byte      // list of operators: + - * / < > == !=
type TokenCategory byte // operator, literal (operand), parentheses
type OperandType byte   // string, number, boolean, null, undefined
type Associativity byte // left, right

const (
	opNone             Operator = '\x00'
	opLogicalOR        Operator = 'O'
	opLogicalAND       Operator = 'A'
	opBitwiseOR        Operator = '|'
	opBitwiseXOR       Operator = '^'
	opBitwiseAND       Operator = '&'
	opEqual            Operator = 'E'
	opStrictEqual      Operator = 'e'
	opNotEqual         Operator = 'N'
	opStrictNotEqual   Operator = 'n'
	opGE               Operator = 'G'
	opG                Operator = 'g'
	opLE               Operator = 'L'
	opL                Operator = 'l'
	opRegexMatch       Operator = 'R'
	opNotRegexMatch    Operator = 'r'
	opShiftRight       Operator = '>'
	opShiftLeft        Operator = '<'
	opPlus             Operator = '+'
	opMinus            Operator = '-'
	opMultiply         Operator = '*'
	opDivide           Operator = '/'
	opExponentiation   Operator = 'P'
	opLogicalNOT       Operator = '!'
	opBitwiseNOT       Operator = '~'
	opUnaryMinus       Operator = '_'
	opLeftParenthesis  Operator = '('
	opRightParenthesis Operator = ')'
)

const (
	tcLiteral          TokenCategory = 1 << iota // string, number, bool
	tcOperator                                   // +-*/^!<=>
	tcLeftParenthesis                            //
	tcRightParenthesis                           //
	tcNodeRef                                    // @.key
)

const (
	otString OperandType = 1 << iota
	otNumber
	otBoolean
	otNull
	otUndefined
	otRegexp
)

const (
	aLeft Associativity = iota
	aRight
)

type OperatorDetail struct {
	Associativity Associativity
	Precedence    int
	Arguments     int
}

var operatorSpelling = []struct {
	Spelling []byte
	Code     Operator
}{ // IMPORTANT: first longest string, then substring(s)! Ex: "!~=", "!~", "!"
	{[]byte("||"), opLogicalOR},
	{[]byte("&&"), opLogicalAND},
	{[]byte("|"), opBitwiseOR},
	{[]byte("&"), opBitwiseAND},
	{[]byte("^"), opBitwiseXOR},
	{[]byte("==="), opStrictEqual},
	{[]byte("=="), opEqual},
	{[]byte("!=="), opStrictNotEqual},
	{[]byte("!="), opNotEqual},
	{[]byte(">>"), opShiftRight},
	{[]byte("<<"), opShiftLeft},
	{[]byte(">="), opGE},
	{[]byte(">"), opG},
	{[]byte("<="), opLE},
	{[]byte("<"), opL},
	{[]byte("~="), opRegexMatch},
	{[]byte("!~="), opNotRegexMatch},
	{[]byte("!~"), opNotRegexMatch},
	{[]byte("**"), opExponentiation},
	{[]byte("+"), opPlus},
	{[]byte("-"), opMinus},
	{[]byte("*"), opMultiply},
	{[]byte("/"), opDivide},
	{[]byte("!"), opLogicalNOT},
	{[]byte("~"), opBitwiseNOT},
	{[]byte("-"), opUnaryMinus},
	{[]byte("("), opLeftParenthesis},
	{[]byte(")"), opRightParenthesis},
}

var operatorDetails = map[Operator]OperatorDetail{
	opLogicalOR:        {aLeft, 1, 2},   // logical OR
	opLogicalAND:       {aLeft, 2, 2},   // logical AND
	opBitwiseOR:        {aLeft, 3, 2},   // bitwise OR
	opBitwiseXOR:       {aLeft, 4, 2},   // bitwise XOR
	opBitwiseAND:       {aLeft, 5, 2},   // bitwise AND
	opEqual:            {aLeft, 6, 2},   // ==
	opStrictEqual:      {aLeft, 6, 2},   // ===
	opNotEqual:         {aLeft, 6, 2},   // !=
	opStrictNotEqual:   {aLeft, 6, 2},   // !==
	opGE:               {aLeft, 7, 2},   // >=
	opG:                {aLeft, 7, 2},   // >
	opLE:               {aLeft, 7, 2},   // <=
	opL:                {aLeft, 7, 2},   // <
	opRegexMatch:       {aLeft, 7, 2},   // ~=
	opNotRegexMatch:    {aLeft, 7, 2},   // !~=
	opShiftRight:       {aLeft, 8, 2},   // >>
	opShiftLeft:        {aLeft, 8, 2},   // <<
	opPlus:             {aLeft, 9, 2},   // +
	opMinus:            {aLeft, 9, 2},   // -
	opMultiply:         {aLeft, 10, 2},  // *
	opDivide:           {aLeft, 10, 2},  // /
	opExponentiation:   {aRight, 11, 2}, // **
	opLogicalNOT:       {aLeft, 12, 1},  // logical NOT (!)
	opBitwiseNOT:       {aLeft, 12, 1},  // bitwise NOT (~)
	opUnaryMinus:       {aLeft, 12, 1},  // unary -
	opLeftParenthesis:  {aLeft, 13, 2},  // (
	opRightParenthesis: {aLeft, 13, 2},  // )
}

type Operand struct {
	Type   OperandType
	Str    []byte
	Number float64
	Bool   bool
	Regexp *regexp.Regexp
	// + node reference
}

type Token struct {
	Category TokenCategory
	Operator Operator
	Operand
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Expression parser.\nUsage: %[1]s <expression>", filepath.Base(os.Args[0]))
		return
	}

	tokens, err := parseExpression([]byte(os.Args[1]), 0)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		printParsedExpression(tokens)

		result, _, err := evaluateExpression(tokens)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			printParsedExpression([]*Token{{Category: tcLiteral, Operand: *result}})
		}
	}
}

func parseExpression(path []byte, i int) ([]*Token, error) {
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
		case tcLiteral:
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
		case tcNodeRef:
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
			if op.Code == opDivide && (prevOperator == opRegexMatch || prevOperator == opNotRegexMatch) {
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

	return i, nil, errUnknownToken
}

func skipSpaces(input []byte, i int) (int, error) {
	l := len(input)
	for ; i < l; i++ {
		if !bytein(input[i], []byte{' ', ',', '\t', '\r', '\n'}) {
			break
		}
	}
	if i == l {
		return i, errUnexpectedEnd
	}
	return i, nil
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

func printParsedExpression(tokens []*Token) {
	for _, tok := range tokens {
		switch tok.Category {
		case tcLiteral:
			switch tok.Type {
			case otString:
				fmt.Printf("\"%s\"", string(tok.Str))
			case otNumber:
				fmt.Printf("%v", tok.Number)
			case otBoolean:
				fmt.Printf("%v", tok.Bool)
			case otRegexp:
				fmt.Printf("/%s/", tok.Regexp.String())
			}
		case tcOperator:
			fmt.Printf("%s", string([]byte{byte(tok.Operator)}))
		}
		fmt.Printf(" ")
	}
	fmt.Println()
}
