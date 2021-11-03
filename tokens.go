package main

import (
	"regexp"
	"strconv"
)

func readNumber(path []byte, i int) (int, *Token, error) {
	e := skipNumber(path, i)
	f, err := strconv.ParseFloat(string(path[i:e]), 64)
	if err != nil {
		return i, nil, err
	}
	return e, &Token{Category: tcLiteral, Operand: Operand{Type: otNumber, Number: f}}, nil
}

func readString(path []byte, i int) (int, *Token, error) {
	bound := path[i]
	done := false
	escaped := false
	i++ // quote
	s := i
	l := len(path)
	for i < l && !done {
		ch := path[i]
		if ch == bound && !escaped {
			break
		}
		escaped = ch == '\\' && !escaped
		i++
	}
	if i == l && !done {
		return s, nil, errUnexpectedEndOfString
	}
	e := i
	i++ // unquote

	return i, &Token{Category: tcLiteral, Operand: Operand{Type: otString, Str: path[s:e]}}, nil
}

func readBool(path []byte, i int) (int, *Token, error) {
	needles := [...][]byte{[]byte("false"), []byte("true")}
	for n := 0; n < len(needles); n++ {
		if matchSubslice(path[i:], needles[n]) {
			return i + len(needles[n]), &Token{Category: tcLiteral, Operand: Operand{Type: otBoolean, Bool: n > 0}}, nil
		}
	}
	return i, nil, errUnknownToken
}

func readNull(path []byte, i int) (int, *Token, error) {
	if matchSubslice(path[i:], []byte("null")) {
		return i + 4, &Token{Category: tcLiteral, Operand: Operand{Type: otNull}}, nil
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

func skipNumber(input []byte, i int) int {
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
