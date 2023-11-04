package charts

import (
	"context"
	"fmt"
	"os/exec"
)

type HelmTemplateCmdOptions struct {
	HelmPath       string
	ReleaseName    string
	Namespace      string
	Chart          string
	ValuesFile     string
	AdditionalArgs []string
}

func (o *HelmTemplateCmdOptions) Execute(ctx context.Context) ([]byte, error) {
	args := []string{
		"template", o.ReleaseName, o.Chart,
	}
	if o.Namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%s", o.Namespace))
	}
	if o.ValuesFile != "" {
		args = append(args, fmt.Sprintf("--values=%s", o.ValuesFile))
	}
	if len(o.AdditionalArgs) > 0 {
		args = append(args, o.AdditionalArgs...)
	}
	cmd := exec.CommandContext(ctx, o.HelmPath, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}
