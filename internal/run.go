package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/build"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {

	failEmpty, returnErr := runInternal(args, stdout, stderr)

	if returnErr != nil {
		if returnErr.Err != nil {
			fmt.Fprintf(stderr, "%v\n", returnErr.Err)
		}
		return returnErr.Code
	}

	if failEmpty {
		return 1
	}

	return 0
}

func runInternal(args []string, stdout io.Writer, stderr io.Writer) (bool, *ReturnError) {

	inputFlags, returnErr := parseFlags(args, stdout)

	if returnErr != nil {
		return false, returnErr
	}

	if inputFlags.version {
		fmt.Fprintf(stdout, "oasdiff version: %s\n", build.Version)
		return false, nil
	}

	if returnErr := validateFlags(inputFlags); returnErr != nil {
		return false, returnErr
	}

	openapi3.CircularReferenceCounter = inputFlags.circularReferenceCounter

	config := generateConfig(inputFlags)

	var diffReport *diff.Diff
	var operationsSources *diff.OperationsSourcesMap

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if inputFlags.composed {
		var err *ReturnError
		if diffReport, operationsSources, err = composedDiff(loader, inputFlags.base, inputFlags.revision, config); err != nil {
			return false, err
		}
	} else {
		var err *ReturnError
		if diffReport, operationsSources, err = normalDiff(loader, inputFlags.base, inputFlags.revision, config); err != nil {
			return false, err
		}
	}

	if inputFlags.checkBreaking || inputFlags.changelog {
		diffEmpty, returnError := handleBreakingChanges(stdout, diffReport, operationsSources, inputFlags)
		return failEmpty(inputFlags.failOnDiff, diffEmpty), returnError
	}

	if inputFlags.summary {
		if err := printYAML(stdout, diffReport.GetSummary()); err != nil {
			return false, getErrFailedPrint("summary", err)
		}
		return failEmpty(inputFlags.failOnDiff, diffReport.Empty()), nil
	}

	return failEmpty(inputFlags.failOnDiff, diffReport.Empty()), handleDiff(stdout, diffReport, inputFlags.format)
}

func normalDiff(loader load.Loader, base, revision string, config *diff.Config) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {
	s1, err := load.LoadOpenAPISpecInfo(loader, base)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("base", base, err)
	}
	s2, err := load.LoadOpenAPISpecInfo(loader, revision)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("revision", revision, err)
	}

	diffReport, operationsSources, err := diff.GetWithOperationsSourcesMap(config, s1, s2)
	if err != nil {
		return nil, nil, getErrDiffFailed(err)
	}

	return diffReport, operationsSources, nil
}

func composedDiff(loader load.Loader, base, revision string, config *diff.Config) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {
	s1, err := load.FromGlob(loader, base)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("base", base, err)
	}

	s2, err := load.FromGlob(loader, revision)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("revision", revision, err)
	}
	diffReport, operationsSources, err := diff.GetPathsDiff(config, s1, s2)
	if err != nil {
		return nil, nil, getErrDiffFailed(err)
	}

	return diffReport, operationsSources, nil
}

func failEmpty(failOnDiff, diffEmpty bool) bool {
	return failOnDiff && !diffEmpty
}
