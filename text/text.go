package text

import (
	"fmt"
	"io"

	"github.com/tufin/oasdiff/diff"
)

// Print outputs a textual diff report
func Print(d *diff.Diff, writer io.Writer) {

	if d.Empty() {
		fmt.Fprintln(writer, "No changes")
		return
	}

	if d.EndpointsDiff.Empty() {
		fmt.Fprintln(writer, "No endpoint changes")
		return
	}

	fmt.Fprintln(writer, "### What's New")
	fmt.Fprintln(writer, "--------------")
	for _, added := range d.EndpointsDiff.Added {
		fmt.Fprintln(writer, added.Method, added.Path)
	}
	fmt.Fprintln(writer, "")

	fmt.Fprintln(writer, "### What's Deprecated")
	fmt.Fprintln(writer, "---------------------")
	for _, deleted := range d.EndpointsDiff.Deleted {
		fmt.Fprintln(writer, deleted.Method, deleted.Path)
	}
	fmt.Fprintln(writer, "")

	fmt.Fprintln(writer, "### What's Changed")
	fmt.Fprintln(writer, "------------------")
	for endpoint, methodDiff := range d.EndpointsDiff.Modified {
		fmt.Fprintln(writer, endpoint.Method, endpoint.Path)
		printMethod(methodDiff, writer)
		fmt.Fprintln(writer, "")
	}
}

func printMethod(d *diff.MethodDiff, writer io.Writer) {
	if d.Empty() {
		return
	}

	if !d.DescriptionDiff.Empty() {
		fmt.Fprintln(writer, "* Description changed from: ", d.DescriptionDiff.From, "To:", d.DescriptionDiff.To)
	}

	printParams(d.ParametersDiff, writer)

	if !d.RequestBodyDiff.Empty() {
		fmt.Fprintln(writer, "* Request body changed")
	}

	if !d.ResponsesDiff.Empty() {
		fmt.Fprintln(writer, "* Response changed")
	}

	if !d.SecurityDiff.Empty() {
		fmt.Fprintln(writer, "* Security changed")
	}
}

func printParams(d *diff.ParametersDiff, writer io.Writer) {
	if d.Empty() {
		return
	}

	for location, params := range d.Added {
		for _, param := range params {
			fmt.Fprintln(writer, "* New", location, "param:", param)
		}
	}

	for location, params := range d.Deleted {
		for _, param := range params {
			fmt.Fprintln(writer, "* Deprecated", location, "param:", param)
		}
	}

	for location, paramDiffs := range d.Modified {
		for param := range paramDiffs {
			fmt.Fprintln(writer, "* Modified", location, "param:", param)
		}
	}
}
