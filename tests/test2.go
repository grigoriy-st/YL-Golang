package orchestrator

import (
	"testing"
)

func TestDistributeAndCompute(t *testing.T) {
	orch := NewOrchestrator([]string{"localhost:8081"})

	expr := "4 * 5 + 3"
	result, err := orch.ProcessExpression(expr)

	if err != nil {
		t.Fatalf("processing error: %v", err)
	}

	if result != 23 {
		t.Errorf("expected 23, got %v", result)
	}
}
