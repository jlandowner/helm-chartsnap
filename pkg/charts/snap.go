package charts

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/jlandowner/helm-chartsnap/pkg/api/v1alpha1"
	"github.com/jlandowner/helm-chartsnap/pkg/snap"
	unstV2 "github.com/jlandowner/helm-chartsnap/pkg/unstructured"
	unstV1 "github.com/jlandowner/helm-chartsnap/pkg/unstructured/v1"
	"github.com/jlandowner/helm-chartsnap/pkg/yaml"
)

const (
	SnapshotVersionV1     = "v1"
	SnapshotVersionV2     = "v2"
	SnapshotVersionV3     = "v3"
	SnapshotVersionLatest = SnapshotVersionV3
)

var (
	logger *slog.Logger
	mutex  sync.Mutex
)

func SetLogger(slogr *slog.Logger) {
	mutex.Lock()
	defer mutex.Unlock()
	logger = slogr
}

func log() *slog.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	if logger == nil {
		logger = slog.Default()
	}
	return logger
}

type ChartSnapshotter struct {
	HelmTemplateCmdOptions HelmTemplateCmdOptions
	SnapshotConfig         v1alpha1.SnapshotConfig
	SnapshotDir            string
	DiffContextLineN       int
	UpdateSnapshot         bool
	HeaderVersion          string
	FailHelmError          bool

	snapshotFile string
}

type SnapshotResult struct {
	Match          bool
	FailureMessage string
}

func (o *ChartSnapshotter) prependSnapshotHeader(data []byte) []byte {
	header := (&v1alpha1.Header{SnapshotVersion: o.SnapshotConfig.SnapshotVersion}).ToString()
	data = append([]byte(header), data...)
	return data
}

func (o *ChartSnapshotter) getVersionFromSnapshotFile() string {
	s, err := snap.ReadFile(o.SnapshotFile())
	if err != nil {
		log().Debug("failed to read snapshot file", "err", err)
		return SnapshotVersionLatest
	}
	split := strings.Split(string(s), "\n")
	return v1alpha1.ParseHeader(split[0]).SnapshotVersion
}

func (o *ChartSnapshotter) snapshotFilePath() string {
	if o.SnapshotDir != "" {
		return snapshotFilePath(o.SnapshotDir, o.HelmTemplateCmdOptions.ValuesFile)
	} else {
		return defaultSnapshotFilePath(o.HelmTemplateCmdOptions.Chart, o.HelmTemplateCmdOptions.ValuesFile)
	}
}

func (o *ChartSnapshotter) SnapshotFile() string {
	if o.snapshotFile != "" {
		return o.snapshotFile
	}

	o.snapshotFile = o.snapshotFilePath()

	if o.SnapshotConfig.SnapshotFileExt != "" {
		// append snapshot file extension
		o.snapshotFile += "." + o.SnapshotConfig.SnapshotFileExt
	}

	if _, err := os.Stat(o.snapshotFile); err == nil {
		log().Debug("snapshot file already exists")
	} else if os.IsNotExist(err) {
		log().Debug("snapshot file does not exist")
	} else {
		log().Error("unexpected error in snapshot file stat", "err", err)
	}

	return o.snapshotFile
}

func (o *ChartSnapshotter) Snap(ctx context.Context) (result *SnapshotResult, err error) {
	// override snapshot config within values file's test spec
	sv := v1alpha1.SnapshotValues{}
	if o.HelmTemplateCmdOptions.ValuesFile != "" {
		err = v1alpha1.FromFile(o.HelmTemplateCmdOptions.ValuesFile, &sv)
		if err != nil {
			return nil, fmt.Errorf("failed to decode values file: %w", err)
		}
	}
	log().Debug("loaded values", "values", sv)
	sv.TestSpec.Merge(o.SnapshotConfig)
	o.SnapshotConfig = sv.TestSpec
	log().Debug("compluted config", "config", o.SnapshotConfig)

	SetLogger(log().With("path", o.SnapshotFile()))

	// execute helm template command
	out, err := o.HelmTemplateCmdOptions.Execute(ctx)
	if err != nil {
		if o.FailHelmError {
			return nil, fmt.Errorf("'helm template' command failed: %w: %s", err, out)
		} else {
			log().Error("helm command failed but snapshot it anyway. use --fail-helm-error if you want error exit code", "err", err, "output", string(out))
		}
	}
	log().Debug("helm template output", "output", string(out))

	// fallback if version is not specified
	if o.SnapshotConfig.SnapshotVersion == "" {
		if snap.IsMultiSnapshots(o.SnapshotFile()) {
			// v1: if snapshot file is v1(multi-snapshot, toml) format, fallback to v1 snapshot matcher
			o.SnapshotConfig.SnapshotVersion = SnapshotVersionV1
		} else if snapVersion := o.getVersionFromSnapshotFile(); snapVersion == "" {
			// v2: if snapshot file have no header, fallback to v2 snapshot matcher
			o.SnapshotConfig.SnapshotVersion = SnapshotVersionV2
		} else {
			// later: use snapshot version from snapshot file if exists
			o.SnapshotConfig.SnapshotVersion = snapVersion
		}
	}

	if o.UpdateSnapshot {
		log().Debug("updating snapshot")
		err := snap.RemoveFile(o.SnapshotFile())
		if err != nil && !os.IsNotExist(err) {
			log().Warn("failed to update snapshot", "err", err)
		}

		if o.SnapshotConfig.SnapshotFileExt != "" {
			// remove snapshot file without extension
			noExtSnapshotFile := o.snapshotFilePath()
			err := snap.RemoveFile(noExtSnapshotFile)
			if err != nil && !os.IsNotExist(err) {
				log().Warn("failed to remove snapshot file without specified extension", "err", err, "file", noExtSnapshotFile)
			}
		}
	}

	// take snapshot
	log().Debug("taking snapshot", "version", o.SnapshotConfig.SnapshotVersion)

	switch o.SnapshotConfig.SnapshotVersion {
	case SnapshotVersionV1:
		log().Warn("legacy format snapshot. it will be deprecated in the future version. please update the snapshots to the latest format")
		return o.snapV1(sv.TestSpec, out)
	case SnapshotVersionV2:
		return o.snapV2(sv.TestSpec, out)
	case SnapshotVersionV3:
		return o.snapV3(sv.TestSpec, out)
	default:
		log().Error("unsupported snapshot version. use latest", "version", o.SnapshotConfig.SnapshotVersion, "latest", SnapshotVersionLatest)
		o.SnapshotConfig.SnapshotVersion = SnapshotVersionLatest
		return o.snapV3(sv.TestSpec, out)
	}
}

