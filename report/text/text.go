package text

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/tufin/oasdiff/diff"
)

// Report is a simplified OpenAPI diff report in text format
type Report struct {
	Writer io.Writer
	level  int
}

func (report *Report) indent() *Report {
	return &Report{
		Writer: report.Writer,
		level:  report.level + 1,
	}
}

func (report *Report) print(output ...interface{}) (n int, err error) {
	return fmt.Fprintln(report.Writer, addPrefix(report.level, output)...)
}

func addPrefix(level int, output []interface{}) []interface{} {
	return append(getPrefix(level), output...)
}

func getPrefix(level int) []interface{} {
	if level == 1 {
		return []interface{}{"*"}
	}

	if level > 1 {
		return []interface{}{strings.Repeat("  ", level-1) + "-"}
	}

	return []interface{}{}
}

// Output outputs a textual diff report
func (report *Report) Output(d *diff.Diff) {

	if d.Empty() {
		report.print("No changes")
		return
	}

	if d.EndpointsDiff.Empty() {
		report.print("No endpoint changes")
		return
	}

	report.print("### New Endpoints")
	report.print("-----------------")
	for _, added := range d.EndpointsDiff.Added {
		report.print(added.Method, added.Path)
	}
	report.print("")

	report.print("### Deleted Endpoints")
	report.print("---------------------")
	for _, deleted := range d.EndpointsDiff.Deleted {
		report.print(deleted.Method, deleted.Path)
	}
	report.print("")

	report.print("### Modified Endpoints")
	report.print("----------------------")
	for endpoint, methodDiff := range d.EndpointsDiff.Modified {
		report.print(endpoint.Method, endpoint.Path)
		report.indent().printMethod(methodDiff)
		report.print("")
	}
}

func (report *Report) printMethod(d *diff.MethodDiff) {
	if d.Empty() {
		return
	}

	report.printValue(d.DescriptionDiff, "Description")

	report.printParams(d.ParametersDiff)

	if !d.RequestBodyDiff.Empty() {
		report.print("Request body changed")
	}

	if !d.ResponsesDiff.Empty() {
		report.print("Response changed")
		report.indent().printResponses(d.ResponsesDiff)
	}

	if !d.CallbacksDiff.Empty() {
		report.print("Callbacks changed")
	}

	if !d.SecurityDiff.Empty() {
		report.print("Security changed")
	}
}

func (report *Report) printParams(d *diff.ParametersDiff) {
	if d.Empty() {
		return
	}

	for location, params := range d.Added {
		for _, param := range params {
			report.print("New", location, "param:", param)
		}
	}

	for location, params := range d.Deleted {
		for _, param := range params {
			report.print("Deleted", location, "param:", param)
		}
	}

	for location, paramDiffs := range d.Modified {
		for param, paramDiff := range paramDiffs {
			report.print("Modified", location, "param:", param)
			report.indent().printParam(paramDiff)
		}
	}
}

func (report *Report) printParam(d *diff.ParameterDiff) {
	if !d.SchemaDiff.Empty() {
		report.print("Schema changed")
		report.indent().printSchema(d.SchemaDiff)
	}

	if !d.ContentDiff.Empty() {
		report.print("Content changed")
		report.indent().printContent(d.ContentDiff)
	}
}

func (report *Report) printSchema(d *diff.SchemaDiff) {
	if d.Empty() {
		return
	}

	if d.SchemaAdded {
		report.print("Schema added")
	}

	if d.SchemaDeleted {
		report.print("Schema deleted")
	}

	report.printValue(d.TypeDiff, "Type")
	report.printValue(d.TitleDiff, "Title")
	report.printValue(d.FormatDiff, "Format")
	report.printValue(d.DescriptionDiff, "Description")
	report.printValue(d.DefaultDiff, "Default")

	if !d.EnumDiff.Empty() {
		if len(d.EnumDiff.Added) > 0 {
			report.print("New enum values:", d.EnumDiff.Added)
		}
		if len(d.EnumDiff.Deleted) > 0 {
			report.print("Deleted enum values:", d.EnumDiff.Deleted)
		}
	}

	if !d.AdditionalPropertiesAllowedDiff.Empty() {
		report.print("Additional properties changed from", d.AdditionalPropertiesAllowedDiff.From, "to", d.AdditionalPropertiesAllowedDiff.To)
	}

	report.printValue(d.DeprecatedDiff, "Deprecated")
	report.printValue(d.MinDiff, "Min")
	report.printValue(d.MaxDiff, "Max")
	report.printValue(d.MultipleOfDiff, "MultipleOf")
	report.printValue(d.MinLengthDiff, "MinLength")
	report.printValue(d.MaxLengthDiff, "MaxLength")
	report.printValue(d.PatternDiff, "Pattern")
	report.printValue(d.MinItemsDiff, "MinItems")
	report.printValue(d.MaxItemsDiff, "MaxItems")

	if !d.ItemsDiff.Empty() {
		report.print("Items changed")
		report.indent().printSchema(d.ItemsDiff)
	}

	if !d.PropertiesDiff.Empty() {
		report.print("Properties changed")
	}
}

func quote(value interface{}) interface{} {
	if reflect.ValueOf(value).Kind() == reflect.String {
		return "'" + value.(string) + "'"
	}
	return value
}

func (report *Report) printResponses(d *diff.ResponsesDiff) {
	if d.Empty() {
		return
	}

	for _, added := range d.Added {
		report.print("New response:", added)
	}

	for _, deleted := range d.Deleted {
		report.print("Deleted response:", deleted)
	}

	for response, responseDiff := range d.Modified {
		report.print("Modified response:", response)
		report.indent().printResponse(responseDiff)
	}
}

func (report *Report) printResponse(d *diff.ResponseDiff) {
	if d.Empty() {
		return
	}

	report.printValue(d.DescriptionDiff, "Description")

	if !d.ContentDiff.Empty() {
		report.print("Content changed")
		report.indent().printContent(d.ContentDiff)
	}

	if !d.HeadersDiff.Empty() {
		report.print("Headers changed")
		report.indent().printHeaders(d.HeadersDiff)
	}
}

func (report *Report) printContent(d *diff.ContentDiff) {
	if d.Empty() {
		return
	}

	if !d.SchemaDiff.Empty() {
		report.print("Schema changed")
		report.indent().printSchema(d.SchemaDiff)
	}

	if !d.EncodingsDiff.Empty() {
		report.print("Encodings changed")
	}
}

func (report *Report) printValue(d *diff.ValueDiff, title string) {
	if d.Empty() {
		return
	}

	report.print(title, "changed from", quote(d.From), "to", quote(d.To))
}

func (report *Report) printHeaders(d *diff.HeadersDiff) {
	if d.Empty() {
		return
	}

	for _, added := range d.Added {
		report.print("New header:", added)
	}

	for _, deleted := range d.Deleted {
		report.print("Deleted header:", deleted)
	}

	for header := range d.Modified {
		report.print("Modified header:", header)
	}
}
