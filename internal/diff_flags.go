package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type DiffFlags struct {
	CommonFlags

	base     *load.Source
	revision *load.Source
}

func NewDiffFlags() *DiffFlags {
	return &DiffFlags{
		CommonFlags: NewCommonFlags(),
	}
}

func (flags *DiffFlags) toConfig() *diff.Config {
	v := flags.getViper()

	config := diff.NewConfig().WithExcludeElements(flags.getExcludeElements())
	config.PathFilter = v.GetString("match-path")
	config.FilterExtension = v.GetString("filter-extension")
	config.PathPrefixBase = v.GetString("prefix-base")
	config.PathPrefixRevision = v.GetString("prefix-revision")
	config.PathStripPrefixBase = v.GetString("strip-prefix-base")
	config.PathStripPrefixRevision = v.GetString("strip-prefix-revision")
	config.IncludePathParams = v.GetBool("include-path-params")

	return config
}

func (flags *DiffFlags) getComposed() bool {
	return flags.getViper().GetBool("composed")
}

func (flags *DiffFlags) getBase() *load.Source {
	return flags.base
}

func (flags *DiffFlags) getRevision() *load.Source {
	return flags.revision
}

func (flags *DiffFlags) getFlattenAllOf() bool {
	return flags.getViper().GetBool("flatten-allof") || flags.getViper().GetBool("flatten")
}

func (flags *DiffFlags) getFlattenParams() bool {
	return flags.getViper().GetBool("flatten-params")
}

func (flags *DiffFlags) getCaseInsensitiveHeaders() bool {
	return flags.getViper().GetBool("case-insensitive-headers")
}

func (flags *DiffFlags) getIncludeChecks() []string {
	return nil
}

func (flags *DiffFlags) getDeprecationDaysBeta() uint {
	return 0
}

func (flags *DiffFlags) getDeprecationDaysStable() uint {
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
	return flags.getViper().GetString("format")
}

func (flags *DiffFlags) getFailOn() string {
	return ""
}

func (flags *DiffFlags) getLevel() string {
	return ""
}

func (flags *DiffFlags) getFailOnDiff() bool {
	return flags.getViper().GetBool("fail-on-diff")
}

func (flags *DiffFlags) getSeverityLevelsFile() string {
	return ""
}

func (flags *DiffFlags) getExcludeElements() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("exclude-elements"))
}

func (flags *DiffFlags) setBase(source *load.Source) {
	flags.base = source
}

func (flags *DiffFlags) setRevision(source *load.Source) {
	flags.revision = source
}

func (flags *DiffFlags) addExcludeElements(element string) {
	flags.getViper().Set("exclude-elements", append(flags.getViper().GetStringSlice("exclude-elements"), element))
}
