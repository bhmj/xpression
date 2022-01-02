package xpression

import (
	"fmt"
	"regexp"
	"strconv"
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
	opRemainder        Operator = '%'
	opExponentiation   Operator = 'P'
	opLogicalNOT       Operator = '!'
	opBitwiseNOT       Operator = '~'
	opUnaryMinus       Operator = '_'
	opLeftParenthesis  Operator = '('
	opRightParenthesis Operator = ')'
)

const (
	tcIntermediateResult TokenCategory = 0         // intermediate result placeholder
	tcLiteral            TokenCategory = 1 << iota // string, number, bool
	tcOperator                                     // +-*/^!<=>
	tcLeftParenthesis                              //
	tcRightParenthesis                             //
	tcVariable                                     // @.key etc
)

const (
	otString OperandType = 1 << iota
	otNumber
	otBoolean
	otNull
	otUndefined
	otRegexp
	otVariable
)

const (
	// public aliases
	StringOperand    = otString
	NumberOperand    = otNumber
	BooleanOperand   = otBoolean
	NullOperand      = otNull
	UndefinedOperand = otUndefined
	RegexpOperand    = otRegexp
	VariableOperand  = otVariable
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
}{ // IMPORTANT: first longest string, then substring(s)! Ex: "!=~", "!=", "!"
	{[]byte("||"), opLogicalOR},
	{[]byte("&&"), opLogicalAND},
	{[]byte("|"), opBitwiseOR},
	{[]byte("&"), opBitwiseAND},
	{[]byte("^"), opBitwiseXOR},
	{[]byte("!=~"), opNotRegexMatch},
	{[]byte("!~"), opNotRegexMatch},
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
	{[]byte("=~"), opRegexMatch},
	{[]byte("**"), opExponentiation},
	{[]byte("+"), opPlus},
	{[]byte("-"), opMinus},
	{[]byte("*"), opMultiply},
	{[]byte("/"), opDivide},
	{[]byte("%"), opRemainder},
	{[]byte("!"), opLogicalNOT},
	{[]byte("~"), opBitwiseNOT},
	{[]byte("-"), opUnaryMinus},
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
	opRegexMatch:       {aLeft, 7, 2},   // =~
	opNotRegexMatch:    {aLeft, 7, 2},   // !=~
	opShiftRight:       {aLeft, 8, 2},   // >>
	opShiftLeft:        {aLeft, 8, 2},   // <<
	opPlus:             {aLeft, 9, 2},   // +
	opMinus:            {aLeft, 9, 2},   // -
	opMultiply:         {aLeft, 10, 2},  // *
	opDivide:           {aLeft, 10, 2},  // /
	opRemainder:        {aLeft, 10, 2},  // %
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

var operatorBound []byte

func init() {
	operatorBound = []byte{' ', '['}
	for _, operator := range operatorSpelling {
		if !bytein(operator.Spelling[0], operatorBound) {
			operatorBound = append(operatorBound, operator.Spelling[0])
		}
	}
}

func (op *Operand) String() string {
	switch op.Type {
	case otNull:
		return "null"
	case otUndefined:
		return "undefined"
	case otString:
		return fmt.Sprintf("\"%s\"", string(op.Str))
	case otNumber:
		return strconv.FormatFloat(op.Number, 'f', -1, 64)
	case otBoolean:
		return fmt.Sprintf("%v", op.Bool)
	case otRegexp:
		return fmt.Sprintf("/%s/", op.Regexp.String())
	}
	return "???"
}

type Token struct {
	Category TokenCategory
	Operator Operator
	Operand
}

func (tok *Token) String() string {
	opCode := func(op Operator) string {
		for _, rec := range operatorSpelling {
			if rec.Code == op {
				return string(rec.Spelling)
			}
		}
		return "???"
	}

	switch tok.Category {
	case tcIntermediateResult:
		return "IR"
	case tcLiteral:
		return tok.Operand.String()
	case tcOperator:
		return opCode(tok.Operator)
	case tcVariable:
		return string(tok.Str)
	}

	return "unknown"
}

func (op *Operand) SetString(s string) {
	op.Type = otString
	op.Str = []byte(s)
}

func (op *Operand) SetNumber(f float64) {
	op.Type = otNumber
	op.Number = f
}

func (op *Operand) SetBoolean(b bool) {
	op.Type = otBoolean
	op.Bool = b
}

func (op *Operand) SetNull() {
	op.Type = otNull
}

func (op *Operand) SetUndefined() {
	op.Type = otUndefined
}

func (op *Operand) SetRegexp(r *regexp.Regexp) {
	op.Type = otRegexp
	op.Regexp = r
}

func String(s string) *Operand {
	return &Operand{Type: otString, Str: []byte(s)}
}

func Number(f float64) *Operand {
	return &Operand{Type: otNumber, Number: f}
}

func Boolean(b bool) *Operand {
	return &Operand{Type: otBoolean, Bool: b}
}

func Null() *Operand {
	return &Operand{Type: otNull}
}

func Undefined() *Operand {
	return &Operand{Type: otUndefined}
}

func Regexp(r *regexp.Regexp) *Operand {
	return &Operand{Type: otRegexp, Regexp: r}
}
