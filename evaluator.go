//
// This package uses https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html for reference of expression evaluation logic.
//
package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

var (
	opsArithmetic = []byte{
		byte(opPlus),
		byte(opMinus),
		byte(opMultiply),
		byte(opDivide),
		byte(opExponentiation),
		byte(opBitwiseOR),
		byte(opBitwiseXOR),
		byte(opBitwiseAND),
		byte(opBitwiseNOT),
		byte(opShiftLeft),
		byte(opShiftRight),
	}

	opsLogic = []byte{
		byte(opLogicalAND),
		byte(opLogicalOR),
		byte(opLogicalNOT),
	}

	opsComparison = []byte{
		byte(opEqual),
		byte(opNotEqual),
		byte(opStrictEqual),
		byte(opStrictNotEqual),
		byte(opGE),
		byte(opG),
		byte(opLE),
		byte(opL),
		byte(opRegexMatch),
		byte(opNotRegexMatch),
	}
)

// EvaluateExpression evaluates the expression starting from the head of the token list.
// Usually it takes an operator from the head of the list and then takes 1 or 2 operands from the list,
// depending of the operator type (unary or binary).
// The extreme case is when there is only one operand exist in the list.
// This function calls itself recursively to evalute operands if needed.
func evaluateExpression(tokens []*Token) (*Operand, []*Token, error) {
	if len(tokens) == 0 {
		return nil, nil, errNotEnoughArguments
	}
	tok := tokens[0]
	if tok.Category == tcLiteral {
		return &tok.Operand, tokens[1:], nil
	}

	var (
		err   error
		left  *Operand
		right *Operand
	)
	left, tokens, err = evaluateExpression(tokens[1:])
	if err != nil {
		return nil, tokens, err
	}
	if operatorDetails[tok.Operator].Arguments > 1 {
		right, tokens, err = evaluateExpression(tokens)
		if err != nil {
			return nil, tokens, err
		}
	}

	result, err := execOperator(tok.Operator, left, right)
	return result, tokens, err
}

// execOperator takes an Operator and one or two Operands (the second one can be `nil` depending on operator type - unary or binary).
// It does evaluate the expression ("operand1 operator operand2" or "operator operand1") and return Operand which is a typed value.
func execOperator(op Operator, left *Operand, right *Operand) (*Operand, error) {
	var res Operand

	if bytes.IndexByte(opsArithmetic, byte(op)) != -1 {
		// arithmetic
		return doArithmetic(op, left, right)
	} else if bytes.IndexByte(opsComparison, byte(op)) != -1 {
		// comparison
		return doComparison(op, left, right)
	} else if bytes.IndexByte(opsLogic, byte(op)) != -1 {
		// logic
		return doLogic(op, left, right)
	}
	return &res, errUnknownToken
}

// doArithmetic actually evaluates the arithmetic operators.
// Note the special case of string concatenation: string + any_type -> string
func doArithmetic(op Operator, left *Operand, right *Operand) (*Operand, error) {
	var res Operand

	if op == opPlus && (left.Type|right.Type)&otString > 0 {
		// string concatenation
		res.Type = otString
		lstr := toString(left)
		rstr := toString(right)
		res.Str = append(lstr.Str, rstr.Str...)
		return &res, nil
	}

	res.Type = otNumber
	switch op {
	case opPlus:
		res.Number = toNumber(left).Number + toNumber(right).Number
	case opMinus:
		res.Number = toNumber(left).Number - toNumber(right).Number
	case opMultiply:
		res.Number = toNumber(left).Number * toNumber(right).Number
	case opDivide:
		res.Number = toNumber(left).Number / toNumber(right).Number
	case opExponentiation:
		res.Number = math.Pow(toNumber(left).Number, toNumber(right).Number)
	case opBitwiseOR:
		res.Number = float64(int64(toNumber(left).Number) | int64(toNumber(right).Number))
	case opBitwiseXOR:
		res.Number = float64(int64(toNumber(left).Number) ^ int64(toNumber(right).Number))
	case opBitwiseNOT:
		res.Number = float64(^int64(toNumber(left).Number))
	case opShiftLeft:
		res.Number = float64(int64(toNumber(left).Number) << int64(toNumber(right).Number))
	case opShiftRight:
		res.Number = float64(int64(toNumber(left).Number) >> int64(toNumber(right).Number))
	}
	return &res, nil
}

// doComaparison compares two operands. The special case is string regexp match which works only on strings.
// Otherwise works like JS comparison.
func doComparison(op Operator, left *Operand, right *Operand) (*Operand, error) {
	result := Operand{Type: otBoolean, Bool: false}

	if left.Type == otNull || right.Type == otNull {
		result.Bool = (left.Type|right.Type == otNull)
		return &result, nil
	}
	if op == opRegexMatch || op == opNotRegexMatch {
		// regexp match only works on string + regexp
		if !(left.Type|right.Type == otString|otRegexp) {
			return &result, nil
		}
	}
	if op == opStrictEqual && left.Type != right.Type {
		// strict comparison: types must match
		return &result, nil
	}
	comparedTypes := left.Type | right.Type

	if comparedTypes&otNumber > 0 {
		return doCompareNumber(op, toNumber(left), toNumber(right))
	}
	if comparedTypes&otBoolean > 0 {
		return doCompareBool(op, toBoolean(left), toBoolean(right))
	}
	if comparedTypes&otString > 0 {
		return doCompareString(op, left, right)
	}
	return &result, nil
}

