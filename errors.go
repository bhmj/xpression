package main

import "errors"

var (
	errUnexpectedEnd,
	errUnknownToken,
	errUnexpectedStringEnd,
	errMismatchedParentheses,
	errNotEnoughArguments error
)

func init() {
	errUnexpectedEnd = errors.New("unexpected end of input")
	errUnknownToken = errors.New("unknown token")
	errUnexpectedStringEnd = errors.New("unexpected end of string")
	errMismatchedParentheses = errors.New("mismatched parentheses")
	errNotEnoughArguments = errors.New("not enough arguments")
}
