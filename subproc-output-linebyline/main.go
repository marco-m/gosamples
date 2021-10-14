package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
)

func main() {
	if err := drive(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func drive() error {
	log := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Debug,
		Color:  hclog.AutoColor,
		Output: os.Stdout,
	})

	cmd := exec.Command("ls", "-1", "/")
	return LogRun(log, cmd)
}

func LogRun(log hclog.Logger, cmd *exec.Cmd) error {
	log = log.Named("LogRun")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %s", err)
	}

	scanOut := bufio.NewScanner(stdout)
	for scanOut.Scan() {
		line := scanOut.Text()
		if line != "" {
			log.Debug("", "stdout", truncate(line, 80))
		}
	}

	scanErr := bufio.NewScanner(stderr)
	for scanErr.Scan() {
		line := scanErr.Text()
		if line != "" {
			log.Debug("", "stderr", truncate(line, 80))
		}
	}

	// Collect errors
	var bld strings.Builder
	fmt.Println(bld.String())

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(&bld, "wait: %s", err)
	}
	if err := scanOut.Err(); err != nil {
		fmt.Fprintf(&bld, "; stdout: %s", err)
	}
	if err := scanErr.Err(); err != nil {
		fmt.Fprintf(&bld, "; stderr: %s", err)
	}
	if bld.Len() > 0 {
		return fmt.Errorf("%s", bld.String())
	}

	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
