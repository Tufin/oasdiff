package internal

import (
	"github.com/spf13/viper"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type Flags interface {
	toConfig() *diff.Config

	getViper() *viper.Viper
	getComposed() bool
	getBase() *load.Source
	getRevision() *load.Source
	getFlattenAllOf() bool
	getFlattenParams() bool
	getCaseInsensitiveHeaders() bool
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
	getSeverityLevelsFile() string
	getAttributes() []string
	getExcludeElements() []string

	setBase(source *load.Source)
	setRevision(source *load.Source)

	addExcludeElements(string)
}

type CommonFlags struct {
	v *viper.Viper
}

func NewCommonFlags() CommonFlags {
	return CommonFlags{
		v: viper.New(),
	}
}

func (flags *CommonFlags) getViper() *viper.Viper {
	return flags.v
}

func (flags *CommonFlags) getAttributes() []string {
	return flags.getViper().GetStringSlice("attributes")
}
