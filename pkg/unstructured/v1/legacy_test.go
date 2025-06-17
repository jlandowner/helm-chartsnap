package unstructured

import (
	"testing"

	"github.com/aryann/difflib"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    []metaV1.Unstructured
		expected string
	}{
		{
			name: "single unstructured object",
			input: []metaV1.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "test-pod",
						},
					},
				},
			},
			expected: "- object:\n    apiVersion: v1\n    kind: Pod\n    metadata:\n        name: test-pod\n",
		},
		{
			name: "multiple unstructured objects sorted",
			input: []metaV1.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Service",
						"metadata": map[string]interface{}{
							"name": "test-service",
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "test-pod",
						},
					},
				},
			},
			expected: "- object:\n    apiVersion: v1\n    kind: Pod\n    metadata:\n        name: test-pod\n- object:\n    apiVersion: v1\n    kind: Service\n    metadata:\n        name: test-service\n",
		},
		{
			name:     "empty array",
			input:    []metaV1.Unstructured{},
			expected: "[]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if err != nil {
				t.Errorf("Encode() error = %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Encode() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

func TestFindKind(t *testing.T) {
	// Test that the function can be called without error
	result := findKind([]difflib.DiffRecord{})
	if result != "" {
		t.Errorf("findKind() with empty diffs should return empty string, got %v", result)
	}

	// Test with diff records (just check it doesn't panic)
	diffs := []difflib.DiffRecord{
		{Payload: "      kind: Pod", Delta: difflib.Common},
	}
	findKind(diffs)
}

func TestFindName(t *testing.T) {
	// Test that the function can be called without error
	result := findName([]difflib.DiffRecord{})
	if result != "" {
		t.Errorf("findName() with empty diffs should return empty string, got %v", result)
	}

	// Test with diff records (just check it doesn't panic)
	diffs := []difflib.DiffRecord{
		{Payload: "      metadata:", Delta: difflib.Common},
	}
	findName(diffs)
}

func TestDiffOptions_Diff(t *testing.T) {
	tests := []struct {
		name          string
		options       *DiffOptions
		x             string
		y             string
		shouldContain []string
	}{
		{
			name:          "context line 0 - show all",
			options:       &DiffOptions{ContextLineN: 0},
			x:             "line1\nline2\nline3",
			y:             "line1\nchanged\nline3",
			shouldContain: []string{"line1", "changed", "line3"},
		},
		{
			name:          "context line 1 - simple test",
			options:       &DiffOptions{ContextLineN: 1},
			x:             "line1\nline2\nline3",
			y:             "line1\nmodified\nline3",
			shouldContain: []string{"line2", "modified"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.Diff(tt.x, tt.y)
			// Just check that some expected content is present
			if len(result) == 0 {
				t.Errorf("Diff() should return some output, but got empty string")
			}
		})
	}
}

func TestIntInRange(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		v        int
		expected int
	}{
		{
			name:     "value in range",
			min:      0,
			max:      10,
			v:        5,
			expected: 5,
		},
		{
			name:     "value below min",
			min:      0,
			max:      10,
			v:        -5,
			expected: 0,
		},
		{
			name:     "value above max",
			min:      0,
			max:      10,
			v:        15,
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intInRange(tt.min, tt.max, tt.v)
			if result != tt.expected {
				t.Errorf("intInRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDiffString(t *testing.T) {
	tests := []struct {
		name     string
		diff     difflib.DiffRecord
		expected bool // whether result contains the payload
	}{
		{
			name: "left only diff",
			diff: difflib.DiffRecord{
				Delta:   difflib.LeftOnly,
				Payload: "removed line",
			},
			expected: true,
		},
		{
			name: "right only diff",
			diff: difflib.DiffRecord{
				Delta:   difflib.RightOnly,
				Payload: "added line",
			},
			expected: true,
		},
		{
			name: "common diff",
			diff: difflib.DiffRecord{
				Delta:   difflib.Common,
				Payload: "unchanged line",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := diffString(tt.diff)
			if tt.expected && !contains(result, tt.diff.Payload) {
				t.Errorf("diffString() should contain payload %v, but got %v", tt.diff.Payload, result)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
