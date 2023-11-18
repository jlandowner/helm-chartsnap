package snap

import (
	"io"
	"os"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aryann/difflib"
)

var _ = Describe("Snapshot", func() {
	f := func(m OmegaMatcher, filePath string) (success bool, err error) {
		f, err := os.Open(filePath)
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		buf, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())

		return m.Match(string(buf))
	}

	It("should match", func() {
		m := SnapShotMatcher("testdata/diff.snap", "default")
		success, err := f(m, "testdata/diff.txt")
		Expect(err).NotTo(HaveOccurred())
		Expect(success).To(BeTrue())
	})

	It("should not match and output diff N=1", func() {
		m := SnapShotMatcher("testdata/diff.snap", "default", WithDiffContextLineN(1))
		success, err := f(m, "testdata/diff_diff.txt")
		Expect(err).NotTo(HaveOccurred())
		Expect(success).To(BeFalse())

		Expect(m.FailureMessage(nil)).Should(Equal(`Expected to match
--- line=13
        --ca-file string                             verify certificates of HTTPS-enabled servers using this CA bundle
-       --cert-file string                           identify HTTPS client using this SSL certificate file
        --create-namespace                           create the release namespace if not present

--- line=23
    -g, --generate-name                              generate the name (and omit the NAME parameter)
-   -h, --help                                       help for template
+   -h, --help                                       help for templates
        --include-crds                               include CRDs in the templated output

--- line=56
    -f, --values strings                             specify values in a YAML file or a URL (can specify multiple)
+       --fake                                       fake
        --verify                                     verify the package before using it


`))
	})
})

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
