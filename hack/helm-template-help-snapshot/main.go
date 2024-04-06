package main

import (
	"log/slog"
	"os"
	"os/exec"

	"github.com/jlandowner/helm-chartsnap/pkg/snap"
)

func main() {
	snapshot("helm-version", execute("helm", "version"))
	snapshot("helm-template", execute("helm", "template", "--help"))
}

func execute(cmd ...string) string {
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		slog.Error("exec error", "err", err)
		os.Exit(9)
	}
	return string(out)
}

func snapshot(id, data string) {
	s := snap.SnapshotMatcher("helm-template.snap", snap.WithSnapshotID(id))
	match, err := s.Match(data)

	if err != nil {
		slog.Error("snapshot error", "err", err)
		os.Exit(9)
	}
	if !match {
		slog.Error(s.FailureMessage(nil))
		os.Exit(1)
	}
}
