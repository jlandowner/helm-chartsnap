package charts

import (
	"context"
	"testing"
)

func TestSnap(t *testing.T) {
	type args struct {
		ctx context.Context
		o   HelmTemplateCmdOptions
	}
	tests := []struct {
		name               string
		args               args
		wantMatch          bool
		wantFailureMessage string
		wantErr            bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, gotFailureMessage, err := Snap(tt.args.ctx, tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("Snap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMatch != tt.wantMatch {
				t.Errorf("Snap() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
			if gotFailureMessage != tt.wantFailureMessage {
				t.Errorf("Snap() gotFailureMessage = %v, want %v", gotFailureMessage, tt.wantFailureMessage)
			}
		})
	}
}

func TestSnapshotID(t *testing.T) {
	type args struct {
		valuesFile string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnapshotID(tt.args.valuesFile); got != tt.want {
				t.Errorf("SnapshotID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnapshotFile(t *testing.T) {
	type args struct {
		chartPath  string
		valuesFile string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnapshotFile(tt.args.chartPath, tt.args.valuesFile); got != tt.want {
				t.Errorf("SnapshotFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
