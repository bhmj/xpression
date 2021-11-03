# Expression parser in Go

## What is it?

This project is a renewed version of expression parser to use in jsonslice.

## Usage

`go run ./... "expression"`

Expression examples:

`1+2`  
`2**1**2`  
`3 + 4 * 2 / (1-5) ** 2 ** 3`  
`5 + -5`  
`1/(3 & 5)`  
`'a' > 'b'`  
`'abc' ~= /a.c/i`  
`!((false))`

## Test coverage

Tests cover the majority of cases described in ECMAScript Language definition (specifically [ECMAScript Language: Expressions reference](https://tc39.es/ecma262/multipage/ecmascript-language-expressions.html) and [Testing and Comparison Operations](https://tc39.es/ecma262/multipage/abstract-operations.html#sec-testing-and-comparison-operations)). 

## Changelog

**0.5.0** (2021-11-04) -- Tests added. Multiple bugs fixed.  
**0.4.0** (2021-11-02) -- Expression evaluation.  
**0.3.0** (2021-11-01) -- MVP.

## Roadmap

- [x] arithmetic operators: `+ - * / **`
- [x] bitwise operators: `| & ^ ~`
- [x] logical operators: `|| && !`
- [x] comparison operators: `> < >= <= == === != !==`
- [x] full support of parentheses
- [x] regular expressions for strings
- [x] unary minus supported
- [x] expression evaluation
- [+] parser test coverage
- [+] expression evaluation
- [+] evaluator test coverage
- [ ] add external reference type (node reference in jsonslice)

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
