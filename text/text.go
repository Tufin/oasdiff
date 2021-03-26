package text

import (
	"fmt"
	"io"

	"github.com/tufin/oasdiff/diff"
)

// Print outputs a simple diff report as text
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
	fmt.Fprintln(writer, "---")
	for _, added := range d.EndpointsDiff.Added {
		fmt.Fprintln(writer, added.Method, added.Path)
	}
	fmt.Fprintln(writer, "")

	fmt.Fprintln(writer, "### What's Deprecated")
	fmt.Fprintln(writer, "---")
	for _, deleted := range d.EndpointsDiff.Deleted {
		fmt.Fprintln(writer, deleted.Method, deleted.Path)
	}
	fmt.Fprintln(writer, "")

	fmt.Fprintln(writer, "### What's Changed")
	fmt.Fprintln(writer, "---")
	for endpoint := range d.EndpointsDiff.Modified {
		fmt.Fprintln(writer, endpoint.Method, endpoint.Path)
	}
}