// doLogic executes binary logical operators following JavaScript conversion rules.
func doLogic(op Operator, left *Operand, right *Operand) (*Operand, error) {
	var res Operand
	res.Type = otBoolean

	lval := toBoolean(left)
	if op == opLogicalAND || op == opLogicalOR {
		if (op == opLogicalAND && !lval.Bool) || (op == opLogicalOR && lval.Bool) { // false AND ..., true OR ... -> result!
			return lval, nil
		}
		return toBoolean(right), nil
	}

	// logical NOT
	lval.Bool = !lval.Bool
	return lval, nil
}

// doCompareBool compares two boolean values.
func doCompareBool(op Operator, left *Operand, right *Operand) (*Operand, error) {
	res := Operand{Type: otBoolean}

	if left.Type|right.Type != otBoolean {
		return &res, nil
	}
	switch op {
	case opG:
		res.Bool = (left.Bool && !right.Bool)
	case opL:
		res.Bool = (!left.Bool && right.Bool)
	case opEqual, opStrictEqual:
		res.Bool = (left.Bool == right.Bool)
	case opNotEqual, opStrictNotEqual:
		res.Bool = (left.Bool != right.Bool)
	case opGE:
		res.Bool = (left.Bool || !right.Bool)
	case opLE:
		res.Bool = (!left.Bool || right.Bool)
	}
	return &res, nil
}

// doCompareNumber compares two numbers.
func doCompareNumber(op Operator, left *Operand, right *Operand) (*Operand, error) {
	var res Operand
	res.Type = otBoolean
	switch op {
	case opG:
		res.Bool = left.Number > right.Number
	case opL:
		res.Bool = left.Number < right.Number
	case opEqual, opStrictEqual:
		res.Bool = left.Number == right.Number
	case opNotEqual, opStrictNotEqual:
		res.Bool = left.Number != right.Number
	case opGE:
		res.Bool = left.Number >= right.Number
	case opLE:
		res.Bool = left.Number <= right.Number
	}
	return &res, nil
}

// doCompareString compares two strings.
func doCompareString(op Operator, left *Operand, right *Operand) (*Operand, error) {
	var res Operand
	res.Type = otBoolean
	switch op {
	case opEqual, opStrictEqual:
		res.Bool = compareSlices(left.Str, right.Str) == 0
	case opG:
		res.Bool = compareSlices(left.Str, right.Str) > 0
	case opL:
		res.Bool = compareSlices(left.Str, right.Str) < 0
	case opGE:
		res.Bool = compareSlices(left.Str, right.Str) >= 0
	case opLE:
		res.Bool = compareSlices(left.Str, right.Str) <= 0
	case opNotEqual, opStrictNotEqual:
		res.Bool = compareSlices(left.Str, right.Str) != 0
	case opRegexMatch:
		res.Bool = right.Regexp.MatchString(string(left.Str))
	case opNotRegexMatch:
		res.Bool = !right.Regexp.MatchString(string(left.Str))
	default:
		return left, errUnknownToken
	}
	return &res, nil
}

// toString converts operand to string following JavaScript conversion rules.
func toString(op *Operand) *Operand {
	if op.Type == otString {
		return op
	}

	result := *op
	result.Type = otString

	switch op.Type {
	case otUndefined:
		result.Str = []byte("undefined")
	case otNull:
		result.Str = []byte("null")
	case otBoolean:
		if op.Bool {
			result.Str = []byte("true")
		} else {
			result.Str = []byte("false")
		}
	case otNumber:
		result.Str = []byte(fmt.Sprintf("%g", result.Number))
	}
	return &result
}

// toNumber converts operand to number following JavaScript conversion rules.
func toNumber(op *Operand) *Operand {
	if op.Type == otNumber {
		return op
	}

	result := *op
	result.Type = otNumber

	switch op.Type {
	case otUndefined:
		result.Number = math.NaN()
	case otNull:
		result.Number = 0
	case otBoolean:
		if op.Bool {
			result.Number = 1
		} else {
			result.Number = 0
		}
	case otString:
		var err error
		result.Number, err = strconv.ParseFloat(string(op.Str), 64)
		if err != nil {
			result.Number = math.NaN()
		}
	case otRegexp:
		result.Number = math.NaN()
	}
	return &result
}

// toBoolean converts operand to boolean following JavaScript conversion rules.
func toBoolean(op *Operand) *Operand {
	if op.Type == otBoolean {
		return op
	}

	result := *op
	result.Type = otBoolean
	switch op.Type {
	case otUndefined, otNull:
		result.Bool = false
	case otString:
		result.Bool = len(op.Str) > 0
	case otNumber:
		result.Bool = op.Number != 0 && !math.IsNaN(op.Number)
	}
	return &result
}

func compareSlices(s1 []byte, s2 []byte) int {
	if len(s1)+len(s2) == 0 {
		return 0
	}
	i := 0
	for i = 0; i < len(s1); i++ {
		if i > len(s2)-1 {
			return 1
		}
		if s1[i] != s2[i] {
			return int(s1[i]) - int(s2[i])
		}
	}
	if i < len(s2) {
		return -1
	}
	return 0
}
