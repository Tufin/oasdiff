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
	Added    int `json:"added,omitempty"`    // how many items were added
	Deleted  int `json:"deleted,omitempty"`  // how many items were deleted
	Modified int `json:"modified,omitempty"` // how many items were modified
}

type componentWithSummary interface {
	getSummary() *SummaryDetails
}

func (summary *Summary) add(component componentWithSummary, name string) {
	if !isNilPointer(component) {
		summary.Components[name] = component.getSummary()
	}
}

func isNilPointer(i interface{}) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()
}
