package text

import (
	"fmt"
	"io"

	"github.com/tufin/oasdiff/diff"
)

// Report is a simplified OpenAPI diff report in text format
type Report struct {
	Writer io.Writer
	level  int
}

// Print outputs a textual diff report
func (report *Report) Print(d *diff.Diff) {

	if d.Empty() {
		fmt.Fprintln(report.Writer, "No changes")
		return
	}

	if d.EndpointsDiff.Empty() {
		fmt.Fprintln(report.Writer, "No endpoint changes")
		return
	}

	fmt.Fprintln(report.Writer, "### New Endpoints")
	fmt.Fprintln(report.Writer, "-----------------")
	for _, added := range d.EndpointsDiff.Added {
		fmt.Fprintln(report.Writer, added.Method, added.Path)
	}
	fmt.Fprintln(report.Writer, "")

	fmt.Fprintln(report.Writer, "### Deleted Endpoints")
	fmt.Fprintln(report.Writer, "---------------------")
	for _, deleted := range d.EndpointsDiff.Deleted {
		fmt.Fprintln(report.Writer, deleted.Method, deleted.Path)
	}
	fmt.Fprintln(report.Writer, "")

	fmt.Fprintln(report.Writer, "### Modified Endpoints")
	fmt.Fprintln(report.Writer, "----------------------")
	for endpoint, methodDiff := range d.EndpointsDiff.Modified {
		fmt.Fprintln(report.Writer, endpoint.Method, endpoint.Path)
		report.printMethod(methodDiff)
		fmt.Fprintln(report.Writer, "")
	}
}

func (report *Report) printMethod(d *diff.MethodDiff) {
	if d.Empty() {
		return
	}

	if !d.DescriptionDiff.Empty() {
		fmt.Fprintln(report.Writer, "* Description changed from: ", d.DescriptionDiff.From, "To:", d.DescriptionDiff.To)
	}

	report.printParams(d.ParametersDiff)

	if !d.RequestBodyDiff.Empty() {
		fmt.Fprintln(report.Writer, "* Request body changed")
	}

	if !d.ResponsesDiff.Empty() {
		fmt.Fprintln(report.Writer, "* Response changed")
		report.printResponses(d.ResponsesDiff)
	}

	if !d.CallbacksDiff.Empty() {
		fmt.Fprintln(report.Writer, "* Callbacks changed")
	}

	if !d.SecurityDiff.Empty() {
		fmt.Fprintln(report.Writer, "* Security changed")
	}
}

func (report *Report) printParams(d *diff.ParametersDiff) {
	if d.Empty() {
		return
	}

	for location, params := range d.Added {
		for _, param := range params {
			fmt.Fprintln(report.Writer, "* New", location, "param:", param)
		}
	}

	for location, params := range d.Deleted {
		for _, param := range params {
			fmt.Fprintln(report.Writer, "* Deleted", location, "param:", param)
		}
	}

	for location, paramDiffs := range d.Modified {
		for param, paramDiff := range paramDiffs {
			fmt.Fprintln(report.Writer, "* Modified", location, "param:", param)
			report.printParam(paramDiff)
		}
	}
}

func (report *Report) printParam(d *diff.ParameterDiff) {
	if !d.SchemaDiff.Empty() {
		fmt.Fprintln(report.Writer, "  - Schema changed")
		report.printSchema(d.SchemaDiff)
	}

	if !d.ContentDiff.Empty() {
		fmt.Fprintln(report.Writer, "  - Content changed")
	}
}

func (report *Report) printSchema(d *diff.SchemaDiff) {
	if d.Empty() {
		return
	}
}

func (report *Report) printResponses(d *diff.ResponsesDiff) {
	if d.Empty() {
		return
	}

	for _, added := range d.Added {
		fmt.Fprintln(report.Writer, "  - New response:", added)
	}

	for _, deleted := range d.Deleted {
		fmt.Fprintln(report.Writer, "  - Deleted response:", deleted)
	}

	for response := range d.Modified {
		fmt.Fprintln(report.Writer, "  - Modified response:", response)
	}
}
