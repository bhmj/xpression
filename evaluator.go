//
// [1] https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html used for reference of expression evaluation logic.
//
package xpression

import (
	"bytes"
	"math"
	"strconv"
)

var (
	opsArithmetic = []byte{
		byte(opUnaryMinus),
		byte(opPlus),
		byte(opMinus),
		byte(opMultiply),
		byte(opDivide),
		byte(opRemainder),
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

type VariableFunc func([]byte) (*Operand, error)

// Evaluate evaluates the previously parsed expression
func Evaluate(tokens []*Token, varFunc VariableFunc) (*Operand, error) {
	op, tokens, err := evaluate(tokens, varFunc)
	if err != nil {
		return nil, err
	}
	if len(tokens) > 0 {
		return nil, errNotEnoughArguments
	}
	return op, nil
}

// evaluate evaluates expression stored in `tokens` in prefix notation (NPN).
// Usually it takes an operator from the head of the list and then takes 1 or 2 operands from the list,
// depending of the operator type (unary or binary).
// The extreme case is when there is only one operand in the list.
// This function calls itself recursively to evalute operands if needed.
func evaluate(tokens []*Token, varFunc VariableFunc) (*Operand, []*Token, error) {
	if len(tokens) == 0 {
		return nil, nil, errNotEnoughArguments
	}
	tok := tokens[0]
	if tok.Category == tcLiteral {
		return &tok.Operand, tokens[1:], nil
	} else if tok.Category == tcVariable {
		if varFunc != nil {
			op, err := varFunc(tok.Str)
			return op, tokens[1:], err
		}
		return nil, tokens, errUnknownToken
	}

	var (
		err   error
		left  *Operand
		right *Operand
	)
	left, tokens, err = evaluate(tokens[1:], varFunc)
	if err != nil {
		return nil, tokens, err
	}
	if operatorDetails[tok.Operator].Arguments > 1 {
		right, tokens, err = evaluate(tokens, varFunc)
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
	return nil, errUnknownToken
}

// doArithmetic actually evaluates the arithmetic operators.
// Note the special case of string concatenation: string + any_type -> string
func doArithmetic(op Operator, left *Operand, right *Operand) (*Operand, error) {
	if op == opPlus && (left.Type|right.Type)&otString > 0 {
		// string concatenation
		toString(left)
		toString(right)
		left.Str = append(left.Str[:len(left.Str):len(left.Str)], right.Str...) // cannot use left buffer, must reallocate!
		return left, nil
	}

	switch op {
	case opUnaryMinus:
		left.Number = -toNumber(left).Number
	case opPlus:
		left.Number = toNumber(left).Number + toNumber(right).Number
	case opMinus:
		left.Number = toNumber(left).Number - toNumber(right).Number
	case opMultiply:
		left.Number = toNumber(left).Number * toNumber(right).Number
	case opDivide:
		left.Number = toNumber(left).Number / toNumber(right).Number
	case opRemainder:
		left.Number = float64(int64(toNumber(left).Number) % int64(toNumber(right).Number))
	case opExponentiation:
		left.Number = math.Pow(toNumber(left).Number, toNumber(right).Number)
	case opBitwiseAND:
		left.Number = float64(int64(toNumber(left).Number) & int64(toNumber(right).Number))
	case opBitwiseOR:
		left.Number = float64(int64(toNumber(left).Number) | int64(toNumber(right).Number))
	case opBitwiseXOR:
		left.Number = float64(int64(toNumber(left).Number) ^ int64(toNumber(right).Number))
	case opBitwiseNOT:
		left.Number = float64(^int64(toNumber(left).Number))
	case opShiftLeft:
		left.Number = float64(int64(toNumber(left).Number) << int64(toNumber(right).Number))
	case opShiftRight:
		left.Number = float64(int64(toNumber(left).Number) >> int64(toNumber(right).Number))
	}
	left.Type = otNumber
	return left, nil
}

// doComaparison compares two operands. The special case is string regexp match which works only on strings.
// Otherwise works like JS comparison.
func doComparison(op Operator, left *Operand, right *Operand) (*Operand, error) {
	result := left

	comparedTypes := left.Type | right.Type

	// [1] 7.2.14 (2,3)
	if op == opEqual && comparedTypes&(otNull|otUndefined) > 0 {
		// at least one side is null or undefined:
		result.Type = otBoolean
		result.Bool = (comparedTypes | otNull | otUndefined) == (otNull | otUndefined) // both are null or undefined
		return result, nil
	}
	if op == opRegexMatch || op == opNotRegexMatch {
		// one of them must be regexp
		if comparedTypes&otRegexp != otRegexp {
			result.Type = otBoolean
			result.Bool = false
			return result, nil
		}
		// other must be non-regexp
		if comparedTypes-otRegexp == 0 {
			result.Type = otBoolean
			result.Bool = false
			return result, nil
		}
		// convert non-regexp part to string and compare
		if right.Type == otRegexp {
			toString(left)
			return doCompareString(op, left, right) // regexp should be second argument
		}
		toString(right)
		return doCompareString(op, right, left) // regexp should be second argument
	}

	// [1] 7.2.15 (1)
	if (op == opStrictEqual || op == opStrictNotEqual) && left.Type != right.Type {
		// strict comparison: types must match
		result.Type = otBoolean
		result.Bool = false
		return result, nil
	}

	// [1] 7.2.15 (4)
	if comparedTypes == otString {
		return doCompareString(op, left, right)
	}

	// [1] 7.2.14 (5,6,7?,8?)
	return doCompareNumber(op, toNumber(left), toNumber(right))
}

// doLogic executes binary logical operators following JavaScript conversion rules.
func doLogic(op Operator, left *Operand, right *Operand) (*Operand, error) {
	lval := toBoolean(left)
	if op == opLogicalAND || op == opLogicalOR {
		if (op == opLogicalAND && !lval) || (op == opLogicalOR && lval) { // false AND ..., true OR ... -> result!
			return left, nil
		}
		return right, nil
	}

	// logical NOT
	left.Type = otBoolean
	left.Bool = !lval
	return left, nil
}

// doCompareNumber compares two numbers.
func doCompareNumber(op Operator, left *Operand, right *Operand) (*Operand, error) {
	res := left
	if math.IsNaN(left.Number) || math.IsNaN(right.Number) { // [1] 7.2.14 (4.h)
		res.Type = otBoolean
		return res, nil
	}
	if math.IsInf(left.Number, -1) || math.IsInf(right.Number, +1) { // [1] 7.2.14 (4.i)
		res.Type = otBoolean
		res.Bool = true
		return res, nil
	}
	if math.IsInf(left.Number, +1) || math.IsInf(right.Number, -1) { // [1] 7.2.14 (4.j)
		res.Type = otBoolean
		res.Bool = false
		return res, nil
	}
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
	res.Type = otBoolean
	return res, nil
}

// doCompareString compares two strings.
func doCompareString(op Operator, left *Operand, right *Operand) (*Operand, error) {
	res := left
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
		return res, errUnknownToken
	}
	res.Type = otBoolean
	return res, nil
}

// toString converts operand to string following JavaScript conversion rules.
func toString(op *Operand) *Operand {
	if op.Type == otString {
		return op
	}

	result := op

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
		result.Str = []byte(strconv.FormatFloat(op.Number, 'f', -1, 64))
	}
	result.Type = otString

	return result
}

// toNumber converts operand to number following JavaScript conversion rules. See [1] 7.1.4
func toNumber(op *Operand) *Operand {
	if op.Type == otNumber {
		return op
	}

	result := op

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
		if len(op.Str) == 0 { // [1] 7.1.4.1.2 (1)
			result.Number = 0
		} else {
			var err error
			result.Number, err = strconv.ParseFloat(string(op.Str), 64)
			if err != nil { // [1] 7.1.4.1.1 (3)
				result.Number = math.NaN()
			}
		}
	case otRegexp:
		result.Number = math.NaN()
	}
	result.Type = otNumber
	return result
}

// toBoolean converts operand to boolean following JavaScript conversion rules.
func toBoolean(op *Operand) bool {
	if op.Type == otBoolean {
		return op.Bool
	}

	result := false
	switch op.Type {
	case otUndefined, otNull, otRegexp:
	case otString:
		result = len(op.Str) > 0
	case otNumber:
		result = op.Number != 0 && !math.IsNaN(op.Number)
	}
	return result
}

// ToBoolean is a public alias
func ToBoolean(op *Operand) bool { return toBoolean(op) }

// logic complies to [1] 7.2.13 (3) only for BYTES not codepoints!
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
