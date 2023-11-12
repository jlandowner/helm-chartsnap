package snap

import (
	"reflect"
	"testing"

	"github.com/aryann/difflib"
)

func TestDiffOptions_ContextLineN(t *testing.T) {
	type fields struct {
		DiffContextLineN int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "zero",
			fields: fields{DiffContextLineN: -1},
			want:   0,
		},
		{
			name:   "default",
			fields: fields{DiffContextLineN: 3},
			want:   3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &DiffOptions{
				DiffContextLineN: tt.fields.DiffContextLineN,
			}
			if got := o.ContextLineN(); got != tt.want {
				t.Errorf("DiffOptions.ContextLineN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeDiffOpts(t *testing.T) {
	type args struct {
		opts []DiffOptions
	}
	tests := []struct {
		name string
		args args
		want DiffOptions
	}{
		{
			name: "single: 0",
			args: args{opts: []DiffOptions{{DiffContextLineN: 0}}},
			want: DiffOptions{DiffContextLineN: 0},
		},
		{
			name: "single: 3",
			args: args{opts: []DiffOptions{{DiffContextLineN: 3}}},
			want: DiffOptions{DiffContextLineN: 3},
		},
		{
			name: "multiple: max 1",
			args: args{
				opts: []DiffOptions{
					{DiffContextLineN: 3},
					{DiffContextLineN: 2},
				},
			},
			want: DiffOptions{DiffContextLineN: 3},
		},
		{
			name: "multiple: max 2",
			args: args{
				opts: []DiffOptions{
					{DiffContextLineN: 2},
					{DiffContextLineN: 3},
				},
			},
			want: DiffOptions{DiffContextLineN: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeDiffOpts(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeDiffOpts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	type args struct {
		x string
		y string
		o DiffOptions
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single: 0",
			args: args{
				x: "a",
				y: "b",
				o: DiffOptions{DiffContextLineN: 0},
			},
			want: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.x, tt.args.y, tt.args.o); got != tt.want {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intInRange(t *testing.T) {
	type args struct {
		min int
		max int
		v   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "min",
			args: args{min: 0, max: 10, v: -1},
			want: 0,
		},
		{
			name: "max",
			args: args{min: 0, max: 10, v: 11},
			want: 10,
		},
		{
			name: "in range",
			args: args{min: 0, max: 10, v: 5},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intInRange(tt.args.min, tt.args.max, tt.args.v); got != tt.want {
				t.Errorf("intInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_diffString(t *testing.T) {
	type args struct {
		d difflib.DiffRecord
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "+",
			args: args{
				difflib.DiffRecord{
					Payload: "###",
					Delta:   difflib.RightOnly,
				},
			},
			want: "+ ###\n",
		},
		{
			name: "-",
			args: args{
				difflib.DiffRecord{
					Payload: "###",
					Delta:   difflib.LeftOnly,
				},
			},
			want: "- ###\n",
		},
		{
			name: "eq",
			args: args{
				difflib.DiffRecord{
					Payload: "###",
					Delta:   difflib.Common,
				},
			},
			want: "  ###\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := diffString(tt.args.d); got != tt.want {
				t.Errorf("diffString() = %v, want %v", got, tt.want)
			}
		})
	}
}
