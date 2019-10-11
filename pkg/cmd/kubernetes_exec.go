package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexellis/go-execute"
)

func fetchChart(path, chart string) error {
	mkErr := os.MkdirAll(path, 0700)

	if mkErr != nil {
		return mkErr
	}

	task := execute.ExecTask{
		Command: fmt.Sprintf("%s fetch %s --untar --untardir %s", localBinary("helm"), chart, path),
		Env:     os.Environ(),
	}
	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func getArchitecture() string {
	res, _ := kubectlTask("get", "nodes", `--output`, `jsonpath={range $.items[0]}{.status.nodeInfo.architecture}`)

	arch := strings.TrimSpace(string(res.Stdout))

	return arch
}

func templateChart(basePath, chart, namespace, outputPath, values string) error {

	mkErr := os.MkdirAll(outputPath, 0700)
	if mkErr != nil {
		return mkErr
	}

	chartRoot := path.Join(basePath, chart)
	task := execute.ExecTask{
		Command: fmt.Sprintf("%s template %s --output-dir %s --values %s --namespace %s",
			localBinary("helm"), chart, outputPath, path.Join(chartRoot, values), namespace),
		Env: os.Environ(),
		Cwd: basePath,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func localBinary(name string) string {
	home := os.Getenv("HOME")
	return path.Join(path.Join(home, ".k3sup/.bin/"), name)
}

func addHelmRepo(name, url string) error {
	task := execute.ExecTask{
		Command: fmt.Sprintf("%s repo add %s %s", localBinary("helm"), name, url),
		Env:     os.Environ(),
	}
	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func updateHelmRepos() error {
	task := execute.ExecTask{
		Command: fmt.Sprintf("%s repo update", localBinary("helm")),
		Env:     os.Environ(),
	}
	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func helmInit() error {
	task := execute.ExecTask{
		Command: fmt.Sprintf("%s", localBinary("helm")),
		Env:     os.Environ(),
		Args:    []string{"init", "--client-only"},
	}
	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func kubectlTask(parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command: "kubectl",
		Args:    parts,
	}

	res, err := task.Execute()

	return res, err
}

func kubectl(parts ...string) error {
	task := execute.ExecTask{
		Command: "kubectl",
		Args:    parts,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}
