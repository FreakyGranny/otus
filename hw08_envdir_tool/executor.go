package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envVar, val := range env {
		if val == "" {
			os.Unsetenv(envVar)

			continue
		}
		os.Setenv(envVar, val)
	}
	cmdName := cmd[0]
	subProc := exec.Command(cmdName)
	subProc.Args = cmd[0:]
	subProc.Stdout = os.Stdout
	subProc.Stderr = os.Stderr
	subProc.Stdin = os.Stdin

	if err := subProc.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}

	return
}
