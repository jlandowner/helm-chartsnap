package charts

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/jlandowner/helm-chartsnap/pkg/snap"
	"github.com/jlandowner/helm-chartsnap/pkg/unstructured"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var log *slog.Logger

func SetLogger(slogr *slog.Logger) {
	log = slogr
}

func Log() *slog.Logger {
	if log == nil {
		log = slog.Default()
	}
	return log
}

type ChartSnapOptions struct {
	HelmTemplateCmdOptions HelmTemplateCmdOptions
	SnapshotConfig         SnapshotConfig
	SnapshotFile           string
	DiffContextLineN       int
}

func Snap(ctx context.Context, o ChartSnapOptions) (match bool, failureMessage string, err error) {
	sv := SnapshotValues{}
	if o.HelmTemplateCmdOptions.ValuesFile != "" {
		f, err := os.Open(o.HelmTemplateCmdOptions.ValuesFile)
		if err != nil {
			return match, "", fmt.Errorf("failed to open values file: %w", err)
		}
		defer f.Close()

		err = yaml.NewDecoder(f).Decode(&sv)
		if err != nil {
			return match, "", fmt.Errorf("failed to decode values file: %w", err)
		}
	}

	// merge snapshot config file and config in snapshot values file
	sv.TestSpec.Merge(o.SnapshotConfig)

	Log().Debug("test spec from values file", "spec", sv.TestSpec)

	out, err := o.HelmTemplateCmdOptions.Execute(ctx)
	if err != nil {
		return match, "", fmt.Errorf("'helm template' command failed: %w: %s", err, out)
	}
	Log().Debug("helm template output", "output", string(out))

	manifests, decodeErrs := unstructured.Decode(string(out))
	if len(decodeErrs) > 0 {
		for _, err := range decodeErrs {
			Log().Info("loading helm output is done with warning")
			fmt.Println(err)
		}
	}

	if err := sv.TestSpec.ApplyFixedValue(manifests); err != nil {
		return match, "", fmt.Errorf("failed to replace json path: %w", err)
	}

	snap.SetLogger(Log())
	s := snap.UnstructuredSnapShotMatcher(
		o.SnapshotFile,
		SnapshotID(o.HelmTemplateCmdOptions.ValuesFile),
		snap.WithDiffContextLineN(o.DiffContextLineN))
	match, err = snap.UnstructuredMatch(s, manifests)

	if err != nil {
		return match, "", fmt.Errorf("failed to get snapshot: %w", err)
	}
	return match, s.FailureMessage(nil), nil
}

func SnapshotID(valuesFile string) string {
	if valuesFile != "" {
		return strings.ReplaceAll(path.Base(valuesFile), ".yaml", "")
	} else {
		return "default"
	}
}

func DefaultSnapshotFilePath(chartPath, valuesFile string) string {
	// if values file is specified, use the directory of the values file as the snapshot directory.
	// otherwise, use the chart directory.
	if valuesFile != "" {
		return SnapshotFilePath(path.Dir(valuesFile), valuesFile)
	} else {
		// if remote chart, create output directory
		if _, err := os.Stat(path.Join(chartPath, "Chart.yaml")); os.IsNotExist(err) {
			chartPath = path.Join("__snapshots__", path.Base(chartPath))
		}
		return SnapshotFilePath(chartPath, valuesFile)
	}
}

func SnapshotFilePath(dir, valuesFile string) string {
	return path.Join(dir, "__snapshots__", SnapshotID(valuesFile)+".snap")
}
