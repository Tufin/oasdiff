package diff

// Summary summarizes the changes between two OAS specs
type Summary struct {
	Diff             bool            `json:"diff"`
	PathSummary      *SummaryDetails `json:"paths,omitempty"`
	SchemaSummary    *SummaryDetails `json:"schemas,omitempty"`
	ParameterSummary *SummaryDetails `json:"parameters,omitempty"`
	HeaderSummary    *SummaryDetails `json:"headers,omitempty"`
	ResponsesSummary *SummaryDetails `json:"responses,omitempty"`
}

// SummaryDetails summarizes the changes between equivalent parts of the two OAS spec: paths, schemas, parameters, headers, responses etc.
type SummaryDetails struct {
	Added    int `json:"added,omitempty"`    // how many items were added
	Deleted  int `json:"deleted,omitempty"`  // how many items were deleted
	Modified int `json:"modified,omitempty"` // how many items were modified
}
