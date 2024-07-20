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
	getIncludeChecks() []string
	getDeprecationDaysBeta() uint
	getDeprecationDaysStable() uint
	getLang() string
	getColor() string
	getWarnIgnoreFile() string
	getErrIgnoreFile() string
	getFormat() string
	getFailOn() string
	getLevel() string
	getFailOnDiff() bool
	getAsymmetric() bool
	getSeverityLevelsFile() string
	getAttributes() []string

	setBase(source *load.Source)
	setRevision(source *load.Source)

	addExcludeElements(string)

	refComposed() *bool
	refExcludeElements() *[]string
	refMatchPath() *string
	refFilterExtension() *string
	refPrefixBase() *string
	refPrefixRevision() *string
	refStripPrefixBase() *string
	refStripPrefixRevision() *string
	refIncludePathParams() *bool
	refFlattenAllOf() *bool
	refFlattenParams() *bool
	refInsensitiveHeaders() *bool
	refLang() *string
	refFormat() *string
	refErrIgnoreFile() *string
	refWarnIgnoreFile() *string
	refIncludeChecks() *[]string
	refDeprecationDaysBeta() *uint
	refDeprecationDaysStable() *uint
	refColor() *string
	refSeverityLevelsFile() *string
	refAttributes() *[]string
}
