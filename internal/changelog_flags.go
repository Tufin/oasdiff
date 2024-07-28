package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type ChangelogFlags struct {
	CommonFlags

	base     *load.Source
	revision *load.Source
}

func NewChangelogFlags() *ChangelogFlags {
	return &ChangelogFlags{
		CommonFlags: NewCommonFlags(),
	}
}

func (flags *ChangelogFlags) toConfig() *diff.Config {
	v := flags.getViper()

	config := diff.NewConfig()
	config.PathFilter = v.GetString("match-path")
	config.FilterExtension = v.GetString("filter-extension")
	config.PathPrefixBase = v.GetString("prefix-base")
	config.PathPrefixRevision = v.GetString("prefix-revision")
	config.PathStripPrefixBase = v.GetString("strip-prefix-base")
	config.PathStripPrefixRevision = v.GetString("strip-prefix-revision")
	config.IncludePathParams = v.GetBool("include-path-params")

	return config
}

func (flags *ChangelogFlags) getComposed() bool {
	return flags.getViper().GetBool("composed")
}

func (flags *ChangelogFlags) getBase() *load.Source {
	return flags.base
}

func (flags *ChangelogFlags) getRevision() *load.Source {
	return flags.revision
}

func (flags *ChangelogFlags) getFlattenAllOf() bool {
	return flags.getViper().GetBool("flatten-allof") || flags.getViper().GetBool("flatten")
}

func (flags *ChangelogFlags) getFlattenParams() bool {
	return flags.getViper().GetBool("flatten-params")
}

func (flags *ChangelogFlags) getCaseInsensitiveHeaders() bool {
	return flags.getViper().GetBool("case-insensitive-headers")
}

func (flags *ChangelogFlags) getIncludeChecks() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("include-checks"))
}

func (flags *ChangelogFlags) getExcludeElements() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("exclude-elements"))
}

func (flags *ChangelogFlags) getDeprecationDaysBeta() uint {
	return flags.getViper().GetUint("deprecation-days-beta")
}

func (flags *ChangelogFlags) getDeprecationDaysStable() uint {
	return flags.getViper().GetUint("deprecation-days-stable")
}

func (flags *ChangelogFlags) getLang() string {
	return flags.getViper().GetString("lang")
}

func (flags *ChangelogFlags) getColor() string {
	return flags.getViper().GetString("color")
}

func (flags *ChangelogFlags) getWarnIgnoreFile() string {
	return flags.getViper().GetString("warn-ignore")
}

func (flags *ChangelogFlags) getErrIgnoreFile() string {
	return flags.getViper().GetString("err-ignore")
}

func (flags *ChangelogFlags) getFormat() string {
	return flags.getViper().GetString("format")
}

func (flags *ChangelogFlags) getFailOn() string {
	return flags.getViper().GetString("fail-on")
}

func (flags *ChangelogFlags) getLevel() string {
	return flags.getViper().GetString("level")
}

func (flags *ChangelogFlags) getFailOnDiff() bool {
	return flags.getViper().GetBool("fail-on-diff")
}

func (flags *ChangelogFlags) getSeverityLevelsFile() string {
	return flags.getViper().GetString("severity-levels")
}

func (flags *ChangelogFlags) setBase(source *load.Source) {
	flags.base = source
}

func (flags *ChangelogFlags) setRevision(source *load.Source) {
	flags.revision = source
}

func (flags *ChangelogFlags) addExcludeElements(element string) {
	flags.getViper().Set("exclude-elements", append(flags.getViper().GetStringSlice("exclude-elements"), element))
}
