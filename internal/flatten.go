package internal

import (
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/flatten"
	"github.com/tufin/oasdiff/formatters"
	"github.com/tufin/oasdiff/load"
)

func getFlattenCmd() *cobra.Command {

	flags := FlattenFlags{}

	cmd := cobra.Command{
		Use:   "flatten spec",
		Short: "Merge allOf",
		Long: `Display a flattened version of the given OpenAPI spec by merging all instances of allOf.
Spec can be a path to a file or a URL.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			flags.spec = args[0]

			// by now flags have been parsed successfully, so we don't need to show usage on any errors
			cmd.Root().SilenceUsage = true

			err := runFlatten(&flags, cmd.OutOrStdout())
			if err != nil {
				setReturnValue(cmd, err.Code)
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().VarP(newEnumValue([]string{string(formatters.FormatYAML), string(formatters.FormatJSON)}, string(formatters.FormatYAML), &flags.format), "format", "f", "output format: yaml or json")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")

	return &cmd
}

func runFlatten(flags *FlattenFlags, stdout io.Writer) *ReturnError {

	openapi3.CircularReferenceCounter = flags.circularReferenceCounter

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	spec, err := load.LoadSpecInfo(loader, flags.spec)
	if err != nil {
		return getErrFailedToLoadSpec("original", flags.spec, err)
	}

	// TODO: get the original format of the spec
	format := flags.format

	flatSpec, err := flatten.MergeSpec(spec.Spec)
	if err != nil {
		return getErrFailedToFlattenSpec("original", flags.spec, err)
	}

	if returnErr := outputFlattenedSpec(formatters.Format(format), stdout, flatSpec); returnErr != nil {
		return returnErr
	}

	return nil
}

func outputFlattenedSpec(format formatters.Format, stdout io.Writer, spec *openapi3.T) *ReturnError {
	switch format {
	case formatters.FormatYAML:
		if err := printYAML(stdout, spec); err != nil {
			return getErrFailedPrint("flattened spec YAML", err)
		}
	case formatters.FormatJSON:
		if err := printJSON(stdout, spec); err != nil {
			return getErrFailedPrint("flattened spec JSON", err)
		}
	default:
		return getErrUnsupportedFormat(string(format))
	}

	return nil
}
