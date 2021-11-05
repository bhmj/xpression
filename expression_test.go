package main

import (
	"testing"
)

func Test_Expressions(t *testing.T) {

	tests := []struct {
		Expression string
		Expected   string
	}{
		// values
		{`1`, `1`},
		{` 1   `, `1`},
		{`-1`, `-1`},
		{`--1`, `1`},
		{`true`, `true`},
		{`/abc/`, `/abc/`},
		{`0.000`, `0`},
		{`"abc"`, `"abc"`},
		{`1.2E2`, `120`},
		{`1.2E+2`, `120`},
		{`1.2E-2`, `0.012`},
		// simple arithmentic
		{`1 + 2`, `3`},
		{`1 - 2`, `-1`},
		{`1--2`, `3`},    // differs from JS: we don't have "--" operator so double minus is treated as "minus negative 2"
		{`1 - - 2`, `3`}, // the same
		{`1 + - 2`, `-1`},
		{`3 * 2`, `6`},
		{`3 / 2`, `1.5`},
		{`1 / 0`, `+Inf`},
		{`-1 / 0`, `-Inf`},
		{`4 % 2`, `0`},
		{`5 % 2`, `1`},
		// bitwise operations
		{`6 | 3`, `7`},
		{`6 & 3`, `2`},
		{`6 ^ 3`, `5`},
		{`~6`, `-7`},
		// bitwise shifts
		{`6 << 2`, `24`},
		{`6 >> 2`, `1`},
		{`6 >> 3`, `0`},
		// operator precedence
		{`1 + 2 * 3`, `7`},
		{`1 + 2 * 3 + 4`, `11`},
		{`2 * 3 - 2`, `4`},
		{`2 - 3 / 4`, `1.25`},
		{`2 ** 3 ** 2`, `512`}, // exponentiation is right associated
		// parentheses
		{`(2)`, `2`},
		{`(1 + 2) * 3`, `9`},
		{`((1 + 2))`, `3`},
		// comparison: numbers
		{`1 == 2`, `false`},
		{`1 > 2`, `false`},
		{`1 < 2`, `true`},
		{`1 <= 2`, `true`},
		{`1 >= 2`, `false`},
		{`1 === 2`, `false`},
		{`1 !== 2`, `true`},
		// comparison: strings
		{`"" == ""`, `true`},
		{`"AA" == "BB"`, `false`},
		{`"AA" != "BB"`, `true`},
		{`"AA" > "BB"`, `false`},
		{`"AAA" >= "AA"`, `true`},
		{`"BB" <= "BBA"`, `true`},
		{`"AAA" < "BB"`, `true`},
		{`"AA" < "BB"`, `true`},
		{`"A" < "BB"`, `true`},
		{`"foobar" =~ /foo.*/`, `true`},
		{`"foobar" =~ /FOO.*/i`, `true`},
		{`"foobar" !=~ /FOO.*/`, `true`},
		// comparison: nulls
		{`null == 1`, `false`},
		{`null == "AA"`, `false`},
		{`null == true`, `false`},
		{`null == false`, `false`},
		{`null == null`, `true`},
		// comparison: different types
		{`"1" == 1`, `true`},
		{`"1" == true`, `true`},
		{`"" == false`, `true`},
		{`"0" == false`, `true`},
		{`"false" == false`, `false`},
		{`1234 =~ /1.*4/`, `true`},
		{`true =~ /TRU?E/i`, `true`},
		{`true > 0`, `true`},
		{`"xxx" > 4`, `false`},
		{`"xxx" > 4`, `false`},
		{`1 < /abc/`, `false`},
		// strict comparison
		{`"1" === 1`, `false`},
		{`null === 0`, `false`},
		{`"" === false`, `false`},
		{`"a" === "a"`, `true`},
		{`42 === 42`, `true`},
		{`null === null`, `true`},
		{`false === false`, `true`},
		// equality and non-equality
		{`1 == null`, `false`}, // note the difference between "equality" and "less than" / "greater than" operations
		{`0 == null`, `false`}, // when number is compared to null. This is due to different paragraphs describing
		{`1 > null`, `true`},   // equality and non-equality in ECMA reference [1]: 7.2.13 and 7.2.14
		{`0 > null`, `false`},  //
		{`1 >= null`, `true`},  //
		{`0 >= null`, `true`},  //
		{`1 <= null`, `false`}, //
		{`0 <= null`, `true`},  //
		// string concatenation
		{`"foo" + "bar"`, `"foobar"`},
		{`null + "able"`, `"nullable"`},
		{`false + " confessions"`, `"false confessions"`},
		// regexps
		{`"a" =~ "b"`, `false`},     // maybe undefined ?
		{`/foo/ =~ /bar/`, `false`}, // maybe undefined ?
		{`/123/ =~ 123`, `true`},
		// logical operations
		{`true && true`, `true`},
		{`true && false`, `false`},
		{`false && true`, `false`},
		{`false && false`, `false`},
		{`true || true`, `true`},
		{`true || false`, `true`},
		{`false || true`, `true`},
		{`false || false`, `false`},
		{`!true`, `false`},
		{`!false`, `true`},
		{`null || true`, `true`},
		{`null && true`, `null`},
		{`"a" || true`, `"a"`},
		{`"" || "b"`, `"b"`},
		{`"a" && true`, `true`},
		{`"a" && "b"`, `"b"`},
		{`1 || "a"`, `1`},
		{`0 || "a"`, `"a"`},
		{`0 || "a"`, `"a"`},
		{`/aa/ || "a"`, `"a"`},
	}

	for _, tst := range tests {
		// println(tst.Query)
		tokens, err := parseExpression([]byte(tst.Expression), 0)
		if err != nil {
			t.Errorf(tst.Expression + " : " + err.Error())
		} else {
			operand, _, err := evaluateExpression(tokens)
			if err != nil {
				t.Errorf(tst.Expression + " : " + err.Error())
			}
			token := &Token{Category: tcLiteral, Operand: *operand}
			result := token.String()

			if result != tst.Expected {
				t.Errorf(tst.Expression + "\n\texpected `" + string(tst.Expected) + "`\n\tbut got  `" + result + "`")
			}
		}
	}
}
func Test_Errors(t *testing.T) {

	tests := []struct {
		Expression string
		Expected   string
	}{
		// simple arithmentic
		{``, errNotEnoughArguments.Error()},
		{`1 + 2 +`, errNotEnoughArguments.Error()},
		{`2 * + 2`, errNotEnoughArguments.Error()},
		{`1..0 +`, `strconv.ParseFloat: parsing "1..0": invalid syntax at 0: 1..0`},
		{`"a" + "b`, errUnexpectedEndOfString.Error() + ` at 7: b`},
		{`"a" @ "b`, errUnknownToken.Error() + ` at 4: @`},
		{`"a" ( "b"`, errMismatchedParentheses.Error()},
		{`troo`, errUnknownToken.Error() + ` at 0: troo`},
		{`nool`, errUnknownToken.Error() + ` at 0: nool`},
		{`"a" =~ /a(b/`, "error parsing regexp: missing closing ): `a(b` at 12: "},
	}

	for _, tst := range tests {
		gotError := false
		tokens, err := parseExpression([]byte(tst.Expression), 0)
		if err != nil {
			gotError = true
			errMessage := err.Error()
			if errMessage != tst.Expected {
				t.Errorf(tst.Expression + "\n\texpected error `" + string(tst.Expected) + "`\n\tbut got `" + errMessage + "`")
			}
		} else {
			_, _, err := evaluateExpression(tokens)
			if err != nil {
				gotError = true
				errMessage := err.Error()
				if errMessage != tst.Expected {
					t.Errorf(tst.Expression + "\n\texpected error `" + string(tst.Expected) + "`\n\tbut got `" + errMessage + "`")
				}
			}
		}
		if !gotError {
			t.Errorf(tst.Expression + "\n\texpected error `" + string(tst.Expected) + "`\n\tbut got nothing")
		}
	}
}

func Benchmark_ModifiedNumericLiteral_WithParsing(b *testing.B) {
	expression := `(2) + (2) == (4)`
	for i := 0; i < b.N; i++ {
		tokens, _ := parseExpression([]byte(expression), 0)
		_, _, _ = evaluateExpression(tokens)
	}
}

func Benchmark_ModifiedNumericLiteral_WithoutParsing(b *testing.B) {
	expression := `(2) + (2) == (4)`
	tokens, _ := parseExpression([]byte(expression), 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = evaluateExpression(tokens)
	}
}