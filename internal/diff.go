package internal

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/diff"
)

func getDiffCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   `diff`,
		Short: "Generate a diff report",
		Run:   getDiff,
	}
	cmd.PersistentFlags().StringP("base", "b", "", "path or URL (or a glob in Composed mode) of original OpenAPI spec in YAML or JSON format")
	cmd.PersistentFlags().StringP("revision", "r", "", "path or URL (or a glob in Composed mode) of revised OpenAPI spec in YAML or JSON format")
	cmd.PersistentFlags().StringP("format", "f", "yaml", "output format: yaml, json, text or html")
	return &cmd
}

func getDiff(cmd *cobra.Command, args []string) {

	failEmpty, returnErr := runDiffInternal(cmd, args)
	exit(failEmpty, returnErr, cmd.ErrOrStderr())
}

func parseDiffFlags(cmd *cobra.Command) (*InputFlags, *ReturnError) {
	inputFlags := InputFlags{}
	inputFlags.base = cmd.Flag("base").Value.String()
	inputFlags.revision = cmd.Flag("revision").Value.String()
	inputFlags.format = cmd.Flag("format").Value.String()
	return &inputFlags, nil
}

func runDiffInternal(cmd *cobra.Command, args []string) (bool, *ReturnError) {

	inputFlags, returnErr := parseDiffFlags(cmd)

	if returnErr != nil {
		return false, returnErr
	}

	// if returnErr := validateCobraFlags(inputFlags); returnErr != nil {
	// 	return false, returnErr
	// }

	openapi3.CircularReferenceCounter = inputFlags.circularReferenceCounter

	config := generateConfig(inputFlags)

	var diffReport *diff.Diff

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if inputFlags.composed {
		var err *ReturnError
		if diffReport, _, err = composedDiff(loader, inputFlags.base, inputFlags.revision, config); err != nil {
			return false, err
		}
	} else {
		var err *ReturnError
		if diffReport, _, err = normalDiff(loader, inputFlags.base, inputFlags.revision, config); err != nil {
			return false, err
		}
	}

	if inputFlags.summary {
		if err := printYAML(cmd.OutOrStdout(), diffReport.GetSummary()); err != nil {
			return false, getErrFailedPrint("summary", err)
		}
		return failEmpty(inputFlags.failOnDiff, diffReport.Empty()), nil
	}

	return failEmpty(inputFlags.failOnDiff, diffReport.Empty()), handleDiff(cmd.OutOrStdout(), diffReport, inputFlags.format)
}
