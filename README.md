# Expression parser/evaluator in Go

## What is it?

This project is a renewed version of expression parser/evaluator used in [jsonslice](https://github.com/bhmj/jsonslice). It is still work in progress so use it with caution.

## Check it out

```
git clone https://github.com/bhmj/xpression.git

cd xpression

make build

./build/xpression "1+2"
````

Expression examples:

`1+2`  
`2**1**2`  
`3 + 4 * 2 / (1-5) ** 2 ** 3`  
`5 + -5`  
`1/(3 & 5)`  
`'a' > 'b'`  
`'abc' =~ /a.c/i`  
`!((false))`

## Usage

```Go
    // simple expression evaluation (error handling skipped)
    tokens, err := xpression.Parse([]byte(`5 - 3 * (6-12)`))
    result, err := xpression.Evaluate(tokens, nil)
    switch result.Type {
    case xpression.NumberOperand:
        fmt.Println(result.Number)
    default:
        fmt.Println("unexpected result")
    }

    // external data in expression (aka variables)
    foobar := 123
    varFunc := func(name []byte) (*xpression.Operand, error) {
        mapper := map[string]*int{
            `foobar`: &foobar
        }
        return xpression.Number(float64(*mapper[string(name)])), nil
    }
    tokens, err := xpression.Parse([]byte(`27 / foobar`))
    result, err := xpression.Evaluate(tokens, varFunc)
    fmt.Println(result.Number)
```
[Run in Go Playground](https://play.golang.com/p/xjS5Tj1_34b)

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
Numeric | 64-bit integers or floats in decimal or hexadecimal form: `123` or `0.123` or `1.2e34` or `0x12a` or `0x12A`
Boolean | `true` or `false`. Comparison results in boolean value.
Regexp | `/expression/` with modifiers:<br>`i` (case-insensitive), `m` (multiline), `s` (single-line), `U` (ungreedy)
Other | `null`

## Test coverage

Tests cover the majority of cases described in ECMAScript Language definition (specifically [ECMAScript Language: Expressions reference](https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html) and [Testing and Comparison Operations](https://tc39.es/ecma262/multipage/abstract-operations.html#sec-testing-and-comparison-operations)). 

## Benchmarks

Evaluate `(2) + (2) == (4)`

```diff
$ go test -bench=. -benchmem -benchtime=4s
goos: darwin
goarch: amd64
pkg: github.com/bhmj/xpression
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
Benchmark_ModifiedNumericLiteral_WithParsing-16        2714204    1853 ns/op    1272 B/op   26 allocs/op
Benchmark_ModifiedNumericLiteral_WithoutParsing-16    31129712   143.9 ns/op     128 B/op    2 allocs/op
PASS
ok      github.com/bhmj/xpression       12.363s
```

The same expression evaluated with [github.com/Knetic/govaluate](https://github.com/Knetic/govaluate) :

```diff
$ go test -bench='LiteralModifiers' -benchmem -benchtime=4s
goos: darwin
goarch: amd64
pkg: github.com/Knetic/govaluate
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkEvaluationLiteralModifiers_WithParsing-16     1000000    4019 ns/op    2208 B/op    43 allocs/op
BenchmarkEvaluationLiteralModifiers-16                30173640   147.2 ns/op       8 B/op     1 allocs/op
PASS
ok      github.com/Knetic/govaluate     9.810s
```


## Changelog

**0.8.0** (2021-11-11) -- hex numbers support. Production ready.  
**0.7.x** (2021-11-11) -- WIP  
**0.7.0** (2021-11-10) -- project renamed to `xpression`  
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
- [x] add external reference type (node reference in jsonslice)
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
