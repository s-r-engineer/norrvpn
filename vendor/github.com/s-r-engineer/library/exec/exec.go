package libraryExec

import (
	"os/exec"
	"strings"
)

func Run(command ...string) (string, int, error) {
	return execCommand(getCommand(command))
}

func RunWithStdin(stdin string, command ...string) (string, int, error) {
	cmd := getCommand(command)
	cmd.Stdin = strings.NewReader(stdin)
	return execCommand(cmd)
}

func getCommand(command []string) *exec.Cmd {
	return exec.Command(command[0], command[1:]...)
}

func execCommand(cmd *exec.Cmd) (string, int, error) {
	output, err := cmd.CombinedOutput()
	return string(output), cmd.ProcessState.ExitCode(), err
}
