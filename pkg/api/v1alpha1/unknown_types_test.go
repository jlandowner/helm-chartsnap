package v1alpha1

import (
	"reflect"
	"testing"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

func TestUnknown_Unstructured(t *testing.T) {
	raw := "some raw data"
	err := NewUnknownError(raw)

	obj := err.Unstructured()
	expectedObj := &metaV1.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "helm-chartsnap.jlandowner.dev/v1alpha1",
			"kind":       "Unknown",
			"metadata": map[string]interface{}{
				"name": "helm-output",
			},
			"raw": "some raw data",
		},
	}

	if !reflect.DeepEqual(obj, expectedObj) {
		t.Errorf("Expected obj to be %v, but got %v", err, obj)
	}

}

func TestUnknown_Error(t *testing.T) {
	type fields struct {
		Raw string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Unknown Error",
			fields: fields{
				Raw: "some raw data",
			},
			want: "failed to recognize a resource in stdout/stderr of helm template command output. snapshot it as Unknown: \n---\nsome raw data\n---",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Unknown{
				Raw: tt.fields.Raw,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Unknown.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnknown_Node(t *testing.T) {
	type fields struct {
		Raw string
	}
	tests := []struct {
		name   string
		fields fields
		want   *yaml.Node
	}{
		{
			name: "Test Unknown Node",
			fields: fields{
				Raw: `some raw
data
`,
			},
			want: &yaml.Node{
				Kind: yaml.MappingNode, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "apiVersion"},
					{Kind: yaml.ScalarNode, Value: "helm-chartsnap.jlandowner.dev/v1alpha1"},
					{Kind: yaml.ScalarNode, Value: "kind"},
					{Kind: yaml.ScalarNode, Value: "Unknown"},
					{Kind: yaml.ScalarNode, Value: "metadata"},
					{Kind: yaml.MappingNode, Content: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "name"},
						{Kind: yaml.ScalarNode, Value: "helm-output"},
					}},
					{Kind: yaml.ScalarNode, Value: "raw"},
					{Kind: yaml.ScalarNode, Value: "some raw\ndata\n", Style: yaml.LiteralStyle},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Unknown{
				Raw: tt.fields.Raw,
			}
			if got := e.Node(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unknown.Node() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnknown_MustString(t *testing.T) {
	type fields struct {
		Raw string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Unknown String",
			fields: fields{
				Raw: `some raw
data
`,
			},
			want: `apiVersion: helm-chartsnap.jlandowner.dev/v1alpha1
kind: Unknown
metadata:
  name: helm-output
raw: |
  some raw
  data
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Unknown{
				Raw: tt.fields.Raw,
			}
			got := e.MustString()
			if got != tt.want {
				t.Errorf("Unknown.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
