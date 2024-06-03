package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/formatters"
	"github.com/tufin/oasdiff/load"
)

const flattenCmd = "flatten"

func getFlattenCmd() *cobra.Command {

	flags := FlattenFlags{}

	cmd := cobra.Command{
		Use:   "flatten spec",
		Short: "Merge allOf",
		Long: `Display a flattened version of the given OpenAPI spec by merging all instances of allOf.
Spec can be a path to a file, a URL or '-' to read standard input.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			flags.spec = load.NewSource(args[0])

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

	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputFlatten), string(formatters.FormatJSON), &flags.format), "format", "f", "output format")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")

	return &cmd
}

func runFlatten(flags *FlattenFlags, stdout io.Writer) *ReturnError {

	openapi3.CircularReferenceCounter = flags.circularReferenceCounter

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	spec, err := load.NewSpecInfo(loader, flags.spec, load.WithFlattenAllOf())
	if err != nil {
		return getErrFailedToLoadSpec("original", flags.spec, err)
	}

	// TODO: get the original format of the spec
	format := flags.format

	if returnErr := outputFlattenedSpec(stdout, spec.Spec, format); returnErr != nil {
		return returnErr
	}

	return nil
}

func outputFlattenedSpec(stdout io.Writer, spec *openapi3.T, format string) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.DefaultFormatterOpts())
	if err != nil {
		return getErrUnsupportedFormat(format, flattenCmd)
	}

	// render
	bytes, err := formatter.RenderFlatten(spec, formatters.NewRenderOpts())
	if err != nil {
		return getErrFailedPrint("flatten "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
