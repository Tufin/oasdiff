package report

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/tufin/oasdiff/diff"
)

type report struct {
	Writer io.Writer
	level  int
}

func (r *report) indent() *report {
	return &report{
		Writer: r.Writer,
		level:  r.level + 1,
	}
}

func (r *report) print(output ...interface{}) (n int, err error) {
	return fmt.Fprintln(r.Writer, addPrefix(r.level, output)...)
}

func addPrefix(level int, output []interface{}) []interface{} {
	return append(getPrefix(level), output...)
}

func getPrefix(level int) []interface{} {
	if level == 1 {
		return []interface{}{"-"}
	}

	if level > 1 {
		return []interface{}{strings.Repeat("  ", level-1) + "-"}
	}

	return []interface{}{}
}

func (r *report) output(d *diff.Diff) {

	if d.Empty() {
		r.print("No changes")
		return
	}

	if d.EndpointsDiff.Empty() {
		r.print("No endpoint changes")
		return
	}

	r.printTitle("New Endpoints", len(d.EndpointsDiff.Added))
	sort.Sort(d.EndpointsDiff.Added)
	for _, added := range d.EndpointsDiff.Added {
		r.print(added.Method, added.Path, " ")
	}
	r.print("")

	r.printTitle("Deleted Endpoints", len(d.EndpointsDiff.Deleted))
	sort.Sort(d.EndpointsDiff.Deleted)
	for _, deleted := range d.EndpointsDiff.Deleted {
		r.print(deleted.Method, deleted.Path, " ")
	}
	r.print("")

	r.printTitle("Modified Endpoints", len(d.EndpointsDiff.Modified))
	keys := d.EndpointsDiff.Modified.ToEndpoints()
	sort.Sort(keys)
	for _, endpoint := range keys {
		r.print(endpoint.Method, endpoint.Path)
		r.indent().printMethod(d.EndpointsDiff.Modified[endpoint])
		r.print("")
	}
}

func (r *report) printMethod(d *diff.MethodDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.DescriptionDiff, "Description")
	r.printParams(d.ParametersDiff)

	if !d.RequestBodyDiff.Empty() {
		r.print("Request body changed")
		r.indent().printRequestBody(d.RequestBodyDiff)
	}

	if !d.ResponsesDiff.Empty() {
		r.print("Responses changed")
		r.indent().printResponses(d.ResponsesDiff)
	}

	r.printMessage(d.CallbacksDiff, "Callbacks changed")

	if !d.SecurityDiff.Empty() {
		r.print("Security changed")
		r.indent().printSecurityRequirements(d.SecurityDiff)
	}
}

func (r *report) printParams(d *diff.ParametersDiff) {
	if d.Empty() {
		return
	}

	for location, params := range d.Added {
		for _, param := range params {
			r.print("New", location, "param:", param)
		}
	}

	for location, params := range d.Deleted {
		for _, param := range params {
			r.print("Deleted", location, "param:", param)
		}
	}

	for location, paramDiffs := range d.Modified {
		for param, paramDiff := range paramDiffs {
			r.print("Modified", location, "param:", param)
			r.indent().printParam(paramDiff)
		}
	}
}

func (r *report) printParam(d *diff.ParameterDiff) {
	if !d.SchemaDiff.Empty() {
		r.print("Schema changed")
		r.indent().printSchema(d.SchemaDiff)
	}

	if !d.ContentDiff.Empty() {
		r.print("Content changed")
		r.indent().printContent(d.ContentDiff)
	}
}

func (r *report) printSchema(d *diff.SchemaDiff) {
	if d.Empty() {
		return
	}

	r.printConditional(d.SchemaAdded, "Schema added")
	r.printConditional(d.SchemaDeleted, "Schema deleted")
	r.printConditional(d.CircularRefDiff, "Schema circular referecnce changed")
	r.printMessage(d.OneOfDiff, "Property 'OneOf' changed")
	r.printMessage(d.AnyOfDiff, "Property 'AnyOf' changed")
	r.printMessage(d.AllOfDiff, "Property 'AllOf' changed")

	if !d.NotDiff.Empty() {
		r.print("Property 'Not' changed")
		r.indent().printSchema(d.NotDiff)
	}

	r.printValue(d.TypeDiff, "Type")
	r.printValue(d.TitleDiff, "Title")
	r.printValue(d.FormatDiff, "Format")
	r.printValue(d.DescriptionDiff, "Description")

	if !d.EnumDiff.Empty() {
		r.printConditional(len(d.EnumDiff.Added) > 0, "New enum values:", d.EnumDiff.Added)
		r.printConditional(len(d.EnumDiff.Deleted) > 0, "Deleted enum values:", d.EnumDiff.Deleted)
	}

	r.printValue(d.DefaultDiff, "Default")
	r.printValue(d.AdditionalPropertiesAllowedDiff, "AdditionalProperties")
	r.printValue(d.UniqueItemsDiff, "UniqueItems")
	r.printValue(d.ExclusiveMinDiff, "ExclusiveMin")
	r.printValue(d.ExclusiveMaxDiff, "ExclusiveMax")
	r.printValue(d.NullableDiff, "Nullable")
	r.printValue(d.ReadOnlyDiff, "ReadOnly")
	r.printValue(d.WriteOnlyDiff, "WriteOnly")
	r.printValue(d.AllowEmptyValueDiff, "AllowEmptyValue")
	r.printValue(d.XMLDiff, "XML")
	r.printValue(d.DeprecatedDiff, "Deprecated")
	r.printValue(d.MinDiff, "Min")
	r.printValue(d.MaxDiff, "Max")
	r.printValue(d.MultipleOfDiff, "MultipleOf")
	r.printValue(d.MinLengthDiff, "MinLength")
	r.printValue(d.MaxLengthDiff, "MaxLength")
	r.printValue(d.PatternDiff, "Pattern")
	r.printValue(d.MinItemsDiff, "MinItems")
	r.printValue(d.MaxItemsDiff, "MaxItems")

	if !d.ItemsDiff.Empty() {
		r.print("Items changed")
		r.indent().printSchema(d.ItemsDiff)
	}

	r.printMessage(d.RequiredDiff, "Required changed")
	r.printValue(d.MinPropsDiff, "MinProps")
	r.printValue(d.MaxPropsDiff, "MaxProps")

	if !d.PropertiesDiff.Empty() {
		r.print("Properties changed")
		r.indent().printProperties(d.PropertiesDiff)
	}

	if !d.AdditionalPropertiesDiff.Empty() {
		r.print("AdditionalProperties changed")
		r.indent().printSchema(d.AdditionalPropertiesDiff)
	}

	r.printMessage(d.DiscriminatorDiff, "Discriminator changed")
}

