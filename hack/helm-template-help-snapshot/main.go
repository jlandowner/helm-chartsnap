package main

import (
	"log/slog"
	"os"
	"os/exec"
	"regexp"

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
	return replaceHomeDir(replaceHelmEnv(string(out)))
}

func replaceHomeDir(bs string) string {
	// Get os.UserHomeDir() and replace it with "###HOME_DIR###"
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(home)
	return re.ReplaceAllString(bs, "###HOME_DIR###")
}

func replaceHelmEnv(bs string) string {
	// Get helm env and replace HELM_REGISTRY_CONFIG, HELM_REPOSITORY_CACHE and HELM_REPOSITORY_CONFIG
	envs := parseHelmEnvOutput()
	re := regexp.MustCompile(envs["HELM_REGISTRY_CONFIG"])
	bs = re.ReplaceAllString(bs, "###HELM_REGISTRY_CONFIG###")
	re = regexp.MustCompile(envs["HELM_REPOSITORY_CACHE"])
	bs = re.ReplaceAllString(bs, "###HELM_REPOSITORY_CACHE###")
	re = regexp.MustCompile(envs["HELM_REPOSITORY_CONFIG"])
	bs = re.ReplaceAllString(bs, "###HELM_REPOSITORY_CONFIG###")
	return bs
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

func parseHelmEnvOutput() map[string]string {
	out, err := exec.Command("helm", "env").CombinedOutput()
	if err != nil {
		slog.Error("exec error", "err", err)
		os.Exit(9)
	}
	re := regexp.MustCompile(`(.*)="(.*)"`)
	matches := re.FindAllStringSubmatch(string(out), -1)
	envs := make(map[string]string)
	for _, match := range matches {
		envs[match[1]] = match[2]
	}
	return envs
}
