package charts

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/jlandowner/helm-chartsnap/pkg/snap"
	unstV2 "github.com/jlandowner/helm-chartsnap/pkg/unstructured"
	unstV1 "github.com/jlandowner/helm-chartsnap/pkg/unstructured/v1"
)

const (
	SnapshotVersionV1 = "v1"
	SnapshotVersionV2 = "v2"
)

var logger *slog.Logger

func SetLogger(slogr *slog.Logger) {
	logger = slogr
}

func log() *slog.Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return logger
}

type ChartSnapshotter struct {
	HelmTemplateCmdOptions HelmTemplateCmdOptions
	SnapshotConfig         SnapshotConfig
	SnapshotFile           string
	SnapshotVersion        string
	DiffContextLineN       int
}

type SnapshotResult struct {
	Match          bool
	FailureMessage string
}

func (o *ChartSnapshotter) Snap(ctx context.Context) (result *SnapshotResult, err error) {
	// override snapshot config within values file's test spec
	sv := SnapshotValues{}
	if o.HelmTemplateCmdOptions.ValuesFile != "" {
		err = LoadSnapshotConfig(o.HelmTemplateCmdOptions.ValuesFile, &sv)
		if err != nil {
			return nil, fmt.Errorf("failed to decode values file: %w", err)
		}
	}
	log().Debug("loaded values", "values", sv, "path", o.SnapshotFile)
	sv.TestSpec.Merge(o.SnapshotConfig)

	// execute helm template command
	out, err := o.HelmTemplateCmdOptions.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("'helm template' command failed: %w: %s", err, out)
	}
	log().Debug("helm template output", "output", string(out), "path", o.SnapshotFile)

	// decode helm output
	manifests, decodeErrs := unstV2.Decode(string(out))
	if len(decodeErrs) > 0 {
		for _, err := range decodeErrs {
			log().Info("WARNING: loading helm output is done with warning")
			fmt.Println(err)
		}
	}

	// apply fixed values to dynamic fields
	if err := sv.TestSpec.ApplyFixedValue(manifests); err != nil {
		return nil, fmt.Errorf("failed to replace json path: %w", err)
	}

	// if snapshot file is v1 format, fallback to v1 snapshot matcher
	if snap.IsMultiSnapshots(o.SnapshotFile) {
		o.SnapshotVersion = SnapshotVersionV1
	}

	// take snapshot
	if o.SnapshotVersion == SnapshotVersionV1 {
		log().Info("WARNING: legacy format snapshot. it will be deprecated in the future version. please update the snapshots to the latest format", "path", o.SnapshotFile)
		return o.snapV1(manifests)
	} else {
		return o.snapV2(manifests)
	}
}

func (o *ChartSnapshotter) snapV1(manifests []metaV1.Unstructured) (result *SnapshotResult, err error) {
	snap.SetLogger(log())

	raw, err := unstV1.Encode(manifests)
	if err != nil {
		return nil, fmt.Errorf("failed to encode manifests: %w", err)
	}

	// v1 snapshot is multi snapshot format with encoding legacy formatted yaml
	matcher := snap.SnapshotMatcher(o.SnapshotFile,
		snap.WithSnapshotID(SnapshotFileName(o.HelmTemplateCmdOptions.ValuesFile)),
		snap.WithDiffFunc((&unstV1.DiffOptions{ContextLineN: o.DiffContextLineN}).Diff))

	match, err := matcher.Match(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}
	return &SnapshotResult{
		Match:          match,
		FailureMessage: matcher.FailureMessage(nil),
	}, nil
}

func (o *ChartSnapshotter) snapV2(manifests []metaV1.Unstructured) (result *SnapshotResult, err error) {
	snap.SetLogger(log())
	unstV2.SetLogger(log())

	raw, err := unstV2.Encode(manifests)
	if err != nil {
		return nil, fmt.Errorf("failed to encode manifests: %w", err)
	}
	matcher := snap.SnapshotMatcher(o.SnapshotFile, snap.WithDiffFunc((&unstV2.DiffOptions{ContextLineN: o.DiffContextLineN}).Diff))

	match, err := matcher.Match(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}
	return &SnapshotResult{
		Match:          match,
		FailureMessage: matcher.FailureMessage(nil),
	}, nil
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

func SnapshotFileName(valuesFile string) string {
	if valuesFile != "" {
		return strings.ReplaceAll(path.Base(valuesFile), ".yaml", "")
	} else {
		return "default"
	}
}

func SnapshotFilePath(dir, valuesFile string) string {
	return path.Join(dir, "__snapshots__", SnapshotFileName(valuesFile)+".snap")
}
