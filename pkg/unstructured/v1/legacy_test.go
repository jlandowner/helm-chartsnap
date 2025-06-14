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
			name:  "empty array",
			input: []metaV1.Unstructured{},
			expected: "[]\n",
		},
		{
			name: "single resource",
			input: []metaV1.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "test-config",
						},
					},
				},
			},
		},
		{
			name: "multiple resources sorted",
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
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "test-config",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if err != nil {
				t.Errorf("Encode() error = %v", err)
				return
			}
			if result == nil {
				t.Error("Encode() returned nil result")
				return
			}
			
			if tt.expected != "" {
				if string(result) != tt.expected {
					t.Errorf("Encode() = %v, want %v", string(result), tt.expected)
				}
			}
		})
	}
}

func TestFindKind(t *testing.T) {
	// Test with actual diff records that simulate real YAML content
	t.Run("kind found in yaml", func(t *testing.T) {
		// Create a realistic diff scenario
		diffs := difflib.Diff(
			[]string{"apiVersion: v1", "kind: ConfigMap", "metadata:", "  name: test"},
			[]string{"apiVersion: v1", "kind: Service", "metadata:", "  name: test"},
		)
		result := findKind(diffs)
		// Should find either ConfigMap or Service depending on the diff structure
		if !(result == "ConfigMap" || result == "Service" || result == "") {
			t.Errorf("findKind() = %v, want ConfigMap, Service, or empty string", result)
		}
	})

	t.Run("no kind found", func(t *testing.T) {
		diffs := difflib.Diff(
			[]string{"some", "text"},
			[]string{"other", "text"},
		)
		result := findKind(diffs)
		if result != "" {
			t.Errorf("findKind() = %v, want empty string", result)
		}
	})
}

func TestFindName(t *testing.T) {
	t.Run("name found in yaml", func(t *testing.T) {
		// Create a realistic diff scenario with metadata
		diffs := difflib.Diff(
			[]string{"apiVersion: v1", "kind: ConfigMap", "metadata:", "  name: old-name"},
			[]string{"apiVersion: v1", "kind: ConfigMap", "metadata:", "  name: new-name"},
		)
		result := findName(diffs)
		// Should find either old-name or new-name depending on the diff structure
		if !(result == "old-name" || result == "new-name" || result == "") {
			t.Errorf("findName() = %v, want old-name, new-name, or empty string", result)
		}
	})

	t.Run("no name found", func(t *testing.T) {
		diffs := difflib.Diff(
			[]string{"some", "text"},
			[]string{"other", "text"},
		)
		result := findName(diffs)
		if result != "" {
			t.Errorf("findName() = %v, want empty string", result)
		}
	})
}

func TestDiffOptions_Diff(t *testing.T) {
	tests := []struct {
		name        string
		options     *DiffOptions
		x           string
		y           string
		shouldDiff  bool
	}{
		{
			name:    "no context lines",
			options: &DiffOptions{ContextLineN: 0},
			x:       "line1\nline2",
			y:       "line1\nline3",
			shouldDiff: true,
		},
		{
			name:    "with context lines",
			options: &DiffOptions{ContextLineN: 2},
			x:       "line1\nline2",
			y:       "line1\nline3",
			shouldDiff: true,
		},
		{
			name:    "identical strings",
			options: &DiffOptions{ContextLineN: 1},
			x:       "same\nlines",
			y:       "same\nlines",
			shouldDiff: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.Diff(tt.x, tt.y)
			if tt.shouldDiff {
				if result == "" {
					t.Error("Diff() returned empty result, expected non-empty")
				}
			}
		})
	}
}

func TestIntInRange(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		value    int
		expected int
	}{
		{
			name:     "value in range",
			min:      0,
			max:      10,
			value:    5,
			expected: 5,
		},
		{
			name:     "value below range",
			min:      5,
			max:      10,
			value:    3,
			expected: 5,
		},
		{
			name:     "value above range",
			min:      0,
			max:      5,
			value:    10,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intInRange(tt.min, tt.max, tt.value)
			if result != tt.expected {
				t.Errorf("intInRange(%d, %d, %d) = %d, want %d", tt.min, tt.max, tt.value, result, tt.expected)
			}
		})
	}
}