package charts

import (
	"testing"
)

func TestSnapshotID(t *testing.T) {
	type args struct {
		valuesFile string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basename of values file",
			args: args{
				valuesFile: "values.yaml",
			},
			want: "values",
		},
		{
			name: "basename of values file with path",
			args: args{
				valuesFile: "chart/test/values.yaml",
			},
			want: "values",
		},
		{
			name: "default if empty",
			args: args{
				valuesFile: "",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnapshotID(tt.args.valuesFile); got != tt.want {
				t.Errorf("SnapshotID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultSnapshotFilePath(t *testing.T) {
	type args struct {
		chartPath  string
		valuesFile string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "same directory with values file",
			args: args{
				chartPath:  "/tmp",
				valuesFile: "chart/test/values.yaml",
			},
			want: "chart/test/__snapshots__/values.snap",
		},
		{
			name: "chart directory with no values file and chart is local",
			args: args{
				chartPath:  "../../example/app1",
				valuesFile: "",
			},
			want: "../../example/app1/__snapshots__/default.snap",
		},
		{
			name: "chart directory with no values file and chart is remote",
			args: args{
				chartPath:  "ingress-nginx/ingress-nginx",
				valuesFile: "",
			},
			want: "__snapshots__/ingress-nginx/ingress-nginx/__snapshots__/default.snap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultSnapshotFilePath(tt.args.chartPath, tt.args.valuesFile); got != tt.want {
				t.Errorf("DefaultSnapshotFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnapshotFilePath(t *testing.T) {
	type args struct {
		dir        string
		valuesFile string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				dir:        "/tmp",
				valuesFile: "values.yaml",
			},
			want: "/tmp/__snapshots__/values.snap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnapshotFilePath(tt.args.dir, tt.args.valuesFile); got != tt.want {
				t.Errorf("SnapshotFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
