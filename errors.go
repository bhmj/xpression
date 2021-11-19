package xpression

import "errors"

var (
	errUnknownToken,
	errUnexpectedEndOfString,
	errMismatchedParentheses,
	errNotEnoughArguments,
	errInvalidHexadecimal,
	errTooLongHexadecimal error
)

func init() {
	errUnknownToken = errors.New("unknown token")
	errUnexpectedEndOfString = errors.New("unexpected end of string")
	errMismatchedParentheses = errors.New("mismatched parentheses")
	errNotEnoughArguments = errors.New("not enough arguments")
	errInvalidHexadecimal = errors.New("invalid hexadecimal")
	errTooLongHexadecimal = errors.New("too long hexadecimal")
}
