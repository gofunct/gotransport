package stencil

import (
	"fmt"
	"github.com/gofunct/gofs"
	"os"
	"os/exec"
	"strings"
)

// Bash executes a bash script (string).
func Bash(script string, options ...map[string]interface{}) (string, error) {
	return bash(script, options)
}

// BashOutput executes a bash script and returns the output
func BashOutput(script string, options ...map[string]interface{}) (string, error) {
	if len(options) == 0 {
		options = append(options, M{"$out": CaptureBoth})
	} else {
		options[0]["$out"] = CaptureBoth
	}
	return bash(script, options)
}

// Bash executes a bash string. Use backticks for multiline. To execute as shell script,
// use Run("bash script.sh")
func bash(script string, options []map[string]interface{}) (output string, err error) {
	m, dir, capture, err := parseOptions(options)
	if err != nil {
		return "", err
	}

	if strings.Contains(script, "{{") {
		script, err = gofs.StrTemplate(script, m)
		if err != nil {
			return "", err
		}
	}

	gcmd := &command{
		executable: "bash",
		argv:       []string{"-c", script},
		wd:         dir,
		capture:    capture,
		commandstr: script,
	}

	return gcmd.run()
}

func Shell(args ...string) (stdout string, err error) {
	stdoutb, err := ShellOutput(args...)
	return strings.TrimSpace(string(stdoutb)), err
}

func ShellOutput(args ...string) (stdout []byte, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return stdoutb, nil
}
