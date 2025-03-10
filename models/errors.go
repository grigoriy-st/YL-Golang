package models

import "errors"

var (
	ErrDivisionByZero              = errors.New("division by zero")
	ErrInvalidExpression           = errors.New("invalid expression")
	ErrExpIsEmpty                  = errors.New("expressions is empty")
	ErrTwoOperandsInRow            = errors.New("two operands in a row")
	ErrTwoOperatorsInRow           = errors.New("two operators in a row")
	ErrDiffNumberOfBrackets        = errors.New("different number of brackets")
	ErrExpEndsWithOperator         = errors.New("the expression ends with the operator")
	ErrIncorrectSeqOfParenthese    = errors.New("an error in the sequence of parentheses")
	ErrExpStartsWithOperator       = errors.New("the expression starts with the operator")
	ErrConvertingNumberToFloatType = errors.New("error converting a number to a float type")
	ErrExpDoesNotMatchRegEx        = errors.New("the string does not match a regular expression")
)
