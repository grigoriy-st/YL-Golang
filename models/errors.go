package calculator

import "errors"

var (
	ErrInvalidExpression           = errors.New("Invalid expression")
	ErrDivisionByZero              = errors.New("Division by zero")
	ErrIncorrectSeqOfParenthese    = errors.New("An error in the sequence of parentheses")
	ErrDiffNumberOfBrackets        = errors.New("Different number of brackets")
	ErrConvertingNumberToFloatType = errors.New("Error converting a number to a float type")
	ErrTwoOperatorsInRow           = errors.New("Two operators in a row")
	ErrTwoOperandsInRow            = errors.New("Two operands in a row")
	ErrExpStartsWithOperator       = errors.New("The expression starts with the operator")
	ErrExpEndsWithOperator         = errors.New("The expression ends with the operator")
)
