package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Expression parser.\nUsage: %[1]s <expression>", filepath.Base(os.Args[0]))
		return
	}

	tokens, err := parseExpression([]byte(os.Args[1]), 0)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		printParsedExpression(tokens)

		result, _, err := evaluateExpression(tokens, variableGetter)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			printParsedExpression([]*Token{{Category: tcLiteral, Operand: *result}})
		}
	}
}

func printParsedExpression(tokens []*Token) {
	for _, tok := range tokens {
		fmt.Printf("%s ", tok.String())
	}
	fmt.Println()
}
