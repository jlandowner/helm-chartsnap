package charts

import (
	"context"
	"reflect"
	"testing"
)

func TestHelmTemplateCmdOptions_Execute(t *testing.T) {
	type fields struct {
		HelmPath       string
		ReleaseName    string
		Namespace      string
		Chart          string
		ValuesFile     string
		AdditionalArgs []string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &HelmTemplateCmdOptions{
				HelmPath:       tt.fields.HelmPath,
				ReleaseName:    tt.fields.ReleaseName,
				Namespace:      tt.fields.Namespace,
				Chart:          tt.fields.Chart,
				ValuesFile:     tt.fields.ValuesFile,
				AdditionalArgs: tt.fields.AdditionalArgs,
			}
			got, err := o.Execute(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("HelmTemplateCmdOptions.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HelmTemplateCmdOptions.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
