package main

import "errors"

var (
	errUnknownToken,
	errUnexpectedEndOfString,
	errMismatchedParentheses,
	errNotEnoughArguments error
)

func init() {
	errUnknownToken = errors.New("unknown token")
	errUnexpectedEndOfString = errors.New("unexpected end of string")
	errMismatchedParentheses = errors.New("mismatched parentheses")
	errNotEnoughArguments = errors.New("not enough arguments")
}
