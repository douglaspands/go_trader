package main

import (
	"os"
	"trader/internal/core"
)

var exit = os.Exit

func run(args []string) int {
	return core.NewApp().Run(args)
}

func main() {
	exit(run(os.Args[1:]))
}
