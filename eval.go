package xpression

// Eval evaluates expression and returns the result. No external variables used. See EvalVar for more.
func Eval(expression []byte) (*Operand, error) {
	tokens, err := Parse(expression)
	if err != nil {
		return nil, err
	}
	return Evaluate(tokens, nil)
}

// EvalVar evaluates expression and returns the result. External variables can be used via varFunc.
func EvalVar(expression []byte, varFunc VariableFunc) (*Operand, error) {
	tokens, err := Parse(expression)
	if err != nil {
		return nil, err
	}
	return Evaluate(tokens, varFunc)
}
