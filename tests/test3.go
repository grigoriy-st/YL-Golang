package agent

import (
	"testing"
)

func TestComputeTask(t *testing.T) {
	task := Task{
		Operator: "*",
		Left:     6,
		Right:    7,
	}

	result, err := Compute(task)
	if err != nil {
		t.Fatalf("compute error: %v", err)
	}

	if result != 42 {
		t.Errorf("expected 42, got %v", result)
	}
}
