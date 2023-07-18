package snap

import (
	"reflect"
	"testing"

	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)

func TestMatchSnapShot(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want types.GomegaMatcher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchSnapShot(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchSnapShot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnapShotMatcher(t *testing.T) {
	type args struct {
		snapFile string
		snapId   string
	}
	tests := []struct {
		name string
		args args
		want *snapShotMatcher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnapShotMatcher(tt.args.snapFile, tt.args.snapId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SnapShotMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snapShotMatcher_Match(t *testing.T) {
	type fields struct {
		snapFilePath string
		snapId       string
		fs           afero.Fs
		expectedJson string
		actualJson   string
	}
	type args struct {
		actual interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantSuccess bool
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &snapShotMatcher{
				snapFilePath: tt.fields.snapFilePath,
				snapId:       tt.fields.snapId,
				fs:           tt.fields.fs,
				expectedJson: tt.fields.expectedJson,
				actualJson:   tt.fields.actualJson,
			}
			gotSuccess, err := m.Match(tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("snapShotMatcher.Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("snapShotMatcher.Match() = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

func Test_snapShotMatcher_FailureMessage(t *testing.T) {
	type fields struct {
		snapFilePath string
		snapId       string
		fs           afero.Fs
		expectedJson string
		actualJson   string
	}
	type args struct {
		actual interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMessage string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &snapShotMatcher{
				snapFilePath: tt.fields.snapFilePath,
				snapId:       tt.fields.snapId,
				fs:           tt.fields.fs,
				expectedJson: tt.fields.expectedJson,
				actualJson:   tt.fields.actualJson,
			}
			if gotMessage := m.FailureMessage(tt.args.actual); gotMessage != tt.wantMessage {
				t.Errorf("snapShotMatcher.FailureMessage() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}

func Test_snapShotMatcher_NegatedFailureMessage(t *testing.T) {
	type fields struct {
		snapFilePath string
		snapId       string
		fs           afero.Fs
		expectedJson string
		actualJson   string
	}
	type args struct {
		actual interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMessage string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &snapShotMatcher{
				snapFilePath: tt.fields.snapFilePath,
				snapId:       tt.fields.snapId,
				fs:           tt.fields.fs,
				expectedJson: tt.fields.expectedJson,
				actualJson:   tt.fields.actualJson,
			}
			if gotMessage := m.NegatedFailureMessage(tt.args.actual); gotMessage != tt.wantMessage {
				t.Errorf("snapShotMatcher.NegatedFailureMessage() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}

func Test_snapShotMatcher_ReadSnapShot(t *testing.T) {
	type fields struct {
		snapFilePath string
		snapId       string
		fs           afero.Fs
		expectedJson string
		actualJson   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &snapShotMatcher{
				snapFilePath: tt.fields.snapFilePath,
				snapId:       tt.fields.snapId,
				fs:           tt.fields.fs,
				expectedJson: tt.fields.expectedJson,
				actualJson:   tt.fields.actualJson,
			}
			got, err := m.ReadSnapShot()
			if (err != nil) != tt.wantErr {
				t.Errorf("snapShotMatcher.ReadSnapShot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("snapShotMatcher.ReadSnapShot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snapShotMatcher_WriteSnapShot(t *testing.T) {
	type fields struct {
		snapFilePath string
		snapId       string
		fs           afero.Fs
		expectedJson string
		actualJson   string
	}
	type args struct {
		snap []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &snapShotMatcher{
				snapFilePath: tt.fields.snapFilePath,
				snapId:       tt.fields.snapId,
				fs:           tt.fields.fs,
				expectedJson: tt.fields.expectedJson,
				actualJson:   tt.fields.actualJson,
			}
			if err := m.WriteSnapShot(tt.args.snap); (err != nil) != tt.wantErr {
				t.Errorf("snapShotMatcher.WriteSnapShot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_snapShotMatcher_readSnapFileData(t *testing.T) {
	type fields struct {
		snapFilePath string
		snapId       string
		fs           afero.Fs
		expectedJson string
		actualJson   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *map[string]data
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &snapShotMatcher{
				snapFilePath: tt.fields.snapFilePath,
				snapId:       tt.fields.snapId,
				fs:           tt.fields.fs,
				expectedJson: tt.fields.expectedJson,
				actualJson:   tt.fields.actualJson,
			}
			got, err := m.readSnapFileData()
			if (err != nil) != tt.wantErr {
				t.Errorf("snapShotMatcher.readSnapFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("snapShotMatcher.readSnapFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snapShotMatcher_writeSnapFileData(t *testing.T) {
	type fields struct {
		snapFilePath string
		snapId       string
		fs           afero.Fs
		expectedJson string
		actualJson   string
	}
	type args struct {
		snapFileData *map[string]data
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &snapShotMatcher{
				snapFilePath: tt.fields.snapFilePath,
				snapId:       tt.fields.snapId,
				fs:           tt.fields.fs,
				expectedJson: tt.fields.expectedJson,
				actualJson:   tt.fields.actualJson,
			}
			if err := m.writeSnapFileData(tt.args.snapFileData); (err != nil) != tt.wantErr {
				t.Errorf("snapShotMatcher.writeSnapFileData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
