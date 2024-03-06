package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type Flags interface {
	toConfig() *diff.Config

	getComposed() bool
	getBase() *load.Source
	getRevision() *load.Source
	getFlattenAllOf() bool
	getFlattenParams() bool
	getInsensitiveHeaders() bool
	getCircularReferenceCounter() int
	getIncludeChecks() []string
	getDeprecationDaysBeta() int
	getDeprecationDaysStable() int
	getLang() string
	getColor() string
	getWarnIgnoreFile() string
	getErrIgnoreFile() string
	getFormat() string
	getFailOn() string
	getFailOnDiff() bool
	getAsymmetric() bool

	setBase(source *load.Source)
	setRevision(source *load.Source)
	setUnchanged(bool)

	addExcludeElements(string)

	refComposed() *bool
	refExcludeElements() *[]string
	refMatchPath() *string
	refFilterExtension() *string
	refCircularReferenceCounter() *int
	refPrefixBase() *string
	refPrefixRevision() *string
	refStripPrefixBase() *string
	refStripPrefixRevision() *string
	refIncludePathParams() *bool
	refFlattenAllOf() *bool
	refFlattenParams() *bool
	refInsensitiveHeaders() *bool
	refLang() *string
	refErrIgnoreFile() *string
	refWarnIgnoreFile() *string
	refIncludeChecks() *[]string
	refDeprecationDaysBeta() *int
	refDeprecationDaysStable() *int
	refColor() *string
}
