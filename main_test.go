package main

import (
	"bytes"
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
			It("should fail", func() {
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
			It("should fail", func() {
				rootCmd.SetArgs([]string{"--chart", "example/app1", "-f", "example/app1/test_latest/notfound.yaml", "--namespace", "default"})
				err := rootCmd.Execute()
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).To(MatchSnapShot())
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
