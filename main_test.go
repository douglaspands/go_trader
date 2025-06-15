package main

import (
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// EXPECT
	expected := `Investor Support Tool

  Designed to assist investors of all levels in making 
  informed decisions and optimizing their portfolios. 
  With an intuitive interface and access to real-time data, 
  you will have the necessary information to track the market and 
  analyze opportunities with confidence.

Usage:
  trader [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  reit        Tool to get reit information
  security    Tool to get security information
  stock       Tool to get stock information
  version     Show version

Flags:
  -h, --help   help for trader

Use "trader [command] --help" for more information about a command.
`
	// THEN
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	main()

	// WHEN
	os.Stdout = oldStdout
	w.Close()
	out, _ := io.ReadAll(r)
	outString := string(out)
	if outString != expected {
		t.Errorf("Output expected: \"%s\", but received: \"%s\"", expected, outString)
	}
}
