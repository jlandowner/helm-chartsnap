package charts

import (
	"reflect"
	"testing"
)

func TestHelmTemplateCmdOptions_Args(t *testing.T) {
	type fields struct {
		HelmPath       string
		ReleaseName    string
		Namespace      string
		Chart          string
		ValuesFile     string
		AdditionalArgs []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "default",
			fields: fields{
				HelmPath:    "helm",
				ReleaseName: "chartsnap",
				Namespace:   "default",
				Chart:       "charts/app1/",
				ValuesFile:  "charts/app1/test/test.values.yaml",
			},
			want: []string{"template", "chartsnap", "charts/app1/", "--namespace=default", "--values=charts/app1/test/test.values.yaml"},
		},
		{
			name: "additional args",
			fields: fields{
				HelmPath:       "helm",
				ReleaseName:    "chartsnap",
				Namespace:      "xxx",
				Chart:          "postgres",
				ValuesFile:     "postgres.values.yaml",
				AdditionalArgs: []string{"--repo", "https://charts.bitnami.com/bitnami", "--skip-tests"},
			},
			want: []string{"template", "chartsnap", "postgres", "--namespace=xxx", "--values=postgres.values.yaml", "--repo", "https://charts.bitnami.com/bitnami", "--skip-tests"},
		},
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
			got := o.Args()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HelmTemplateCmdOptions.Args() = %v, want %v", got, tt.want)
			}
		})
	}
}
