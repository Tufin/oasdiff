package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/build"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {

	failEmpty, returnErr := runInternal(args, stdout)

	if returnErr != nil {
		fmt.Fprintf(stderr, "%v\n", returnErr.Err)
		return returnErr.Code
	}

	if failEmpty {
		return 1
	}

	return 0
}

func runInternal(args []string, stdout io.Writer) (bool, *ReturnError) {

	inputFlags, returnErr := parseFlags(args)
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

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	config := diff.NewConfig()
	config.PathFilter = inputFlags.filter
	config.FilterExtension = inputFlags.filterExtension
	config.PathPrefixBase = inputFlags.prefixBase
	config.PathPrefixRevision = inputFlags.prefixRevision
	config.PathStripPrefixBase = inputFlags.stripPrefixBase
	config.PathStripPrefixRevision = inputFlags.strip_prefix_revision
	config.BreakingOnly = inputFlags.breakingOnly
	config.DeprecationDays = inputFlags.deprecationDays
	config.SetExcludeElements(inputFlags.excludeElements.ToStringSet(), inputFlags.excludeExamples, inputFlags.excludeDescription, inputFlags.excludeEndpoints)

	var diffReport *diff.Diff
	var operationsSources *diff.OperationsSourcesMap

	if inputFlags.checkBreaking {
		config.IncludeExtensions.Add(checker.XStabilityLevelExtension)
		config.IncludeExtensions.Add(diff.SunsetExtension)
		config.IncludeExtensions.Add(checker.XExtensibleEnumExtension)
	}

	if inputFlags.composed {
		s1, err := load.FromGlob(loader, inputFlags.base)
		if err != nil {
			return false, getErrFailedToLoadSpec("base", inputFlags.base, err)
		}

		s2, err := load.FromGlob(loader, inputFlags.revision)
		if err != nil {
			return false, getErrFailedToLoadSpec("revision", inputFlags.revision, err)
		}
		diffReport, operationsSources, err = diff.GetPathsDiff(config, s1, s2)
		if err != nil {
			return false, getErrDiffFailed(err)
		}
	} else {
		s1, err := checker.LoadOpenAPISpecInfo(inputFlags.base)
		if err != nil {
			return false, getErrFailedToLoadSpec("base", inputFlags.base, err)
		}

		s2, err := checker.LoadOpenAPISpecInfo(inputFlags.revision)
		if err != nil {
			return false, getErrFailedToLoadSpec("revision", inputFlags.revision, err)
		}
		diffReport, operationsSources, err = diff.GetWithOperationsSourcesMap(config, s1, s2)
		if err != nil {
			return false, getErrDiffFailed(err)
		}
	}

	if inputFlags.checkBreaking {
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

func failEmpty(failOnDiff, diffEmpty bool) bool {
	return failOnDiff && !diffEmpty
}
