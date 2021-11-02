package main

import "errors"

var (
	errUnrecognizedValue,
	errUnexpectedEnd,
	errUnknownToken,
	errInvalidCharacter,
	errUnexpectedStringEnd,
	errMismatchedParentheses error
)

func init() {
	errUnrecognizedValue = errors.New("unrecognized value: true, false or null expected")
	errUnexpectedEnd = errors.New("unexpected end of input")
	errUnknownToken = errors.New("unknown token")
	errInvalidCharacter = errors.New("invalid character")
	errUnexpectedStringEnd = errors.New("unexpected end of string")
	errMismatchedParentheses = errors.New("mismatched parentheses")
}
