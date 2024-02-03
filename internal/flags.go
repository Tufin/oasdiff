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

	addExcludeElements(string)
}
