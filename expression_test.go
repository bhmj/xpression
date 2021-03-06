package xpression

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
		// variables
		{`@.foo`, `123`},
		{`@.foo.length() + 1`, `457`},
		// complex comparisons
		{`(123 == "123") == 123`, `false`},
		// hex numbers
		{`123 == 0x7b`, `true`},
		{`0x7B == 123`, `true`},
		// infinity comparison
		{`1/0 > 0`, `true`},
		{`1/0 >= 0`, `true`},
		{`1/0 > 1/0`, `false`},
		{`1/0 < 1/0`, `false`},
		{`1/0 >= 1/0`, `true`},
		{`1/0 <= 1/0`, `true`},
		{`1/0 == 1/0`, `true`},
		{`1/0 === 1/0`, `true`},
		{`-1/0 < 0`, `true`},
		{`-1/0 <= 0`, `true`},
		{`-1/0 > -1/0`, `false`},
		{`-1/0 < -1/0`, `false`},
		{`-1/0 >= -1/0`, `true`},
		{`-1/0 <= -1/0`, `true`},
		{`-1/0 == -1/0`, `true`},
		{`-1/0 === -1/0`, `true`},
	}

	varFunc := func(str []byte, result *Operand) error {
		if string(str) == "@.foo" {
			result.SetNumber(123)
			return nil
		}
		if string(str) == "@.foo.length()" {
			result.SetNumber(456)
			return nil
		}
		return errUnknownToken
	}

	l := len(tests)
	for i := 0; i < l; i++ {
		// same tests without spaces (except in strings)
		tests = append(tests, struct {
			Expression string
			Expected   string
		}{string(unspace([]byte(tests[i].Expression))), tests[i].Expected})
	}

	for _, tst := range tests {
		tokens, err := Parse([]byte(tst.Expression))
		if err != nil {
			t.Errorf(tst.Expression + " : " + err.Error())
		} else {
			operand, err := Evaluate(tokens, varFunc)
			if err != nil {
				t.Errorf(tst.Expression + " : " + err.Error())
			}
			result := operand.String()
			if result != tst.Expected {
				t.Errorf(tst.Expression + "\n\texpected `" + string(tst.Expected) + "`\n\tbut got  `" + result + "`")
			}
		}
	}
}

func Test_MultipleEvaluation(t *testing.T) {

	type probePair struct {
		Variable string
		Expected string
	}
	type testBundle struct {
		Expression string
		Probes     []probePair
	}
	testBundles := []testBundle{
		{
			`@.var == "1"`,
			[]probePair{
				{`null`, `false`},
				{`"1"`, `true`},
				{`"abc"`, `false`},
				{`1`, `true`},
			},
		},
		{
			`(@.var == "1") == false`,
			[]probePair{
				{`null`, `true`},
				{`"abc"`, `true`},
				{`"1"`, `false`},
			},
		},
	}

	for _, bundle := range testBundles {
		expression := bundle.Expression
		for _, tst := range bundle.Probes {
			varFunc := func(str []byte, result *Operand) error {
				switch string(tst.Variable) {
				case "null":
					result.SetNull()
					return nil
				case `"1"`:
					result.SetString("1")
					return nil
				case `"abc"`:
					result.SetString("abc")
					return nil
				case "1":
					result.SetNumber(1)
					return nil
				}
				return errUnknownToken
			}

			tokens, err := Parse([]byte(expression))
			if err != nil {
				t.Errorf(expression + " : " + err.Error())
			} else {
				operand, err := Evaluate(tokens, varFunc)
				if err != nil {
					t.Errorf(expression + " : " + err.Error())
				} else {
					token := &Token{Category: tcLiteral, Operand: *operand}
					result := token.String()

					if result != tst.Expected {
						t.Errorf(expression + "\n\texpected `" + string(tst.Expected) + "`\n\tbut got  `" + result + "`")
					}
				}
			}
		}
	}
}

