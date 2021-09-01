package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func build() (string, bool) {
	buildLog("Building...")

	cmd := exec.Command("go", "build", "-gcflags=\"all=-N -l\"", "-o", buildPath(), root())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	buildLog("Delve debuggins...")

	cmd = exec.Command("dlv", "--listen=:40000", "--continue", "--headless", "--accept-multiclient", "--only-same-user=false", "--check-go-version=false", "--api-version=2", "exec", "--output", buildPath())

	stderr, err = cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	io.Copy(os.Stdout, stdout)
	errBuf, _ = ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}
