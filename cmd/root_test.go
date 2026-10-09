package cmd_test

import (
	"strings"
	"testing"
)

func TestRootWithoutArguments(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run()

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if !strings.Contains(tr.stdout.String(), "Investor Support Tool") {
		t.Errorf("expected the description in stdout and received %q", tr.stdout.String())
	}
	for _, command := range []string{"reit", "security", "stock", "version", "completion", "help"} {
		if !strings.Contains(tr.stdout.String(), command) {
			t.Errorf("expected help to list %q and received %q", command, tr.stdout.String())
		}
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestRootVersion(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("version")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if tr.stdout.String() != "development\n" {
		t.Errorf(`expected stdout "development\n" and received %q`, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}
