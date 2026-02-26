package main

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/jlandowner/helm-chartsnap/pkg/api/v1alpha1"
	"github.com/jlandowner/helm-chartsnap/pkg/charts"
	. "github.com/jlandowner/helm-chartsnap/pkg/snap/gomega"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = Describe("overrideSnapshotConfig", func() {
	type testCase struct {
		name     string
		opt      *option
		cfg      *v1alpha1.SnapshotConfig
		expected *v1alpha1.SnapshotConfig
	}

	DescribeTable("overriding snapshot config",
		func(tc testCase) {
			tc.opt.overrideSnapshotConfig(tc.cfg)

			Expect(tc.cfg.DynamicFields).To(Equal(tc.expected.DynamicFields))
			Expect(tc.cfg.SnapshotFileExt).To(Equal(tc.expected.SnapshotFileExt))
			Expect(tc.cfg.SnapshotVersion).To(Equal(tc.expected.SnapshotVersion))
		},
		Entry("when SnapshotFileExt is set in option", testCase{
			name: "SnapshotFileExt is set",
			opt: &option{
				SnapshotFileExt: "yaml",
			},
			cfg: &v1alpha1.SnapshotConfig{},
			expected: &v1alpha1.SnapshotConfig{
				SnapshotFileExt: "yaml",
			},
		}),
		Entry("when SnapshotVersion is set in option", testCase{
			name: "SnapshotVersion is set",
			opt: &option{
				SnapshotVersion: "v2",
			},
			cfg: &v1alpha1.SnapshotConfig{},
			expected: &v1alpha1.SnapshotConfig{
				SnapshotVersion: "v2",
			},
		}),
		Entry("when SnapshotFileExt is set in config", testCase{
			name: "SnapshotFileExt is set",
			opt:  &option{},
			cfg: &v1alpha1.SnapshotConfig{
				SnapshotFileExt: "yaml",
			},
			expected: &v1alpha1.SnapshotConfig{
				SnapshotFileExt: "yaml",
			},
		}),
		Entry("when SnapshotVersion is set in config", testCase{
			name: "SnapshotVersion is set",
			opt:  &option{},
			cfg: &v1alpha1.SnapshotConfig{
				SnapshotVersion: "v2",
			},
			expected: &v1alpha1.SnapshotConfig{
				SnapshotVersion: "v2",
			},
		}),
		Entry("when LegacySnapshot is true", testCase{
			name: "LegacySnapshot is true",
			opt: &option{
				LegacySnapshot: true,
			},
			cfg: &v1alpha1.SnapshotConfig{},
			expected: &v1alpha1.SnapshotConfig{
				SnapshotVersion: charts.SnapshotVersionV1,
			},
		}),
		Entry("when both SnapshotFileExt and SnapshotVersion are set", testCase{
			name: "both SnapshotFileExt and SnapshotVersion are set",
			opt: &option{
				SnapshotFileExt: "yaml",
				SnapshotVersion: "v3",
			},
			cfg: &v1alpha1.SnapshotConfig{
				SnapshotFileExt: "yml",
				SnapshotVersion: "v2",
			},
			expected: &v1alpha1.SnapshotConfig{
				SnapshotFileExt: "yaml",
				SnapshotVersion: "v3",
			},
		}),
	)
})

var _ = Describe("option methods", func() {
	Describe("OK", func() {
		It("should return 'updated' when UpdateSnapshot is true", func() {
			opt := &option{UpdateSnapshot: true}
			Expect(opt.OK()).To(Equal("updated"))
		})

		It("should return 'matched' when UpdateSnapshot is false", func() {
			opt := &option{UpdateSnapshot: false}
			Expect(opt.OK()).To(Equal("matched"))
		})
	})

	Describe("HelmBin", func() {
		It("should return HELM_BIN env var when set", func() {
			os.Setenv("HELM_BIN", "/custom/helm")
			defer os.Unsetenv("HELM_BIN")
			opt := &option{}
			Expect(opt.HelmBin()).To(Equal("/custom/helm"))
		})

		It("should return 'helm' when HELM_BIN env var is not set", func() {
			os.Unsetenv("HELM_BIN")
			opt := &option{}
			Expect(opt.HelmBin()).To(Equal("helm"))
		})
	})

	Describe("Namespace", func() {
		It("should return HELM_NAMESPACE env var when set", func() {
			os.Setenv("HELM_NAMESPACE", "custom-namespace")
			defer os.Unsetenv("HELM_NAMESPACE")
			opt := &option{NamespaceFlag: "default"}
			Expect(opt.Namespace()).To(Equal("custom-namespace"))
		})

		It("should return NamespaceFlag when HELM_NAMESPACE env var is not set", func() {
			os.Unsetenv("HELM_NAMESPACE")
			opt := &option{NamespaceFlag: "test-namespace"}
			Expect(opt.Namespace()).To(Equal("test-namespace"))
		})
	})

	Describe("snapshotVersion", func() {
		It("should return v1 when LegacySnapshot is true", func() {
			opt := &option{LegacySnapshot: true, SnapshotVersion: "v3"}
			Expect(opt.snapshotVersion()).To(Equal(charts.SnapshotVersionV1))
		})

		It("should return SnapshotVersion when LegacySnapshot is false", func() {
			opt := &option{LegacySnapshot: false, SnapshotVersion: "v2"}
			Expect(opt.snapshotVersion()).To(Equal("v2"))
		})
	})
})

