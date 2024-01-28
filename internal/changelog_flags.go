package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type ChangelogFlags struct {
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
	color                    string
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

func (flags *ChangelogFlags) getBase() *load.Source {
	return flags.base
}

func (flags *ChangelogFlags) getRevision() *load.Source {
	return flags.revision
}

func (flags *ChangelogFlags) getFlatten() bool {
	return flags.flatten
}

func (flags *ChangelogFlags) getCircularReferenceCounter() int {
	return flags.circularReferenceCounter
}

func (flags *ChangelogFlags) getIncludeChecks() []string {
	return flags.includeChecks
}

func (flags *ChangelogFlags) getDeprecationDaysBeta() int {
	return flags.deprecationDaysBeta
}

func (flags *ChangelogFlags) getDeprecationDaysStable() int {
	return flags.deprecationDaysStable
}

func (flags *ChangelogFlags) getLang() string {
	return flags.lang
}

func (flags *ChangelogFlags) getColor() string {
	return flags.color
}

func (flags *ChangelogFlags) getWarnIgnoreFile() string {
	return flags.warnIgnoreFile
}

func (flags *ChangelogFlags) getErrIgnoreFile() string {
	return flags.errIgnoreFile
}

func (flags *ChangelogFlags) getFormat() string {
	return flags.format
}

func (flags *ChangelogFlags) getFailOn() string {
	return flags.failOn
}

func (flags *ChangelogFlags) getFailOnDiff() bool {
	return false
}

func (flags *ChangelogFlags) getAsymmetric() bool {
	return false
}

func (flags *ChangelogFlags) setBase(source *load.Source) {
	flags.base = source
}

func (flags *ChangelogFlags) setRevision(source *load.Source) {
	flags.revision = source
}

func (flags *ChangelogFlags) addExcludeElements(element string) {
	flags.excludeElements = append(flags.excludeElements, element)
}