func Test_VariableParser(t *testing.T) {

	type testPair struct {
		Expression string
		Expected   string
	}
	testPairs := []testPair{
		{`@.var`, `@.var`},
		{`@.var+1`, `@.var`},
		{`@.var[1]+1`, `@.var[1]`},
		{`@.var[$.i]+1`, `@.var[$.i]`},
		{`@.var[$.i+1]+1`, `@.var[$.i+1]`},
		{`@[1]+@[2]`, `@[1]`},
		{`@[-1][2]/4`, `@[-1][2]`},
	}

	for _, pair := range testPairs {
		_, token, err := readVar([]byte(pair.Expression), 0)
		if err != nil {
			t.Errorf(pair.Expression + ": expected `" + pair.Expected + "` but got error `" + err.Error() + "`")
		} else {
			if token.String() != pair.Expected {
				t.Errorf(pair.Expression + ": expected `" + pair.Expected + "` but got `" + token.String() + "`")
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
		{`"a" + "b`, errUnexpectedEndOfString.Error() + ` at 6: "b`},
		{`"a" # "b`, errUnknownToken.Error() + ` at 4: #`},
		{`"a" ( "b"`, errMismatchedParentheses.Error()},
		{`"a" =~ /a(b/`, "error parsing regexp: missing closing ): `a(b` at 12: "},
		{`?`, errUnknownToken.Error() + ` at 0: ?`},
		{`ABC`, errUnknownToken.Error()},
		{`0x123456789ABCDEF012345`, errTooLongHexadecimal.Error() + ` at 0: 0x123456789ABCDEF012345`},
		{`1 + @.'foo`, errUnexpectedEndOfString.Error() + ` at 6: 'foo`},
		{`1 + 0xABCDEFG`, errInvalidHexadecimal.Error() + ` at 4: 0xABCDEFG`},
	}

	for _, tst := range tests {
		gotError := false
		tokens, err := Parse([]byte(tst.Expression))
		if err != nil {
			gotError = true
			errMessage := err.Error()
			if errMessage != tst.Expected {
				t.Errorf(tst.Expression + "\n\texpected error `" + string(tst.Expected) + "`\n\tbut got `" + errMessage + "`")
			}
		} else {
			_, err := Evaluate(tokens, nil)
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

func Test_Bugfixes(t *testing.T) {

	tests := []struct {
		Expression string
		Expected   string
	}{
		// all vars are 123
		{`(0+a)`, "123"},              // closing bracket is a variable bound
		{`@.foo.length() + 1`, `124`}, // function call brackets are a part of the variable
		{`@.var(fn()) + 2`, `125`},    // nested function calls are allowed as a part of the variable
		{`"0x123" * 2`, `582`},        // string containing hex number is converted to number on evaluation
	}

	varFunc := func(str []byte, result *Operand) error {
		expectedVars := map[string]int{
			"a":              123,
			"@.foo.length()": 123,
			"@.var(fn())":    123,
		}
		v, found := expectedVars[string(str)]
		if !found {
			return errUnknownToken
		}
		result.SetNumber(float64(v))
		return nil
	}

	for _, tst := range tests {
		tokens, err := Parse([]byte(tst.Expression))
		if err != nil {
			t.Errorf(tst.Expression + " : " + err.Error())
		} else {
			operand, err := Evaluate(tokens, varFunc)
			if err != nil {
				t.Errorf(tst.Expression + " : " + err.Error())
			}
			result := operand.String()
			if result != tst.Expected {
				t.Errorf(tst.Expression + "\n\texpected `" + string(tst.Expected) + "`\n\tbut got  `" + result + "`")
			}
		}
	}
}

func Benchmark_ModifiedNumericLiteral_WithParsing(b *testing.B) {
	expression := `(2) + (2) == (4)`
	for i := 0; i < b.N; i++ {
		tokens, _ := Parse([]byte(expression))
		_, _ = Evaluate(tokens, nil)
	}
}

func Benchmark_ModifiedNumericLiteral_WithoutParsing(b *testing.B) {
	expression := `(2) + (2) == (4)`
	tokens, _ := Parse([]byte(expression))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens, nil)
	}
}

func unspace(buf []byte) []byte {
	var result []byte
	r := 0
	bound := byte(0)
	for r < len(buf) {
		if (buf[r] == '\'' || buf[r] == '"') && bound == 0 {
			bound = buf[r]
		} else if buf[r] == bound {
			bound = 0
		}
		if (buf[r] != ' ' && buf[r] != '\t') || bound > 0 {
			result = append(result, buf[r])
		}
		r++
	}
	return result
}
