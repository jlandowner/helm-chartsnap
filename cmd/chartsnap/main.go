package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/cosmo-workspace/controller-testtools/pkg/charts"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

var (
	// goreleaser default https://goreleaser.com/customization/builds/
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
	o       = &option{}
	log     *slog.Logger
	values  []string
)

type option struct {
	HelmPath       string
	ReleaseName    string
	Namespace      string
	Chart          string
	ValuesFile     string
	Debug          bool
	UpdateSnapshot bool
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "chart-snapshot",
		Short: "helm chart snapshotter",
		Long: `
Snapshot test like Jest for Helm charts.

MIT 2023 cosmo-workspace/controller-testtools
`,
		Version: fmt.Sprintf("version=%s commit=%s date=%s", version, commit, date),
		RunE:    run,
	}
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().BoolVar(&o.Debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().BoolVarP(&o.UpdateSnapshot, "update-snapshot", "u", false, "update snapshot mode")
	rootCmd.PersistentFlags().StringVarP(&o.Chart, "chart", "c", "", "path to the chart directory")
	if err := rootCmd.MarkPersistentFlagDirname("chart"); err != nil {
		panic(err)
	}
	if err := rootCmd.MarkPersistentFlagRequired("chart"); err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().StringVar(&o.ReleaseName, "release-name", "testrelease", "release name")
	rootCmd.PersistentFlags().StringVar(&o.Namespace, "namespace", "testns", "namespace")
	rootCmd.PersistentFlags().StringVar(&o.HelmPath, "helm-path", "helm", "path to the helm command")
	rootCmd.PersistentFlags().StringVarP(&o.ValuesFile, "values", "f", "", "path to the test values file. this flag is passed to `helm template --values=values.yaml`")

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: func() slog.Leveler {
			if o.Debug {
				return slog.LevelDebug
			}
			return slog.LevelInfo
		}(),
	}))
	log.Debug("options", printOptions(*o)...)

	if o.ValuesFile == "" {
		values = []string{""}
	} else {
		if s, err := os.Stat(o.ValuesFile); os.IsNotExist(err) {
			return fmt.Errorf("values file '%s' not found", o.ValuesFile)
		} else if s.IsDir() {
			// get all values files in the directory
			files, err := os.ReadDir(o.ValuesFile)
			if err != nil {
				return fmt.Errorf("failed to read values file directory: %w", err)
			}
			values = make([]string, 0)
			for _, f := range files {
				// read only *.yaml
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".yaml") {
					values = append(values, path.Join(o.ValuesFile, f.Name()))
				}
			}
		} else {
			values = []string{o.ValuesFile}
		}
	}

	eg, ctx := errgroup.WithContext(cmd.Context())
	for _, v := range values {
		ht := charts.HelmTemplateCmdOptions{
			HelmPath:    o.HelmPath,
			ReleaseName: o.ReleaseName,
			Namespace:   o.Namespace,
			Chart:       o.Chart,
			ValuesFile:  v,
		}
		bannerPrintln("RUNS",
			fmt.Sprintf("Snapshot testing chart=%s values=%s", ht.Chart, ht.ValuesFile), 0, color.BgBlue)
		eg.Go(func() error {
			if o.UpdateSnapshot {
				if err := os.Remove(charts.SnapshotFile(ht.Chart, ht.ValuesFile)); err != nil {
					return fmt.Errorf("failed to replace snapshot file: %w", err)
				}
			}
			matched, failureMessage, err := charts.Snap(ctx, ht)
			if err != nil {
				bannerPrintln("FAIL", fmt.Sprintf("%v chart=%s values=%s", err, ht.Chart, ht.ValuesFile), color.FgRed, color.BgRed)
				return fmt.Errorf("failed to get snapshot: %w chart=%s values=%s", err, ht.Chart, ht.ValuesFile)
			}
			if !matched {
				bannerPrintln("FAIL", failureMessage, color.FgRed, color.BgRed)
				return fmt.Errorf("not match snapshot chart=%s values=%s", ht.Chart, ht.ValuesFile)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	bannerPrintln("PASS", "Snapshot matched", color.FgGreen, color.BgGreen)

	return nil
}

func bannerPrintln(banner string, message string, fgColor color.Attribute, bgColor color.Attribute) {
	color.New(color.FgWhite, bgColor).Printf(" %s ", banner)
	color.New(fgColor).Printf(" %s\n", message)
}

func printOptions(o option) []any {
	rv := reflect.ValueOf(o)
	rt := rv.Type()
	options := make([]any, rt.NumField()*2)

	for i := 0; i < rt.NumField(); i++ {
		options[i*2] = rt.Field(i).Name
		options[i*2+1] = rv.Field(i).Interface()
	}

	return options
}
