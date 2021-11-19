//
// [1] https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html used for reference of expression evaluation logic.
//
package xpression

import (
	"bytes"
	"math"
	"regexp"
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

type VariableFunc func([]byte, *Operand) (error)

// Evaluate evaluates the previously parsed expression
func Evaluate(tokens []*Token, varFunc VariableFunc) (*Operand, error) {
	if len(tokens) == 0 {
		return nil, errNotEnoughArguments
	}
	op, tokens, err := evaluate(tokens, varFunc)
	if err != nil {
		return nil, err
	}
	if len(tokens) > 0 {
		return nil, errNotEnoughArguments
	}
	return op, nil
}

const (
	tokenOperand int = 0
	tokenResult int = 1
	tokensRest int = 2
)

// evaluate evaluates expression stored in `tokens` in prefix notation (NPN).
// Usually it takes an operator from the head of the list and then takes 1 or 2 operands from the list,
// depending of the operator type (unary or binary).
// The extreme case is when there is only one operand in the list.
// This function calls itself recursively to evalute operands if needed.
func evaluate(tokens []*Token, varFunc VariableFunc) (*Operand, []*Token, error) {
	if len(tokens) == 0 {
		return nil, nil, errNotEnoughArguments
	}
	tok := tokens[tokenOperand]
	if tok.Category == tcLiteral {
		return &tok.Operand, tokens[tokenResult:], nil
	} else if tok.Category == tcVariable {
		if varFunc != nil {
			err := varFunc(tok.Str, &tokens[tokenResult].Operand)
			return &tokens[tokenResult].Operand, tokens[tokensRest:], err
		}
		return nil, tokens, errUnknownToken
	}

	var (
		err   error
		left  *Operand
		right *Operand
		result *Operand
	)
	result = &tokens[tokenResult].Operand
	left, tokens, err = evaluate(tokens[tokensRest:], varFunc)
	if err != nil {
		return nil, tokens, err
	}
	if operatorDetails[tok.Operator].Arguments > 1 {
		right, tokens, err = evaluate(tokens, varFunc)
		if err != nil {
			return nil, tokens, err
		}
	}

	err = execOperator(tok.Operator, left, right, result)
	return result, tokens, err
}

// execOperator takes an Operator and one or two Operands (the second one can be `nil` depending on operator type - unary or binary).
// It does evaluate the expression ("operand1 operator operand2" or "operator operand1") and return Operand which is a typed value.
func execOperator(op Operator, left *Operand, right *Operand, result *Operand) error {
	if bytes.IndexByte(opsArithmetic, byte(op)) != -1 {
		// arithmetic
		return doArithmetic(op, left, right, result)
	} else if bytes.IndexByte(opsComparison, byte(op)) != -1 {
		// comparison
		return doComparison(op, left, right, result)
	} else if bytes.IndexByte(opsLogic, byte(op)) != -1 {
		// logic
		return doLogic(op, left, right, result)
	}
	return errUnknownToken
}

// doArithmetic actually evaluates the arithmetic operators.
// Note the special case of string concatenation: string + any_type -> string
func doArithmetic(op Operator, left *Operand, right *Operand, result *Operand) error {
	if op == opPlus && (left.Type|right.Type)&otString > 0 {
		// string concatenation
		lval := toString(left)
		rval := toString(right)
		result.Type = otString
		result.Str = append(lval[:len(lval):len(lval)], rval...) // cannot use left buffer, must reallocate!
		return nil
	}

	switch op {
	case opUnaryMinus:
		result.Number = -toNumber(left)
	case opPlus:
		result.Number = toNumber(left) + toNumber(right)
	case opMinus:
		result.Number = toNumber(left) - toNumber(right)
	case opMultiply:
		result.Number = toNumber(left) * toNumber(right)
	case opDivide:
		result.Number = toNumber(left) / toNumber(right)
	case opRemainder:
		result.Number = float64(int64(toNumber(left)) % int64(toNumber(right)))
	case opExponentiation:
		result.Number = math.Pow(toNumber(left), toNumber(right))
	case opBitwiseAND:
		result.Number = float64(int64(toNumber(left)) & int64(toNumber(right)))
	case opBitwiseOR:
		result.Number = float64(int64(toNumber(left)) | int64(toNumber(right)))
	case opBitwiseXOR:
		result.Number = float64(int64(toNumber(left)) ^ int64(toNumber(right)))
	case opBitwiseNOT:
		result.Number = float64(^int64(toNumber(left)))
	case opShiftLeft:
		result.Number = float64(int64(toNumber(left)) << int64(toNumber(right)))
	case opShiftRight:
		result.Number = float64(int64(toNumber(left)) >> int64(toNumber(right)))
	}
	result.Type = otNumber
	return nil
}

// doComaparison compares two operands. The special case is string regexp match which works only on strings.
// Otherwise works like JS comparison.
func doComparison(op Operator, left *Operand, right *Operand, result *Operand) error {
	comparedTypes := left.Type | right.Type
	result.Type = otBoolean

	// [1] 7.2.14 (2,3)
	if op == opEqual && comparedTypes&(otNull|otUndefined) > 0 {
		// at least one side is null or undefined:
		result.Bool = (comparedTypes | otNull | otUndefined) == (otNull | otUndefined) // both are null or undefined
		return nil
	}
	if op == opRegexMatch || op == opNotRegexMatch {
		// one of them must be regexp
		if comparedTypes&otRegexp != otRegexp {
			result.Bool = false
			return nil
		}
		// other must be non-regexp
		if comparedTypes-otRegexp == 0 {
			result.Bool = false
			return nil
		}
		// convert non-regexp part to string and compare
		if right.Type == otRegexp {
			lval := toString(left)
			return doCompareRegexp(op, lval, right.Regexp, result) // regexp should be second argument
		}
		rval := toString(right)
		return doCompareRegexp(op, rval, left.Regexp, result) // regexp should be second argument
	}

	// [1] 7.2.15 (1)
	if (op == opStrictEqual || op == opStrictNotEqual) && left.Type != right.Type {
		// strict comparison: types must match
		result.Bool = false
		return nil
	}

	// [1] 7.2.15 (4)
	if comparedTypes == otString {
		return doCompareString(op, toString(left), toString(right), result)
	}

	// [1] 7.2.14 (5,6,7?,8?)
	return doCompareNumber(op, toNumber(left), toNumber(right), result)
}

// doLogic executes binary logical operators following JavaScript conversion rules.
func doLogic(op Operator, left *Operand, right *Operand, result *Operand) error {
	lval := toBoolean(left)
	if op == opLogicalAND || op == opLogicalOR {
		if (op == opLogicalAND && !lval) || (op == opLogicalOR && lval) { // false AND ..., true OR ... -> result!
			*result = *left
			return nil
		}
		*result = *right
		return nil
	}

	// logical NOT
	result.Type = otBoolean
	result.Bool = !lval
	return nil
}

// doCompareNumber compares two numbers.
func doCompareNumber(op Operator, left float64, right float64, result *Operand) error {
	result.Type = otBoolean
	if math.IsNaN(left) || math.IsNaN(right) { // [1] 7.2.14 (4.h)
		result.Bool = false
		return nil
	}
	if math.IsInf(left, -1) || math.IsInf(right, +1) { // [1] 7.2.14 (4.i)
		if sameInfinities(left, right) {
			result.Bool = opEquality(op)
		} else {
			result.Bool = opLessNotEqual(op)
		}
		return nil
	}
	if math.IsInf(left, +1) || math.IsInf(right, -1) { // [1] 7.2.14 (4.j)
		// NB: same infinities processed above
		result.Bool = opGreaterNotEqual(op)
		return nil
	}
	switch op {
	case opG:
		result.Bool = left > right
	case opL:
		result.Bool = left < right
	case opEqual, opStrictEqual:
		result.Bool = left == right
	case opNotEqual, opStrictNotEqual:
		result.Bool = left != right
	case opGE:
		result.Bool = left >= right
	case opLE:
		result.Bool = left <= right
	}
	return nil
}

func opEquality(op Operator) bool {
	return op==opEqual || op==opStrictEqual || op==opLE || op==opGE
}

func opLessNotEqual(op Operator) bool {
	return op==opL || op==opLE || op==opNotEqual || op==opStrictNotEqual
}

func opGreaterNotEqual(op Operator) bool {
	return op==opG || op==opGE || op==opNotEqual || op==opStrictNotEqual
}

func sameInfinities(left, right float64) bool {
	return math.IsInf(left, -1) && math.IsInf(right, -1) || math.IsInf(left, +1) && math.IsInf(right, +1)
}

// doCompareString compares two strings.
func doCompareString(op Operator, left []byte, right []byte, result *Operand) error {
	result.Type = otBoolean
	switch op {
	case opEqual, opStrictEqual:
		result.Bool = compareSlices(left, right) == 0
	case opG:
		result.Bool = compareSlices(left, right) > 0
	case opL:
		result.Bool = compareSlices(left, right) < 0
	case opGE:
		result.Bool = compareSlices(left, right) >= 0
	case opLE:
		result.Bool = compareSlices(left, right) <= 0
	case opNotEqual, opStrictNotEqual:
		result.Bool = compareSlices(left, right) != 0
	}
	return nil
}

// doCompareRegexp matches a string to regexp
func doCompareRegexp(op Operator, left []byte, right *regexp.Regexp, result *Operand) error {
	result.Type = otBoolean
	if op == opRegexMatch {
		result.Bool = right.MatchString(string(left))
	} else {
		result.Bool = !right.MatchString(string(left))
	}
	return nil
}

// toString converts operand to string following JavaScript conversion rules.
func toString(op *Operand) []byte {
	if op.Type == otString {
		return op.Str
	}

	switch op.Type {
	case otUndefined:
		return []byte("undefined")
	case otNull:
		return []byte("null")
	case otBoolean:
		if op.Bool {
			return []byte("true")
		} else {
			return []byte("false")
		}
	case otNumber:
		return []byte(strconv.FormatFloat(op.Number, 'f', -1, 64))
	}

	return nil // not reaching here
}

// toNumber converts operand to number following JavaScript conversion rules. See [1] 7.1.4
func toNumber(op *Operand) float64 {
	if op.Type == otNumber {
		return op.Number
	}

	switch op.Type {
	case otUndefined:
		return math.NaN()
	case otNull:
		return 0
	case otBoolean:
		if op.Bool {
			return 1
		} else {
			return 0
		}
	case otString:
		if len(op.Str) == 0 { // [1] 7.1.4.1.2 (1)
			return 0
		} else {
			f, err := strconv.ParseFloat(string(op.Str), 64)
			if err != nil { // [1] 7.1.4.1.1 (3)
				return math.NaN()
			}
			return f
		}
	case otRegexp:
		return math.NaN()
	}
	return 0 // not reaching here
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
