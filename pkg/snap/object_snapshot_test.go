package snap

import (
	"reflect"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestObjectSnapshot(t *testing.T) {
	type args struct {
		obj client.Object
	}
	tests := []struct {
		name string
		args args
		want client.Object
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ObjectSnapshot(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectSnapshot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDynamicFields(t *testing.T) {
	type args struct {
		o client.Object
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveDynamicFields(tt.args.o)
		})
	}
}
