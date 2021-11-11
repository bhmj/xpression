package xpression

import (
	"regexp"
	"strconv"
)

func SetLiteral(tok *Token) {
	tok.Category = tcLiteral
}

func readNumber(path []byte, i int) (int, *Token, error) {
	e := skipNumber(path, i)
	f, err := strconv.ParseFloat(string(path[i:e]), 64)
	if err != nil {
		return i, nil, err
	}
	return e, &Token{Category: tcLiteral, Operand: Operand{Type: otNumber, Number: f}}, nil
}

func readString(path []byte, i int) (int, *Token, error) {
	e, err := skipString(path, i)
	if err != nil {
		return i, nil, err
	}
	return e, &Token{Category: tcLiteral, Operand: Operand{Type: otString, Str: path[i+1 : e-1]}}, nil
}

func readBoolNull(path []byte, i int) (int, *Token, error) {
	needles := [...][]byte{[]byte("false"), []byte("true"), []byte("null")}
	for n := 0; n < len(needles); n++ {
		if matchSubslice(path[i:], needles[n]) {
			if n == 2 {
				return i + len(needles[n]), &Token{Category: tcLiteral, Operand: Operand{Type: otNull}}, nil
			} else {
				return i + len(needles[n]), &Token{Category: tcLiteral, Operand: Operand{Type: otBoolean, Bool: n > 0}}, nil
			}
		}
	}
	return i, nil, errUnknownToken
}

func readRegexp(path []byte, i int) (int, *Token, error) {
	l := len(path)
	prev := byte(0)
	re := make([]byte, 0, 32)
	flags := make([]byte, 0, 8)
	i++
	for i < l && !(path[i] == '/' && prev != '\\') {
		prev = path[i]
		re = append(re, prev)
		i++
	}
	if i < l { // skip trailing '/'
		i++
	}
	flags = append(flags, '(', '?')
	for i < l && len(flags) < 4 && (path[i] == 'i' || path[i] == 'm' || path[i] == 's' || path[i] == 'U') {
		flags = append(flags, path[i])
		i++
	}
	flags = append(flags, ')')
	rex := ""
	if len(flags) > 3 {
		rex = string(flags) + string(re)
	} else {
		rex = string(re)
	}
	reg, err := regexp.Compile(rex)
	if err != nil {
		return i, nil, err
	}
	return i, &Token{Category: tcLiteral, Operand: Operand{Type: otRegexp, Regexp: reg}}, nil
}

func readVar(path []byte, i int) (int, *Token, error) {
	var err error
	l := len(path)
	s := i
	for i < l && !bytein(path[i], operatorBound) {
		if path[i] == '\'' || path[i] == '"' {
			i, err = skipString(path, i)
			if err != nil {
				return i, nil, err // unexpected EOL
			}
		} else {
			i++
		}
	}
	return i, &Token{Category: tcVariable, Operand: Operand{Type: otVariable, Str: path[s:i]}}, nil
}

func skipNumber(input []byte, i int) int {
	// numbers: -2  0.3  .3  1e2  -0.1e-2
	// [-][0[.[0]]][e[-]0]
	l := len(input)
	if i < l && input[i] == '-' {
		i++
	}
	for ; i < l && input[i] >= '0' && input[i] <= '9'; i++ {
	}
	for ; i < l && input[i] == '.'; i++ {
	}
	for ; i < l && input[i] >= '0' && input[i] <= '9'; i++ {
	}
	if i < l && (input[i] == 'E' || input[i] == 'e') {
		i++
	} else {
		return i
	}
	if i < l && (input[i] == '+' || input[i] == '-') {
		i++
	}
	for ; i < l && (input[i] >= '0' && input[i] <= '9'); i++ {
	}
	return i
}

func skipString(input []byte, i int) (int, error) {
	bound := input[i]
	done := false
	escaped := false
	s := i
	i++ // bound
	l := len(input)
	for i < l && !done {
		ch := input[i]
		if ch == bound && !escaped {
			done = true
		}
		escaped = ch == '\\' && !escaped
		i++
	}
	if i == l && !done {
		return s, errUnexpectedEndOfString
	}
	return i, nil
}

func matchSubslice(str, needle []byte) bool {
	l := len(needle)
	if len(str) < l {
		return false
	}
	for i := 0; i < l; i++ {
		if str[i] != needle[i] {
			return false
		}
	}
	return true
}

func getLastWord(path []byte) string {
	i := 0
	for ; i < len(path) && path[i] != ' ' && path[i] != '\t' && path[i] != '\r'; i++ {
	}
	return string(path[:i])
}
