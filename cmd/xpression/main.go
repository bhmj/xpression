package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bhmj/readline"
	"github.com/bhmj/xpression"
)

var (
	variableAssignment = regexp.MustCompile(`(^\w+)\s*=\s*([^=~].*)`)
	variables          map[string]xpression.Operand
	verbose            bool
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--help":
			fmt.Printf("Expression evaluator.\n")
			fmt.Printf("Usage: %[1]s [expression]\n  or just run '%[1]s' without arguments for interactive mode.\n", filepath.Base(os.Args[0]))
			return
		case "--verbose":
			verbose = true
		default:
			evaluateAndPrint(os.Args[1])
			return
		}
	}

	variables = make(map[string]xpression.Operand)

	fmt.Printf("Expression evaluator.\n")
	fmt.Printf("Run with --verbose parameter to print a NPN stack for each expression.\n")
	fmt.Printf("You can assign variables using 'var = expression' syntax.\n")
	fmt.Printf("Example:\n")
	fmt.Printf("  > pi = 3.1415926536 * 2\n")
	fmt.Printf("  > pi / 2\n")
	fmt.Printf("\nUse Up or Down to navigate through history\n\n")
	fmt.Printf("Enter expression or 'q' to quit\n")
	for {
		fmt.Print("> ")
		input, err := readline.Read()
		if err != nil {
			fmt.Printf("readline: %v\n", err)
			break
		}
		if input == "q" {
			break
		}
		if input != "" {
			if !parseVariable(input) {
				evaluateAndPrint(input)
			}
		}
	}
}

func evaluateAndPrint(str string) {
	if verbose {
		evalVerbose(str)
	} else {
		evalSilent(str)
	}
}

func evalVerbose(str string) {
	tokens, err := xpression.Parse([]byte(str))
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

func evalSilent(str string) {
	result, err := xpression.EvalVar([]byte(str), variableGetter)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(result.String())
	}
}

func parseVariable(str string) bool {
	matches := variableAssignment.FindStringSubmatch(str)
	if len(matches) == 0 {
		return false
	}
	varName := matches[1]
	varValue := matches[2]
	value, err := xpression.EvalVar([]byte(varValue), variableGetter)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}

	fmt.Println(value.String())

	variables[varName] = *value

	return true
}

func variableGetter(variable []byte, result *xpression.Operand) error {
	val, found := variables[string(variable)]
	if !found {
		return errors.New("variable " + string(variable) + " not found")
	}
	*result = val
	return nil
}

func printParsedExpression(tokens []*xpression.Token) {
	for _, tok := range tokens {
		fmt.Printf("%s ", tok.String())
	}
	fmt.Println()
}
