package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type ChangelogFlags struct {
	base                  *load.Source
	revision              *load.Source
	composed              bool
	prefixBase            string
	prefixRevision        string
	stripPrefixBase       string
	stripPrefixRevision   string
	matchPath             string
	filterExtension       string
	format                string
	includePathParams     bool
	excludeElements       []string
	includeChecks         []string
	failOn                string
	level                 string
	flattenAllOf          bool
	flattenParams         bool
	insensitiveHeaders    bool
	lang                  string
	errIgnoreFile         string
	warnIgnoreFile        string
	deprecationDaysBeta   uint
	deprecationDaysStable uint
	color                 string
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

func (flags *ChangelogFlags) getFlattenAllOf() bool {
	return flags.flattenAllOf
}

func (flags *ChangelogFlags) getFlattenParams() bool {
	return flags.flattenParams
}

func (flags *ChangelogFlags) getInsensitiveHeaders() bool {
	return flags.insensitiveHeaders
}

func (flags *ChangelogFlags) getIncludeChecks() []string {
	return flags.includeChecks
}

func (flags *ChangelogFlags) getDeprecationDaysBeta() uint {
	return flags.deprecationDaysBeta
}

func (flags *ChangelogFlags) getDeprecationDaysStable() uint {
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

func (flags *ChangelogFlags) getLevel() string {
	return flags.level
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

func (flags *ChangelogFlags) refComposed() *bool {
	return &flags.composed
}

func (flags *ChangelogFlags) refExcludeElements() *[]string {
	return &flags.excludeElements
}

func (flags *ChangelogFlags) refMatchPath() *string {
	return &flags.matchPath
}

func (flags *ChangelogFlags) refFilterExtension() *string {
	return &flags.filterExtension
}

func (flags *ChangelogFlags) refPrefixBase() *string {
	return &flags.prefixBase
}

func (flags *ChangelogFlags) refPrefixRevision() *string {
	return &flags.prefixRevision
}

func (flags *ChangelogFlags) refStripPrefixBase() *string {
	return &flags.stripPrefixBase
}

func (flags *ChangelogFlags) refStripPrefixRevision() *string {
	return &flags.stripPrefixRevision
}

func (flags *ChangelogFlags) refIncludePathParams() *bool {
	return &flags.includePathParams
}

func (flags *ChangelogFlags) refFlattenAllOf() *bool {
	return &flags.flattenAllOf
}

func (flags *ChangelogFlags) refFlattenParams() *bool {
	return &flags.flattenParams
}

func (flags *ChangelogFlags) refInsensitiveHeaders() *bool {
	return &flags.insensitiveHeaders
}

func (flags *ChangelogFlags) refLang() *string {
	return &flags.lang
}

func (flags *ChangelogFlags) refFormat() *string {
	return &flags.format
}

func (flags *ChangelogFlags) refErrIgnoreFile() *string {
	return &flags.errIgnoreFile
}

func (flags *ChangelogFlags) refWarnIgnoreFile() *string {
	return &flags.warnIgnoreFile
}

func (flags *ChangelogFlags) refIncludeChecks() *[]string {
	return &flags.includeChecks
}

func (flags *ChangelogFlags) refDeprecationDaysBeta() *uint {
	return &flags.deprecationDaysBeta
}

func (flags *ChangelogFlags) refDeprecationDaysStable() *uint {
	return &flags.deprecationDaysStable
}

func (flags *ChangelogFlags) refColor() *string {
	return &flags.color
}
