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
	flattenAllOf             bool
	flattenParams            bool
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

func (flags *DiffFlags) getFlattenAllOf() bool {
	return flags.flattenAllOf
}

func (flags *DiffFlags) getFlattenParams() bool {
	return flags.flattenParams
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

func (flags *DiffFlags) getAsymmetric() bool {
	return false
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

func (flags *DiffFlags) refComposed() *bool {
	return &flags.composed
}

func (flags *DiffFlags) refExcludeElements() *[]string {
	return &flags.excludeElements
}

func (flags *DiffFlags) refMatchPath() *string {
	return &flags.matchPath
}

func (flags *DiffFlags) refFilterExtension() *string {
	return &flags.filterExtension
}

func (flags *DiffFlags) refCircularReferenceCounter() *int {
	return &flags.circularReferenceCounter
}

func (flags *DiffFlags) refPrefixBase() *string {
	return &flags.prefixBase
}

func (flags *DiffFlags) refPrefixRevision() *string {
	return &flags.prefixRevision
}

func (flags *DiffFlags) refStripPrefixBase() *string {
	return &flags.stripPrefixBase
}

func (flags *DiffFlags) refStripPrefixRevision() *string {
	return &flags.stripPrefixRevision
}

func (flags *DiffFlags) refIncludePathParams() *bool {
	return &flags.includePathParams
}

func (flags *DiffFlags) refFlattenAllOf() *bool {
	return &flags.flattenAllOf
}

func (flags *DiffFlags) refFlattenParams() *bool {
	return &flags.flattenParams
}

func (flags *DiffFlags) refLang() *string {
	return nil
}

func (flags *DiffFlags) refErrIgnoreFile() *string {
	return nil
}

func (flags *DiffFlags) refWarnIgnoreFile() *string {
	return nil
}

func (flags *DiffFlags) refIncludeChecks() *[]string {
	return nil
}

func (flags *DiffFlags) refDeprecationDaysBeta() *int {
	return nil
}

func (flags *DiffFlags) refDeprecationDaysStable() *int {
	return nil
}

func (flags *DiffFlags) refColor() *string {
	return nil
}
