package internal

import (
	"github.com/tufin/oasdiff/diff"
)

type DiffFlags struct {
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
	failOnDiff               bool
	mergeAllOf               bool
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
	config.MergeAllOf = flags.mergeAllOf

	return config
}

func (flags *DiffFlags) getComposed() bool {
	return flags.composed
}

func (flags *DiffFlags) getBase() string {
	return flags.base
}

func (flags *DiffFlags) getRevision() string {
	return flags.revision
}