func (r *report) printProperties(d *diff.SchemasDiff) {
	if d.Empty() {
		return
	}

	for _, property := range d.Added {
		r.print("New property:", property)
	}

	for _, property := range d.Deleted {
		r.print("Deleted property:", property)
	}

	for property, schemaDiff := range d.Modified {
		r.print("Modified property:", property)
		r.indent().printSchema(schemaDiff)
	}
}

func quote(value interface{}) interface{} {
	if reflect.ValueOf(value).Kind() == reflect.String {
		return "'" + value.(string) + "'"
	}
	return value
}

func (r *report) printResponses(d *diff.ResponsesDiff) {
	if d.Empty() {
		return
	}

	for _, added := range d.Added {
		r.print("New response:", added)
	}

	for _, deleted := range d.Deleted {
		r.print("Deleted response:", deleted)
	}

	for response, responseDiff := range d.Modified {
		r.print("Modified response:", response)
		r.indent().printResponse(responseDiff)
	}
}

func (r *report) printResponse(d *diff.ResponseDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.DescriptionDiff, "Description")

	if !d.ContentDiff.Empty() {
		r.print("Content changed")
		r.indent().printContent(d.ContentDiff)
	}

	if !d.HeadersDiff.Empty() {
		r.print("Headers changed")
		r.indent().printHeaders(d.HeadersDiff)
	}
}

func (r *report) printRequestBody(d *diff.RequestBodyDiff) {
	if d.Empty() {
		return
	}

	r.printValue(d.DescriptionDiff, "Description")

	if !d.ContentDiff.Empty() {
		r.print("Content changed")
		r.indent().printContent(d.ContentDiff)
	}
}

func (r *report) printContent(d *diff.ContentDiff) {
	if d.Empty() {
		return
	}

	for _, name := range d.MediaTypeAdded {
		r.print("New media type:", name)
	}

	for _, name := range d.MediaTypeDeleted {
		r.print("Deleted media type:", name)
	}

	for name, d := range d.MediaTypeModified {
		r.print("Modified media type:", name)
		r.indent().printMediaType(d)
	}
}

func (r *report) printMediaType(d *diff.MediaTypeDiff) {
	if d.Empty() {
		return
	}

	if !d.SchemaDiff.Empty() {
		r.print("Schema changed")
		r.indent().printSchema(d.SchemaDiff)
	}

	r.printMessage(d.EncodingsDiff, "Encodings changed")
}

func (r *report) printValue(d *diff.ValueDiff, title string) {
	if d.Empty() {
		return
	}

	r.print(title, "changed from", quote(d.From), "to", quote(d.To))
}

func (r *report) printHeaders(d *diff.HeadersDiff) {
	if d.Empty() {
		return
	}

	for _, added := range d.Added {
		r.print("New header:", added)
	}

	for _, deleted := range d.Deleted {
		r.print("Deleted header:", deleted)
	}

	for header := range d.Modified {
		r.print("Modified header:", header)
	}
}

func (r *report) printSecurityRequirements(d *diff.SecurityRequirementsDiff) {
	if d.Empty() {
		return
	}

	for _, added := range d.Added {
		r.print("New security requirements:", added)
	}

	for _, deleted := range d.Deleted {
		r.print("Deleted security requirements:", deleted)
	}

	for securityRequirementID, securityScopesDiff := range d.Modified {
		r.print("Modified security requirements:", securityRequirementID)
		r.indent().printSecurityScopes(securityScopesDiff)
	}
}

func (r *report) printSecurityScopes(d diff.SecurityScopesDiff) {
	for scheme, scopeDiff := range d {
		r.printConditional(len(scopeDiff.Added) > 0, "Scheme", scheme, "Added scopes:", scopeDiff.Added)
		r.printConditional(len(scopeDiff.Deleted) > 0, "Scheme", scheme, "Deleted scopes:", scopeDiff.Deleted)
	}
}

func (r *report) printTitle(title string, count int) {
	text := ""
	if count == 0 {
		text = fmt.Sprintf("### %s: None", title)
	} else {
		text = fmt.Sprintf("### %s: %d", title, count)
	}

	r.print(text)
	r.print(strings.Repeat("-", len(text)))
}

func (r *report) printMessage(d diff.IDiff, output ...interface{}) {
	r.printConditional(!d.Empty(), output...)
}

func (r *report) printConditional(b bool, output ...interface{}) {
	if b {
		r.print(output...)
	}
}