func (o *ChartSnapshotter) snapV1(cfg v1alpha1.SnapshotConfig, data []byte) (result *SnapshotResult, err error) {
	// decode helm output to unstructured
	manifests, decodeErrs := unstV2.Decode(string(data))
	if len(decodeErrs) > 0 {
		for _, err := range decodeErrs {
			log().Warn("loading helm output is done with error")
			fmt.Println(err)
		}
	}

	// apply fixed values to dynamic fields
	if err := unstV2.ApplyFixedValue(cfg, manifests); err != nil {
		return nil, fmt.Errorf("failed to replace json path: %w", err)
	}

	snap.SetLogger(log())

	raw, err := unstV1.Encode(manifests)
	if err != nil {
		return nil, fmt.Errorf("failed to encode manifests: %w", err)
	}

	// v1 snapshot is multi snapshot format with encoding legacy formatted yaml
	matcher := snap.SnapshotMatcher(o.SnapshotFile(),
		snap.WithSnapshotID(snapshotFileName(o.HelmTemplateCmdOptions.ValuesFile)),
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

func (o *ChartSnapshotter) snapV2(cfg v1alpha1.SnapshotConfig, data []byte) (result *SnapshotResult, err error) {
	// decode helm output to unstructured
	manifests, decodeErrs := unstV2.Decode(string(data))
	if len(decodeErrs) > 0 {
		for _, err := range decodeErrs {
			log().Warn("loading helm output is done with error")
			fmt.Println(err)
		}
	}

	// apply fixed values to dynamic fields
	if err := unstV2.ApplyFixedValue(cfg, manifests); err != nil {
		return nil, fmt.Errorf("failed to replace json path: %w", err)
	}

	snap.SetLogger(log())
	unstV2.SetLogger(log())

	raw, err := unstV2.Encode(manifests)
	if err != nil {
		return nil, fmt.Errorf("failed to encode manifests: %w", err)
	}
	matcher := snap.SnapshotMatcher(o.SnapshotFile(), snap.WithDiffFunc((&unstV2.DiffOptions{ContextLineN: o.DiffContextLineN}).Diff))

	match, err := matcher.Match(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}
	return &SnapshotResult{
		Match:          match,
		FailureMessage: matcher.FailureMessage(nil),
	}, nil
}

func (o *ChartSnapshotter) snapV3(cfg v1alpha1.SnapshotConfig, data []byte) (result *SnapshotResult, err error) {
	snap.SetLogger(log())
	yaml.SetLogger(log())

	// decode helm output to kustomize/kyaml Nodes
	manifests, err := yaml.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode manifests: %w", err)
	}

	// apply fixed values to dynamic fields
	if err := yaml.ApplyFixedValueToDynamicFieleds(cfg, manifests); err != nil {
		return nil, fmt.Errorf("failed to replace json path: %w", err)
	}

	raw, err := yaml.Encode(manifests)
	if err != nil {
		return nil, fmt.Errorf("failed to encode manifests: %w", err)
	}

	matcher := snap.SnapshotMatcher(o.SnapshotFile(), snap.WithDiffFunc((&yaml.DiffOptions{ContextLineN: o.DiffContextLineN}).Diff))

	// add snapshot header on the top of the snapshot file from v3
	raw = o.prependSnapshotHeader(raw)

	match, err := matcher.Match(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}
	return &SnapshotResult{
		Match:          match,
		FailureMessage: matcher.FailureMessage(nil),
	}, nil
}

func defaultSnapshotFilePath(chartPath, valuesFile string) string {
	// if values file is specified, use the directory of the values file as the snapshot directory.
	// otherwise, use the chart directory.
	if valuesFile != "" {
		return snapshotFilePath(path.Dir(valuesFile), valuesFile)
	} else {
		// if remote chart, create output directory
		if _, err := os.Stat(path.Join(chartPath, "Chart.yaml")); os.IsNotExist(err) {
			chartPath = path.Join("__snapshots__", path.Base(chartPath))
		}
		return snapshotFilePath(chartPath, "")
	}
}

func snapshotFileName(valuesFile string) string {
	if valuesFile != "" {
		name := path.Base(valuesFile)
		name = strings.ReplaceAll(name, ".yaml", "")
		return name
	} else {
		return "default"
	}
}

func snapshotFilePath(dir, valuesFile string) string {
	return path.Join(dir, "__snapshots__", snapshotFileName(valuesFile)+".snap")
}
