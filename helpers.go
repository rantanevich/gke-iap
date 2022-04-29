package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func CheckPackagesInstalled(pkgs []string) {
	ok := true

	for _, pkg := range pkgs {
		if !IsCommandAvailable(pkg) {
			fmt.Fprintf(os.Stderr, "%q: executable file not found in PATH", pkg)
		}
	}

	if !ok {
		os.Exit(1)
	}
}

func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	if err != nil {
		return false
	}
	return true
}

func Exec(args []string, msg string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if msg != "" {
		fmt.Printf(msg)
	}

	if err := cmd.Run(); err != nil {
		return "", ConcatenateError(err, stderr.String())
	}

	return stdout.String(), nil
}

func ConcatenateError(err error, stderr string) error {
	if len(stderr) == 0 {
		return err
	}
	return fmt.Errorf("%w - %s", err, stderr)
}
