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

	cmd := cobra.Command{
		Use:   "flatten spec",
		Short: "Merge allOf",
		Long: `Display a flattened version of the given OpenAPI spec by merging all instances of allOf.
Spec can be a path to a file, a URL or '-' to read standard input.
`,
		Args: cobra.ExactArgs(1),
		RunE: getRun(runFlatten),
	}

	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputFlatten), string(formatters.FormatJSON)), "format", "f", "output format")
	addHiddenCircularDepFlag(&cmd)

	return &cmd
}

func runFlatten(flags *Flags, stdout io.Writer) (bool, *ReturnError) {

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	spec, err := load.NewSpecInfo(loader, flags.getBase(), load.WithFlattenAllOf())
	if err != nil {
		return false, getErrFailedToLoadSpec("original", flags.getBase(), err)
	}

	// TODO: get the original format of the spec
	format := flags.getFormat()

	if returnErr := outputFlattenedSpec(stdout, spec.Spec, format); returnErr != nil {
		return false, returnErr
	}

	return false, nil
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
