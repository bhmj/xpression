package main

import "errors"

var (
	errUnexpectedEnd,
	errUnknownToken,
	errUnexpectedStringEnd,
	errMismatchedParentheses error
)

func init() {
	errUnexpectedEnd = errors.New("unexpected end of input")
	errUnknownToken = errors.New("unknown token")
	errUnexpectedStringEnd = errors.New("unexpected end of string")
	errMismatchedParentheses = errors.New("mismatched parentheses")
}
