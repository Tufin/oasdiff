package main

import (
	"fmt"
	"os"

	"github.com/tufin/oasdiff/internal"
)

func main() {
	exit(internal.Run(os.Args, os.Stdout))
}

func exit(failOnDiff bool, diffEmpty bool, err *internal.ReturnError) {
	if err != nil {
		exitWithError(err)
	}
	exitNormally(failOnDiff, diffEmpty)
}

func exitWithError(err *internal.ReturnError) {
	fmt.Fprintf(os.Stderr, "%v\n", err.Err)
	os.Exit(err.Code)
}

func exitNormally(failOnDiff, diffEmpty bool) {
	if failOnDiff && !diffEmpty {
		os.Exit(1)
	}
	os.Exit(0)
}
