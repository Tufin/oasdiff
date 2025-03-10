package main

import (
	"os"

	"github.com/oasdiff/oasdiff/internal"
)

func main() {
	os.Exit(internal.Run(os.Args, os.Stdout, os.Stderr))
}
