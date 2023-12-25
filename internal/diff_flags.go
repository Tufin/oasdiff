package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type DiffFlags struct {
	base                     *load.Source
	revision                 *load.Source
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

func (flags *DiffFlags) getBase() *load.Source {
	return flags.base
}

func (flags *DiffFlags) getRevision() *load.Source {
	return flags.revision
}

func (flags *DiffFlags) getFlatten() bool {
	return flags.flatten
}

func (flags *DiffFlags) getCircularReferenceCounter() int {
	return flags.circularReferenceCounter
}

func (flags *DiffFlags) getIncludeChecks() []string {
	return nil
}

func (flags *DiffFlags) getDeprecationDaysBeta() int {
	return 0
}

func (flags *DiffFlags) getDeprecationDaysStable() int {
	return 0
}

func (flags *DiffFlags) getLang() string {
	return ""
}

func (flags *DiffFlags) getColor() string {
	return ""
}

func (flags *DiffFlags) getTags() []string {
	return []string{""}
}

func (flags *DiffFlags) getWarnIgnoreFile() string {
	return ""
}

func (flags *DiffFlags) getErrIgnoreFile() string {
	return ""
}

func (flags *DiffFlags) getFormat() string {
	return flags.format
}

func (flags *DiffFlags) getFailOn() string {
	return ""
}

func (flags *DiffFlags) getFailOnDiff() bool {
	return flags.failOnDiff
}

func (flags *DiffFlags) setBase(source *load.Source) {
	flags.base = source
}

func (flags *DiffFlags) setRevision(source *load.Source) {
	flags.revision = source
}

func (flags *DiffFlags) addExcludeElements(element string) {
	flags.excludeElements = append(flags.excludeElements, element)
}
