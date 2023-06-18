package internal

import (
	"fmt"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

type ChangelogFlags struct {
	base                     string
	revision                 string
	composed                 bool
	prefixBase               string
	prefixRevision           string
	stripPrefixBase          string
	stripPrefixRevision      string
	matchPath                string
	filterExtension          string
	format                   string
	circularReferenceCounter int
	includePathParams        bool
	excludeElements          []string
	includeChecks            []string
	failOn                   checker.LevelArg
	lang                     Lang
	errIgnoreFile            string
	warnIgnoreFile           string
}

func (flags *ChangelogFlags) toConfig() *diff.Config {
	config := diff.NewConfig()
	config.PathFilter = flags.matchPath
	config.FilterExtension = flags.filterExtension
	config.PathPrefixBase = flags.prefixBase
	config.PathPrefixRevision = flags.prefixRevision
	config.PathStripPrefixBase = flags.stripPrefixBase
	config.PathStripPrefixRevision = flags.stripPrefixRevision
	config.IncludePathParams = flags.includePathParams
	config.SetExcludeElements(flags.excludeElements)

	return config
}

func (flags *ChangelogFlags) validate() *ReturnError {

	if invalidChecks := checker.ValidateIncludeChecks(flags.includeChecks); len(invalidChecks) > 0 {
		return getErrInvalidFlags(fmt.Errorf("invalid include-checks=%s", flags.includeChecks))
	}
	// if invalidElements := diff.ValidateExcludeElements(flags.excludeElements); len(invalidElements) > 0 {
	// 	return getErrInvalidFlags(fmt.Errorf("invalid exclude-elements=%s", flags.excludeElements))
	// }

	return nil
}
