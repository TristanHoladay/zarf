package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Change terminal colors
const colorReset = "\x1b[0m"
const colorGray = "\x1b[30;1m"
const colorRed = "\x1b[31;1m"
const colorGreen = "\x1b[32;1m"
const colorYellow = "\x1b[33;1m"
const colorBlue = "\x1b[34;1m"
const colorPink = "\x1b[35;1m"
const colorCyan = "\x1b[36;1m"
const colorWhite = "\x1b[37;1m"

// SilentExecCommand
func SilentExecCommand(envVariables []string, commandName string, args ...string) error {
	_, err := execCommand(false, envVariables, commandName, args...)
	return err
}

// ExecCommand performs a process execution outside of Go
func ExecCommand(envVariables []string, commandName string, args ...string) (string, error) {
	return execCommand(true, envVariables, commandName, args...)
}

func execCommand(showLogs bool, envVariables []string, commandName string, args ...string) (string, error) {
	if showLogs {
		fmt.Println()
		fmt.Printf("%s", colorGreen)
		fmt.Print(commandName + " ")
		fmt.Printf("%s", colorCyan)
		fmt.Printf("%v", args)
		fmt.Printf("%s", colorWhite)
		fmt.Printf("%s", colorReset)
		fmt.Println("")
	}

	cmd := exec.Command(commandName, args...)
	env := os.Environ()
	env = append(env, envVariables...)
	cmd.Env = env

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if showLogs {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			_, errStdout = io.Copy(stdout, stdoutIn)
			wg.Done()
		}()

		_, errStderr = io.Copy(stderr, stderrIn)
		wg.Wait()
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	if showLogs {
		if errStdout != nil || errStderr != nil {
			return "", errors.New("Unable to capture stdOut or stdErr")
		}
	}

	return stdoutBuf.String(), nil
}