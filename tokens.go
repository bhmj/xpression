package main

import (
	"regexp"
	"strconv"
)

func readNumber(path []byte, i int) (int, *Token, error) {
	e := skipNumber(path, i)
	f, err := strconv.ParseFloat(string(path[i:e]), 64)
	if err != nil {
		return e, nil, err
	}
	return e, &Token{Category: tcLiteral, Type: ltNumber, Number: f}, nil
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
		return i, nil, errUnexpectedStringEnd
	}
	e := i
	i++ // unquote

	return i, &Token{Category: tcLiteral, Type: ltString, Str: path[s:e]}, nil
}

func readBool(path []byte, i int) (int, *Token, error) {
	s := i
	e, err := skipBoolNull(path, i)
	if err != nil {
		return i, nil, err
	}
	b, ok := map[string]bool{"true": true, "false": false}[string(path[s:e])]
	if !ok {
		return i, nil, errUnknownToken
	}
	return e, &Token{Category: tcLiteral, Type: ltBoolean, Bool: b}, nil
}

func readNull(path []byte, i int) (int, *Token, error) {
	s := i
	l := len(path)
	if l-i < 4 {
		return i, nil, errUnknownToken
	}
	null := []byte("null")
	diff := false
	for i-s < 4 && !diff {
		if path[i] != null[i-s] {
			diff = true
		}
		i++
	}
	if diff {
		return i, nil, errUnrecognizedValue
	}

	return i, &Token{Category: tcLiteral, Type: ltNull}, nil
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
	return i, &Token{Category: tcLiteral, Type: ltRegexp, Regexp: reg}, nil
}

func skipNumber(input []byte, i int) int {
	l := len(input)
	for ; i < l; i++ {
		ch := input[i]
		if !((ch >= '0' && ch <= '9') || ch == '.' || ch == '-' || ch == 'E' || ch == 'e') {
			break
		}
	}
	return i
}

func skipBoolNull(input []byte, i int) (int, error) {
	needles := [...][]byte{[]byte("true"), []byte("false"), []byte("null")}
	for n := 0; n < len(needles); n++ {
		if matchSubslice(input[i:], needles[n]) {
			return i + len(needles[n]), nil
		}
	}
	return i, errUnrecognizedValue
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
