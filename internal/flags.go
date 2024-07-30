package internal

import (
	"github.com/spf13/viper"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type Flags struct {
	v        *viper.Viper
	base     *load.Source
	revision *load.Source
}

func NewFlags() *Flags {
	return &Flags{
		v: viper.New(),
	}
}

func (flags *Flags) toConfig() *diff.Config {
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

func (flags *Flags) getViper() *viper.Viper {
	return flags.v
}

func (flags *Flags) getAttributes() []string {
	return flags.getViper().GetStringSlice("attributes")
}

func (flags *Flags) getComposed() bool {
	return flags.getViper().GetBool("composed")
}

func (flags *Flags) getBase() *load.Source {
	return flags.base
}

func (flags *Flags) getRevision() *load.Source {
	return flags.revision
}

func (flags *Flags) getFlattenAllOf() bool {
	return flags.getViper().GetBool("flatten-allof") || flags.getViper().GetBool("flatten")
}

func (flags *Flags) getFlattenParams() bool {
	return flags.getViper().GetBool("flatten-params")
}

func (flags *Flags) getCaseInsensitiveHeaders() bool {
	return flags.getViper().GetBool("case-insensitive-headers")
}

func (flags *Flags) getIncludeChecks() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("include-checks"))
}

func (flags *Flags) getDeprecationDaysBeta() uint {
	return flags.getViper().GetUint("deprecation-days-beta")
}

func (flags *Flags) getDeprecationDaysStable() uint {
	return flags.getViper().GetUint("deprecation-days-stable")
}

func (flags *Flags) getLang() string {
	return flags.getViper().GetString("lang")
}

func (flags *Flags) getColor() string {
	return flags.getViper().GetString("color")
}

func (flags *Flags) getWarnIgnoreFile() string {
	return flags.getViper().GetString("warn-ignore")
}

func (flags *Flags) getErrIgnoreFile() string {
	return flags.getViper().GetString("err-ignore")
}

func (flags *Flags) getFormat() string {
	return flags.getViper().GetString("format")
}

func (flags *Flags) getFailOn() string {
	return flags.getViper().GetString("fail-on")
}

func (flags *Flags) getLevel() string {
	return flags.getViper().GetString("level")
}

func (flags *Flags) getFailOnDiff() bool {
	return flags.getViper().GetBool("fail-on-diff")
}

func (flags *Flags) getSeverityLevelsFile() string {
	return flags.getViper().GetString("severity-levels")
}

func (flags *Flags) getExcludeElements() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("exclude-elements"))
}

func (flags *Flags) setBase(source *load.Source) {
	flags.base = source
}

func (flags *Flags) setRevision(source *load.Source) {
	flags.revision = source
}

func (flags *Flags) addExcludeElements(element string) {
	flags.getViper().Set("exclude-elements", append(flags.getViper().GetStringSlice("exclude-elements"), element))
}

func (flags *Flags) getSeverity() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("severity"))
}

func (flags *Flags) getTags() []string {
	return fixViperStringSlice(flags.getViper().GetStringSlice("tags"))
}
