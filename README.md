# Expression parser/evaluator in Go

## What is it?

This project is a renewed version of expression parser/evaluator used in [jsonslice](https://github.com/bhmj/jsonslice). It allows you to evaluate simple arithmetic expressions with variable support.

## Check it out

```
git clone https://github.com/bhmj/xpression.git

cd xpression

make build

./bin/xpression "1+2"
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
    result, _ := xpression.EvalStr(`5-3*(6-12)`)
    fmt.Println(result.String())

    // external data in expression (aka variables)
    foobar := 123
    varMapper := map[string]*int{
        `foobar`: &foobar,
    }
    varFunc := func(name []byte, result *xpression.Operand) error {
        ref := varMapper[string(name)]
        if ref == nil {
            return errors.New("unknown variable")
        }
        result.SetNumber(float64(*ref))
        return nil
    }
    result, _ = xpression.EvalVarStr(`27 / foobar`, varFunc)
    fmt.Println(result.String())
```
[Run in Go Playground](https://play.golang.com/p/YaweDHqUJKa)

## Functions

```Go
func Eval(expression []byte) (*Operand, error)
```
Eval evaluates expression and returns the result. No external variables used. See EvalVar for more.
Returns a pointer to Operand or error.
Use Operand's method `.String()` to get a serializable value or see Operand's `.OperandType` field to determine the result type and get the value from `.Str`, `.Number` or `.Bool` field.

```Go
func EvalVar(expression []byte, varFunc VariableFunc) (*Operand, error)
```
`EvalVar` evaluates expression and returns the result. External variables can be used via `varFunc`.

```Go
func EvalStr(expression string) (*Operand, error)
```
Same as `Eval` but receives a `string` instead of `[]byte`.

```Go
func EvalVarStr(expression string, varFunc VariableFunc) (*Operand, error)
```
Same as `EvalVar` but receives a `string` instead of `[]byte`.

## xpression CLI

You can find a simple and dumb expression evaluation CLI tool in cmd/xpression.  
Build it with `make build` and run with `./bin/xpression`.  

Feel free to make your own, better implementation.  

Remember that this evaluator complies with JavaScript rules, so the following example is legitimate and behaves JavaScript-like:  
```
> a = 123
123
> ("0x" + a) * 2
582
```

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

```golang
$ go test -bench=. -benchmem -benchtime=4s
goos: darwin
goarch: amd64
pkg: github.com/bhmj/xpression
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
Benchmark_ModifiedNumericLiteral_WithParsing-16       2622770    1811 ns/op   1272 B/op   26 allocs/op
Benchmark_ModifiedNumericLiteral_WithoutParsing-16   77455698   57.55 ns/op      0 B/op    0 allocs/op
PASS
ok      github.com/bhmj/xpression       11.548s
```

The same expression evaluated with [github.com/Knetic/govaluate](https://github.com/Knetic/govaluate) :

```golang
$ go test -bench='LiteralModifiers' -benchmem -benchtime=4s
goos: darwin
goarch: amd64
pkg: github.com/Knetic/govaluate
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkEvaluationLiteralModifiers_WithParsing-16    1000000    4019 ns/op   2208 B/op   43 allocs/op
BenchmarkEvaluationLiteralModifiers-16               30173640   147.2 ns/op      8 B/op    1 allocs/op
PASS
ok      github.com/Knetic/govaluate     9.810s
```


## Changelog

**0.9.4** (2022-11-28) -- Go 1.19 support. `EvalStr` and `EvalVarStr` helper functions added.  
**0.9.3** (2022-07-06) -- bugfixes: string to hex conversion `"0x123"*2`, closing bracket after a variable `(0+a)`.  
**0.9.2** (2022-02-23) -- Minor code refactoring. One-liner functions added (Eval, EvalVar).  
**0.9.1** (2022-01-02) -- Variable bounds refined.  
**0.9.0** (2021-11-19) -- Memory allocation reduced. Speed optimization.  
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
- [x] evaluator test coverage
- [x] add external reference type aka variable (node reference in jsonslice)
- [x] optimize memory allocations
- [ ] Unicode support
- [ ] add math functions support (?)
- [ ] add math/big support (?)

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
