package charts

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cosmo-workspace/controller-testtools/pkg/snap"
	"github.com/cosmo-workspace/controller-testtools/pkg/unstructured"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v2"
)

func Snap(ctx context.Context, o HelmTemplateCmdOptions) (match bool, failureMessage string, err error) {
	sv := SnapshotValues{}
	if o.ValuesFile != "" {
		f, err := os.Open(o.ValuesFile)
		if err != nil {
			return match, "", fmt.Errorf("failed to open values file: %w", err)
		}
		defer f.Close()

		err = yaml.NewDecoder(f).Decode(&sv)
		if err != nil {
			return match, "", fmt.Errorf("failed to decode values file: %w", err)
		}
	}
	slog.Debug("test spec from values file", "spec", sv.TestSpec)

	helmTemplate := HelmTemplateCmdOptions{
		HelmPath:    o.HelmPath,
		ReleaseName: o.ReleaseName,
		Namespace:   o.Namespace,
		Chart:       o.Chart,
		ValuesFile:  o.ValuesFile,
	}
	out, err := helmTemplate.Execute(ctx)
	if err != nil {
		return match, "", fmt.Errorf("failed to execute helm template: %w: %s", err, out)
	}
	slog.Debug("helm template output", "output", string(out))

	manifests, err := unstructured.Decode(string(out))
	if err != nil {
		return match, "", fmt.Errorf("failed to load helm output: %w: out='%s'", err, string(out))
	}

	for _, v := range sv.TestSpec.DynamicFields {
		for i, obj := range manifests {
			if v.APIVersion == obj.GetAPIVersion() &&
				v.Kind == obj.GetKind() &&
				v.Name == obj.GetName() {
				for _, p := range v.JSONPath {
					newObj, err := unstructured.Replace(manifests[i], p, "###DYNAMIC_FIELD###")
					if err != nil {
						return match, "", fmt.Errorf("failed to replace json path: %w", err)
					}
					manifests[i] = *newObj
				}
			}
		}
	}
	res, err := unstructured.Encode(manifests)
	if err != nil {
		return match, "", fmt.Errorf("failed to encode manifests: %w", err)
	}

	s := snap.SnapShotMatcher(SnapshotFile(o.Chart, o.ValuesFile), SnapshotID(o.ValuesFile))
	match, err = s.Match(string(res))

	if err != nil {
		return match, "", fmt.Errorf("failed to get snapshot: %w", err)
	}
	return match, s.FailureMessage(out), nil
}

func SnapshotID(valuesFile string) string {
	if valuesFile != "" {
		return strings.ReplaceAll(path.Base(valuesFile), ".yaml", "")
	} else {
		return "default"
	}
}

func SnapshotFile(chartPath, valuesFile string) string {
	if valuesFile != "" {
		return path.Join(path.Dir(valuesFile), "__snapshots__", SnapshotID(valuesFile)+".snap")
	} else {
		return path.Join(chartPath, "__snapshots__", SnapshotID(valuesFile)+".snap")
	}
}
