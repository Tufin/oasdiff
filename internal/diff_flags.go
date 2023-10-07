package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type DiffFlags struct {
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
	failOnDiff               bool
	flatten                  bool
	circularReferenceCounter int
	includePathParams        bool
	excludeElements          []string
}

func (flags *DiffFlags) toConfig() *diff.Config {
	config := diff.NewConfig().WithExcludeElements(flags.excludeElements)
	config.PathFilter = flags.matchPath
	config.FilterExtension = flags.filterExtension
	config.PathPrefixBase = flags.prefixBase
	config.PathPrefixRevision = flags.prefixRevision
	config.PathStripPrefixBase = flags.stripPrefixBase
	config.PathStripPrefixRevision = flags.stripPrefixRevision
	config.IncludePathParams = flags.includePathParams

	return config
}

func (flags *DiffFlags) getComposed() bool {
	return flags.composed
}

func (flags *DiffFlags) getBase() load.Source {
	return flags.base
}

func (flags *DiffFlags) getRevision() load.Source {
	return flags.revision
}

func (flags *DiffFlags) getFlatten() bool {
	return flags.flatten
}
