package internal

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func Run(shell string, commands []string, out io.Writer, env map[string]string) error {
	shellCmd := strings.Split(shell, " ")
	envSlice := make([]string, 0, len(env))
	for key, value := range env {
		envSlice = append(envSlice, key+"="+value)
	}
	for _, command := range commands {
		Debug("running command: %s\n", command)
		cmdLine := append(shellCmd, command)
		cmd := exec.Command(cmdLine[0], cmdLine[1:]...)
		cmd.Env = append(os.Environ(), envSlice...)
		cmd.Stdout = out
		cmd.Stderr = out
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
