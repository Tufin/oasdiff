package main

import (
	"os"

	"github.com/tufin/oasdiff/internal"
)

func main() {
	internal.Run(os.Args, os.Stdout, os.Stderr)
}
