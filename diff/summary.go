package diff

import "reflect"

// Summary summarizes the changes between two OpenAPI specifications
type Summary struct {
	Diff       bool                       `json:"diff"`
	Components map[string]*SummaryDetails `json:"components,omitempty"`
}

func newSummary() *Summary {
	return &Summary{
		Components: map[string]*SummaryDetails{},
	}
}

// SummaryDetails summarizes the changes between equivalent parts of the two OpenAPI specifications: paths, schemas, parameters, headers, responses etc.
type SummaryDetails struct {
	Added    int `json:"added,omitempty"`    // number of added items
	Deleted  int `json:"deleted,omitempty"`  // number of deleted items
	Modified int `json:"modified,omitempty"` // number of modified items
}

type componentWithSummary interface {
	getSummary() *SummaryDetails
}

// Component names
const (
	PathsComponent         = "paths"
	TagsComponent          = "tags"
	SchemasComponent       = "schemas"
	ParametersComponent    = "parameters"
	HeadersComponent       = "headers"
	RequestBodiesComponent = "requestBodies"
	ResponsesComponent     = "responses"
	CallbacksComponent     = "callbacks"
)

// GetSummaryDetails returns the summary for a specific component
func (summary *Summary) GetSummaryDetails(component string) SummaryDetails {
	if details, ok := summary.Components[component]; ok {
		if details != nil {
			return *details
		}
	}

	return SummaryDetails{}
}

func (summary *Summary) add(component componentWithSummary, name string) {
	if !isNilPointer(component) {
		summary.Components[name] = component.getSummary()
	}
}

func isNilPointer(i interface{}) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()
}
