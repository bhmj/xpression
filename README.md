# Expression parser in Go

## What is it?

This project is a renewed version of expression parser used in [jsonslice](https://github.com/bhmj/jsonslice). It is still work in progress so use it with caution.

## Check it out

`go run ./... "expression"`

Expression examples:

`1+2`  
`2**1**2`  
`3 + 4 * 2 / (1-5) ** 2 ** 3`  
`5 + -5`  
`1/(3 & 5)`  
`'a' > 'b'`  
`'abc' =~ /a.c/i`  
`!((false))`

## Operators and data types supported

Operators | &nbsp;
--- | ---
Arithmetic | `+` `-` `*` `/` `**` `%`
Bitwise | `\|` `&` `^` `~` `<<` `>>`
Logical | `&&` `\|\|` `!`
Comparison | `==` `!=` `===` `!==` `>=` `>` `<=` `<`
Regexp | `=~` `!=~` `!~`
Parentheses | `(` `)`

<b>Data types</b> | &nbsp;
--- | ---
String constants | `'string'` or `"string"`
Numeric | 64-bit integers or floats
Boolean | `true` or `false`. Comparison results in boolean value.
Regexp | `/expression/` with modifiers:<br>`i` (case-insensitive), `m` (multiline), `s` (single-line), `U` (ungreedy)
Other | `null`

## Test coverage

Tests cover the majority of cases described in ECMAScript Language definition (specifically [ECMAScript Language: Expressions reference](https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html) and [Testing and Comparison Operations](https://tc39.es/ecma262/multipage/abstract-operations.html#sec-testing-and-comparison-operations)). 

## Benchmarks

Evaluate `(2) + (2) == (4)`

```diff
$ go test -bench=. -benchmem -benchtime=4s
goos: windows
goarch: amd64
pkg: github.com/bhmj/expression_parser
Benchmark_ModifiedNumericLiteral_WithParsing-4         1252026    3753 ns/op    1144 B/op   24 allocs/op
Benchmark_ModifiedNumericLiteral_WithoutParsing-4     32295195     158 ns/op       0 B/op    0 allocs/op
PASS
ok      github.com/bhmj/expression_parser       16.605
```

The same expression evaluated with [github.com/Knetic/govaluate](https://github.com/Knetic/govaluate) :

```diff
$ go test -bench='LiteralModifiers' -benchmem -benchtime=4s
goos: windows
goarch: amd64
pkg: github.com/Knetic/govaluate
BenchmarkEvaluationLiteralModifiers_WithParsing-4       559497    8499 ns/op    2272 B/op   49 allocs/op
BenchmarkEvaluationLiteralModifiers-4                 10908465     403 ns/op       8 B/op    1 allocs/op
PASS
ok      github.com/Knetic/govaluate     12.603s
```


## Changelog

**0.6.0** (2021-11-05) -- a remainder operator `%` added. Benchmarks added. Some optimization done.  
**0.5.0** (2021-11-04) -- Tests added. Multiple bugs fixed.  
**0.4.0** (2021-11-02) -- Expression evaluation.  
**0.3.0** (2021-11-01) -- MVP.

## Roadmap

- [x] arithmetic operators: `+ - * / **` `%`
- [x] bitwise operators: `| & ^ ~`
- [x] logical operators: `|| && !`
- [x] comparison operators: `> < >= <= == === != !==`
- [x] full support of parentheses
- [x] regular expressions for strings
- [x] unary minus supported
- [x] expression evaluation
- [x] parser test coverage
- [x] expression evaluation
- [x] evaluator test coverage
- [ ] refactor: operatorSpelling + operatorDetails -> map[spelling]{code, assoc, prec, args}. Store pointer to struct in Operand.
- [ ] add external reference type (node reference in jsonslice)
- [ ] Unicode support!

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :)

## Licence

[MIT](http://opensource.org/licenses/MIT)

## Author

Michael Gurov aka BHMJ
