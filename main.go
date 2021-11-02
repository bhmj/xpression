package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Operator byte

type TokenCategory byte
type LiteralType byte

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

	tcLiteral          TokenCategory = iota // string, number, bool
	tcOperator                              // +-*/^!<=>
	tcLeftParenthesis                       //
	tcRightParenthesis                      //
	tcNodeRef                               // @.key

	ltString LiteralType = iota
	ltNumber
	ltBoolean
	ltNull
	ltUndefined
	ltRegexp
)

type Associativity byte

const (
	aLeft Associativity = iota
	aRight
)

type OperatorDetail struct {
	Associativity Associativity
	Precedence    int
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
	opLogicalOR:        {aLeft, 1},   // logical OR
	opLogicalAND:       {aLeft, 2},   // logical AND
	opBitwiseOR:        {aLeft, 3},   // bitwise OR
	opBitwiseXOR:       {aLeft, 4},   // bitwise XOR
	opBitwiseAND:       {aLeft, 5},   // bitwise AND
	opEqual:            {aLeft, 6},   // ==
	opStrictEqual:      {aLeft, 6},   // ===
	opNotEqual:         {aLeft, 6},   // !=
	opStrictNotEqual:   {aLeft, 6},   // !==
	opGE:               {aLeft, 7},   // >=
	opG:                {aLeft, 7},   // >
	opLE:               {aLeft, 7},   // <=
	opL:                {aLeft, 7},   // <
	opRegexMatch:       {aLeft, 7},   // ~=
	opNotRegexMatch:    {aLeft, 7},   // !~=
	opShiftRight:       {aLeft, 8},   // >>
	opShiftLeft:        {aLeft, 8},   // <<
	opPlus:             {aLeft, 9},   // +
	opMinus:            {aLeft, 9},   // -
	opMultiply:         {aLeft, 10},  // *
	opDivide:           {aLeft, 10},  // /
	opExponentiation:   {aRight, 11}, // **
	opLogicalNOT:       {aLeft, 12},  // logical NOT (!)
	opBitwiseNOT:       {aLeft, 12},  // bitwise NOT (~)
	opUnaryMinus:       {aLeft, 12},  // unary -
	opLeftParenthesis:  {aLeft, 13},  // (
	opRightParenthesis: {aLeft, 13},  // )
}

type Token struct {
	Category TokenCategory
	Operator Operator
	Type     LiteralType
	Str      []byte
	Number   float64
	Bool     bool
	Regexp   *regexp.Regexp
	// + node reference
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
			return nil, fmt.Errorf("%w at %d: %c", err, i, path[i])
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
					return nil, errMismatchedParenthesis
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
	// parenthesis
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

	return i, nil, errInvalidCharacter
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
			case ltString:
				fmt.Printf("\"%s\"", string(tok.Str))
			case ltNumber:
				fmt.Printf("%v", tok.Number)
			case ltBoolean:
				fmt.Printf("%v", tok.Bool)
			case ltRegexp:
				fmt.Printf("/%s/", tok.Regexp.String())
			}
		case tcOperator:
			fmt.Printf("%s", string([]byte{byte(tok.Operator)}))
		}
		fmt.Printf(" ")
	}
	fmt.Println()
}
