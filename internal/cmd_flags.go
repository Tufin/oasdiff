package internal

import (
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/formatters"
)

func addCommonDiffFlags(cmd *cobra.Command, flags Flags) {
	cmd.PersistentFlags().BoolVarP(flags.refComposed(), "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	cmd.PersistentFlags().StringVarP(flags.refMatchPath(), "match-path", "p", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(flags.refFilterExtension(), "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().StringVarP(flags.refPrefixBase(), "prefix-base", "", "", "add this prefix to paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(flags.refPrefixRevision(), "prefix-revision", "", "", "add this prefix to paths in revised-spec before comparison")
	cmd.PersistentFlags().StringVarP(flags.refStripPrefixBase(), "strip-prefix-base", "", "", "strip this prefix from paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(flags.refStripPrefixRevision(), "strip-prefix-revision", "", "", "strip this prefix from paths in revised-spec before comparison")
	cmd.PersistentFlags().BoolVarP(flags.refIncludePathParams(), "include-path-params", "", false, "include path parameter names in endpoint matching")
	cmd.PersistentFlags().BoolVarP(flags.refFlattenAllOf(), "flatten-allof", "", false, "merge subschemas under allOf before diff")
	cmd.PersistentFlags().BoolVarP(flags.refFlattenParams(), "flatten-params", "", false, "merge common parameters at path level with operation parameters")
	cmd.PersistentFlags().BoolVarP(flags.refInsensitiveHeaders(), "case-insensitive-headers", "", false, "case-insensitive header name comparison")

	addHiddenFlattenFlag(cmd, flags)
	addHiddenCircularDepFlag(cmd)
}

// addHiddenFlattenFlag adds --flatten as a hidden flag
// --flatten was replaced by --flatten-allof
// we still accept --flatten as a synonym for --flatten-allof to avoid breaking existing scripts
func addHiddenFlattenFlag(cmd *cobra.Command, flags Flags) {
	cmd.PersistentFlags().BoolVarP(flags.refFlattenAllOf(), "flatten", "", false, "merge subschemas under allOf before diff")
	hideFlag(cmd, "flatten")
}

// addHiddenCircularDepFlag adds --max-circular-dep as a hidden flag
// --max-circular-dep is no longer needed because kin-openapi3 handles circular references automatically since https://github.com/getkin/kin-openapi/pull/970
// we still accept --max-circular-dep to avoid breaking existing scripts, but we ignore this flag
func addHiddenCircularDepFlag(cmd *cobra.Command) {
	var dummy int
	cmd.PersistentFlags().IntVarP(&dummy, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	hideFlag(cmd, "max-circular-dep")
}

// hideFlag hides a flag from the help
// this is an alternative to marking the flag as deprecated
// marking the flag as deprecated is problematic because it causes cobra to write an error message to stdout which messes up the json and yaml output
func hideFlag(cmd *cobra.Command, flag string) {
	if err := cmd.PersistentFlags().MarkHidden(flag); err != nil {
		// we can ignore this error safely
		_ = err
	}
}

func addCommonBreakingFlags(cmd *cobra.Command, flags Flags) {
	enumWithOptions(cmd, newEnumValue(localizations.GetSupportedLanguages(), localizations.LangDefault, flags.refLang()), "lang", "l", "language for localized output")
	cmd.PersistentFlags().StringVarP(flags.refErrIgnoreFile(), "err-ignore", "", "", "configuration file for ignoring errors")
	cmd.PersistentFlags().StringVarP(flags.refWarnIgnoreFile(), "warn-ignore", "", "", "configuration file for ignoring warnings")
	cmd.PersistentFlags().VarP(newEnumSliceValue(checker.GetOptionalRuleIds(), nil, flags.refIncludeChecks()), "include-checks", "i", "comma-separated list of optional checks (run 'oasdiff checks --required false' to see options)")
	hideFlag(cmd, "include-checks")
	cmd.PersistentFlags().UintVarP(flags.refDeprecationDaysBeta(), "deprecation-days-beta", "", checker.DefaultBetaDeprecationDays, "min days required between deprecating a beta resource and removing it")
	cmd.PersistentFlags().UintVarP(flags.refDeprecationDaysStable(), "deprecation-days-stable", "", checker.DefaultStableDeprecationDays, "min days required between deprecating a stable resource and removing it")
	enumWithOptions(cmd, newEnumValue([]string{"auto", "always", "never"}, "auto", flags.refColor()), "color", "", "when to colorize textual output")
	enumWithOptions(cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputChangelog), string(formatters.FormatText), flags.refFormat()), "format", "f", "output format")
	cmd.PersistentFlags().StringVarP(flags.refSeverityLevelsFile(), "severity-levels", "", "", "configuration file for custom severity levels")
}
