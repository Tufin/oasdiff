package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type ChangelogFlags struct {
	base                     load.Source
	revision                 load.Source
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
	failOn                   string
	flatten                  bool
	lang                     string
	errIgnoreFile            string
	warnIgnoreFile           string
	deprecationDaysBeta      int
	deprecationDaysStable    int
}

func (flags *ChangelogFlags) toConfig() *diff.Config {
	config := diff.NewConfig().WithCheckBreaking().WithExcludeElements(flags.excludeElements)
	config.PathFilter = flags.matchPath
	config.FilterExtension = flags.filterExtension
	config.PathPrefixBase = flags.prefixBase
	config.PathPrefixRevision = flags.prefixRevision
	config.PathStripPrefixBase = flags.stripPrefixBase
	config.PathStripPrefixRevision = flags.stripPrefixRevision
	config.IncludePathParams = flags.includePathParams

	return config
}

func (flags *ChangelogFlags) getComposed() bool {
	return flags.composed
}

func (flags *ChangelogFlags) getBase() load.Source {
	return flags.base
}

func (flags *ChangelogFlags) getRevision() load.Source {
	return flags.revision
}

func (flags *ChangelogFlags) getFlatten() bool {
	return flags.flatten
}
