package internal

import (
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
	failOn                   Level
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

func (flags *ChangelogFlags) getComposed() bool {
	return flags.composed
}

func (flags *ChangelogFlags) getBase() string {
	return flags.base
}

func (flags *ChangelogFlags) getRevision() string {
	return flags.revision
}
