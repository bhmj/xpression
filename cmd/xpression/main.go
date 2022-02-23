package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bhmj/xpression"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Expression evaluator.\nUsage: %[1]s <expression>\n", filepath.Base(os.Args[0]))
		return
	}

	tokens, err := xpression.Parse([]byte(os.Args[1]))

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		printParsedExpression(tokens)

		result, err := xpression.Evaluate(tokens, variableGetter)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println(result.String())
		}
	}
}

func variableGetter(variable []byte, result *xpression.Operand) error {
	if string(variable) == "@.foobar" {
		result.SetNumber(123)
		return nil
	}
	result.SetBoolean(false)
	return nil
}

func printParsedExpression(tokens []*xpression.Token) {
	for _, tok := range tokens {
		fmt.Printf("%s ", tok.String())
	}
	fmt.Println()
}
