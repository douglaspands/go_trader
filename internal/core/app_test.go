package core_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"trader/internal/core"
)

// capture runs fn while redirecting stdout and stderr, and returns what was written to each.
func capture(t *testing.T, fn func()) (string, string) {
	t.Helper()
	oldStdout, oldStderr := os.Stdout, os.Stderr
	stdoutReader, stdoutWriter, _ := os.Pipe()
	stderrReader, stderrWriter, _ := os.Pipe()
	os.Stdout, os.Stderr = stdoutWriter, stderrWriter
	defer func() { os.Stdout, os.Stderr = oldStdout, oldStderr }()

	fn()

	stdoutWriter.Close()
	stderrWriter.Close()
	var stdout, stderr bytes.Buffer
	io.Copy(&stdout, stdoutReader)
	io.Copy(&stderr, stderrReader)
	return stdout.String(), stderr.String()
}

func TestRunVersion(t *testing.T) {
	var code int

	// WHEN
	stdout, stderr := capture(t, func() { code = core.NewApp().Run([]string{"version"}) })

	// THEN
	if code != 0 {
		t.Errorf("expected exit code 0 and received %d", code)
	}
	if stdout != "development\n" {
		t.Errorf(`expected stdout "development\n" and received %q`, stdout)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr and received %q", stderr)
	}
}

func TestRunWithoutArguments(t *testing.T) {
	var code int

	// WHEN
	stdout, stderr := capture(t, func() { code = core.NewApp().Run([]string{}) })

	// THEN
	if code != 0 {
		t.Errorf("expected exit code 0 and received %d", code)
	}
	for _, command := range []string{"reit", "security", "stock", "version"} {
		if !strings.Contains(stdout, command) {
			t.Errorf("expected help to list %q and received %q", command, stdout)
		}
	}
	if stderr != "" {
		t.Errorf("expected empty stderr and received %q", stderr)
	}
}

func TestRunPurchaseBalanceInfiniteAmount(t *testing.T) {
	var code int

	// WHEN
	stdout, stderr := capture(t, func() { code = core.NewApp().Run([]string{"stock", "purchase-balance", "PETR4", "--amount", "inf"}) })

	// THEN
	if code != 1 {
		t.Errorf("expected exit code 1 and received %d", code)
	}
	if !strings.Contains(stderr, "Error: invalid amount") {
		t.Errorf("expected an invalid amount error in stderr and received %q", stderr)
	}
	if stdout != "" {
		t.Errorf("expected empty stdout and received %q", stdout)
	}
}

func TestRunStockGetWithoutTicker(t *testing.T) {
	var code int

	// WHEN
	_, stderr := capture(t, func() { code = core.NewApp().Run([]string{"stock", "get"}) })

	// THEN
	if code != 1 {
		t.Errorf("expected exit code 1 and received %d", code)
	}
	if !strings.Contains(stderr, "Error: accepts 1 arg(s), received 0") {
		t.Errorf("expected an arguments error in stderr and received %q", stderr)
	}
}
