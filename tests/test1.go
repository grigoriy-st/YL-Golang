package calculator

import (
	"testing"
)

func TestEvaluateSimpleExpression(t *testing.T) {
	expr := "2 + 3 * 4"
	expected := 14.0

	result, err := Evaluate(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestEvaluateInvalidExpression(t *testing.T) {
	_, err := Evaluate("2 + ")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}
