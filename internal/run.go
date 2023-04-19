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

func Run(args []string, stdout io.Writer) (bool, bool, *ReturnError) {

	inputFlags, returnErr := parseFlags(args)
	if returnErr != nil {
		return false, false, returnErr
	}

	if inputFlags.version {
		fmt.Fprintf(stdout, "oasdiff version: %s\n", build.Version)
		return false, false, nil
	}

	if err := validateFlags(inputFlags); err != nil {
		return false, false, err
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
			return false, false, getErrFailedToLoadSpec("base", inputFlags.base, err)
		}

		s2, err := load.FromGlob(loader, inputFlags.revision)
		if err != nil {
			return false, false, getErrFailedToLoadSpec("revision", inputFlags.revision, err)
		}
		diffReport, operationsSources, err = diff.GetPathsDiff(config, s1, s2)
		if err != nil {
			return false, false, getErrDiffFailed(err)
		}
	} else {
		s1, err := checker.LoadOpenAPISpecInfo(inputFlags.base)
		if err != nil {
			return false, false, getErrFailedToLoadSpec("base", inputFlags.base, err)
		}

		s2, err := checker.LoadOpenAPISpecInfo(inputFlags.revision)
		if err != nil {
			return false, false, getErrFailedToLoadSpec("revision", inputFlags.revision, err)
		}
		diffReport, operationsSources, err = diff.GetWithOperationsSourcesMap(config, s1, s2)
		if err != nil {
			return false, false, getErrDiffFailed(err)
		}
	}

	if inputFlags.checkBreaking {
		diffEmpty, returnError := handleBreakingChanges(stdout, diffReport, operationsSources, inputFlags.includeChecks, inputFlags.format, inputFlags.failOnWarns, inputFlags.lang, inputFlags.warnIgnoreFile, inputFlags.errIgnoreFile)
		return inputFlags.failOnDiff, diffEmpty, returnError
	}

	if inputFlags.summary {
		if err := printYAML(stdout, diffReport.GetSummary()); err != nil {
			return false, false, getErrFailedPrint("summary", err)
		}
		return inputFlags.failOnDiff, diffReport.Empty(), nil
	}

	return inputFlags.failOnDiff, diffReport.Empty(), handleDiff(stdout, diffReport, inputFlags.format)
}