var _ = Describe("loadSnapshotConfig", func() {
	var tempDir string
	var configPath string
	var cfg *v1alpha1.SnapshotConfig

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "chartsnap-test")
		Expect(err).ShouldNot(HaveOccurred())
		configPath = tempDir + "/.chartsnap.yaml"
		cfg = &v1alpha1.SnapshotConfig{}
		o = &option{FailFast: false}
		// Initialize the logger to prevent panic
		log = slog.New(slogHandler())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
		log = nil
	})

	Context("when config file exists and is valid", func() {
		It("should load config successfully", func() {
			configContent := `snapshotVersion: v2
snapshotFileExt: yaml`
			err := os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ShouldNot(HaveOccurred())

			err = loadSnapshotConfig(configPath, cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfg.SnapshotVersion).To(Equal("v2"))
			Expect(cfg.SnapshotFileExt).To(Equal("yaml"))
		})
	})

	Context("when config file does not exist", func() {
		It("should not return error", func() {
			err := loadSnapshotConfig(configPath, cfg)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when config file exists but is invalid and FailFast is true", func() {
		It("should return error", func() {
			o.FailFast = true
			invalidContent := `invalid: yaml: content:`
			err := os.WriteFile(configPath, []byte(invalidContent), 0644)
			Expect(err).ShouldNot(HaveOccurred())

			err = loadSnapshotConfig(configPath, cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to load snapshot config"))
		})
	})

	Context("when config file exists but is invalid and FailFast is false", func() {
		It("should not return error but log warning", func() {
			o.FailFast = false
			invalidContent := `invalid: yaml: content:`
			err := os.WriteFile(configPath, []byte(invalidContent), 0644)
			Expect(err).ShouldNot(HaveOccurred())

			err = loadSnapshotConfig(configPath, cfg)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("rootCmd", func() {
	BeforeEach(func() {
		initRootCmd()
	})
	Context("success", func() {
		Context("snapshot local chart with single values file", func() {
			It("should pass", func() {
				var output bytes.Buffer
				color.Output = &output
				DeferCleanup(func() {
					color.Output = os.Stdout
				})
				rootCmd.SetArgs([]string{"--chart", "example/app1", "-f", "example/app1/test_latest/test_ingress_enabled.yaml", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
				Ω(output.String()).To(MatchSnapShot())
			})
		})

		Context("snapshot local chart with values directory", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "-f", "example/app1/test_latest/", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("snapshot OCI chart", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "oci://ghcr.io/nginxinc/charts/nginx-gateway-fabric", "-f", "example/remote/nginx-gateway-fabric.values.yaml", "--", "--namespace", "nginx-gateway"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("snapshot remote chart 1", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "cilium", "-f", "example/remote/cilium.values.yaml", "--", "--namespace", "kube-system", "--repo", "https://helm.cilium.io"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("snapshot remote chart 2", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "ingress-nginx", "-f", "example/remote/ingress-nginx.values.yaml", "--", "--namespace", "ingress-nginx", "--repo", "https://kubernetes.github.io/ingress-nginx", "--skip-tests"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("snapshot empty chart", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app2", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("error handling for specifying the empty directory", func() {
			It("should return an error for empty values", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "-f", "scripts", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
			})
		})

		Context("snapshot helm error", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app3", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("snapshot empty chart with no config file", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app2", "--namespace", "default", "--config-file", "notfound"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("snapshot empty chart with debug mode", func() {
			It("should pass", func() {
				os.Setenv("HELM_DEBUG", "true")
				defer os.Unsetenv("HELM_DEBUG")
				rootCmd.SetArgs([]string{"--chart", "example/app2", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("env FORCE_COLOR is enabled", func() {
			It("should force a colorized output", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "-f", "example/app1/test_latest/test_ingress_enabled.yaml", "--namespace", "default"})
				var output bytes.Buffer
				color.Output = &output
				os.Setenv("FORCE_COLOR", "1")
				DeferCleanup(func() {
					color.Output = os.Stdout
					color.NoColor = true
					os.Unsetenv("FORCE_COLOR")
				})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
				Ω(output.String()).To(MatchSnapShot())
			})
		})

		Context("snapshot file ext is yaml", func() {
			It("should pass", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app3", "--namespace", "default", "--snapshot-file-ext", "yaml"})
				err := rootCmd.Execute()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("fail", func() {
		Context("including dynamic outputs", func() {
			It("should fail", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
			})
		})

		Context("snapshot is different", func() {
			It("should fail with snapshotNotMatchError for exit code 2", func() {
				var output bytes.Buffer
				color.Output = &output
				DeferCleanup(func() {
					color.Output = os.Stdout
				})
				rootCmd.SetArgs([]string{"--chart", "example/app1", "--namespace", "default", "-f", "example/app1/testfail/test_ingress_enabled.yaml"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
				Ω(output.String()).To(MatchSnapShot())
				var notMatchErr *snapshotNotMatchError
				Expect(errors.As(err, &notMatchErr)).To(BeTrue())
			})
		})

		Context("values directory contains not matched snapshots", func() {
			It("should fail", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "--namespace", "default", "-f", "example/app1/testfail"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
			})
		})

		Context("values file not found", func() {
			It("should fail with general error for exit code 1", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "-f", "example/app1/test_latest/notfound.yaml", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
				var notMatchErr *snapshotNotMatchError
				Expect(errors.As(err, &notMatchErr)).To(BeFalse())
			})
		})

		Context("snapshot helm error with --fail-helm-error", func() {
			It("should fail", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app3", "--namespace", "default", "--fail-helm-error"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
			})
		})

		Context("required flag is not set", func() {
			It("should fail", func() {
				rootCmd.SetArgs([]string{})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
			})
		})

		Context("invalid flag", func() {
			It("should fail", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "-f", "example/app1/test_latest/test_ingress_enabled.yaml", "--namespace", "default", "--invalid"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
			})
		})
	})

	Context("--help", func() {
		It("should show help", func() {
			rootCmd.SetArgs([]string{"--help"})
			help := rootCmd.UsageString()
			Ω(help).To(MatchSnapShot())
		})
	})
})
