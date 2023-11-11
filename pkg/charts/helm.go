package charts

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

type HelmTemplateCmdOptions struct {
	HelmPath       string
	ReleaseName    string
	Namespace      string
	Chart          string
	ValuesFile     string
	AdditionalArgs []string

	log *slog.Logger
}

func (o *HelmTemplateCmdOptions) SetLogger(log *slog.Logger) {
	o.log = log
}

func (o *HelmTemplateCmdOptions) Log() *slog.Logger {
	if o.log == nil {
		o.log = slog.Default()
	}
	return o.log
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
	o.Log().DebugContext(ctx, "executing 'helm template' command", "args", args, "AdditionalArgs", o.AdditionalArgs)

	// helm template should not be executed in debug mode because YAML parser fails.
	os.Setenv("HELM_DEBUG", "false")

	cmd := exec.CommandContext(ctx, o.HelmPath, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}
