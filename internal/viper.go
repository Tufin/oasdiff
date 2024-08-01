package internal

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/formatters"
	"golang.org/x/exp/slices"
)

type IViper interface {
	SetConfigName(in string)
	AddConfigPath(in string)
	ReadInConfig() error
	BindPFlag(key string, flag *pflag.Flag) error
	UnmarshalExact(rawVal any, opts ...viper.DecoderConfigOption) error
}

func RunViper(cmd *cobra.Command, v IViper) *ReturnError {
	if err := readConfFile(v); err != nil {
		return getErrConfigFileProblem(err)
	}

	if err := validate(v); err != nil {
		return getErrConfigFileProblem(err)
	}

	if err := bindFlags(cmd, v); err != nil {
		return getErrConfigFileProblem(err)
	}

	return nil
}

func readConfFile(v IViper) error {
	v.SetConfigName(".oasdiff")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			return fmt.Errorf("read error: %s \n", err)
		}
	}

	return nil
}

func bindFlags(cmd *cobra.Command, v IViper) error {
	var result error
	persitentFlags := cmd.PersistentFlags()
	persitentFlags.VisitAll(func(flag *pflag.Flag) {
		name := flag.Name
		if err := v.BindPFlag(name, persitentFlags.Lookup(name)); err != nil {
			result = fmt.Errorf("error binding flag %q to viper: %v", name, err)
			return
		}
	})

	return result
}

// fixViperStringSlice fixes a limitation in viper that doesn't handle custom flags with multiple values
func fixViperStringSlice(viperString []string) []string {
	// viper returns a slice with a single element if the flag was set with a comma-separated list
	if len(viperString) == 1 {
		return strings.Split(viperString[0], ",")
	}

	return viperString
}

type Config struct {
	Attributes             []string `mapstructure:"attributes"`
	Composed               bool     `mapstructure:"composed"`
	FlattenAllof           bool     `mapstructure:"flatten-allof"`
	FlattenParams          bool     `mapstructure:"flatten-params"`
	CaseInsensitiveHeaders bool     `mapstructure:"case-insensitive-headers"`
	DeprecationDaysBeta    uint     `mapstructure:"deprecation-days-beta"`
	DeprecationDaysStable  uint     `mapstructure:"deprecation-days-stable"`
	Lang                   string   `mapstructure:"lang"`
	Color                  string   `mapstructure:"color"`
	WarnIgnore             string   `mapstructure:"warn-ignore"`
	ErrIgnore              string   `mapstructure:"err-ignore"`
	Format                 string   `mapstructure:"format"`
	FailOn                 string   `mapstructure:"fail-on"`
	Level                  string   `mapstructure:"level"`
	FailOnDiff             bool     `mapstructure:"fail-on-diff"`
	SeverityLevels         string   `mapstructure:"severity-levels"`
	ExcludeElements        []string `mapstructure:"exclude-elements"`
	Severity               []string `mapstructure:"severity"`
	Tags                   []string `mapstructure:"tags"`
	MatchPath              string   `mapstructure:"match-path"`
	FilterExtension        string   `mapstructure:"filter-extension"`
	PrefixBase             string   `mapstructure:"prefix-base"`
	PrefixRevision         string   `mapstructure:"prefix-revision"`
	StripPrefixBase        string   `mapstructure:"strip-prefix-base"`
	StripPrefixRevision    string   `mapstructure:"strip-prefix-revision"`
	IncludePathParams      bool     `mapstructure:"include-path-params"`
}

func validate(v IViper) error {
	var config Config

	if err := v.UnmarshalExact(&config); err != nil {
		return fmt.Errorf("validation error: %s \n", err)
	}

	if err := validateString(localizations.SupportedLanguages, config.Lang, "lang"); err != nil {
		return err
	}

	if err := validateString(checker.SupportedColorValues, config.Color, "color"); err != nil {
		return err
	}

	if err := validateString(formatters.SupportedFormats, config.Format, "format"); err != nil {
		return err
	}

	if err := validateString(SupportedLevels, config.FailOn, "fail-on"); err != nil {
		return err
	}

	if err := validateString(SupportedLevels, config.Level, "level"); err != nil {
		return err
	}

	if err := validateStrings(diff.ExcludeDiffOptions, config.ExcludeElements, "exclude-elements"); err != nil {
		return err
	}

	if err := validateStrings([]string{"info", "warn", "error"}, config.Severity, "severity"); err != nil {
		return err
	}

	if err := validateStrings(getAllTags(), config.Tags, "tags"); err != nil {
		return err
	}

	return nil
}

func validateStrings(allowedValues []string, values []string, name string) error {
	for _, value := range values {
		if err := validateString(allowedValues, value, name); err != nil {
			return err
		}
	}
	return nil
}

func validateString(allowedValues []string, value string, name string) error {
	if value == "" {
		return nil
	}

	if slices.Contains(allowedValues, value) {
		return nil
	}

	return fmt.Errorf("invalid %s %q, allowed values: %v", name, value, strings.Join(allowedValues, ", "))
}
