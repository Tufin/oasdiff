package diff

// Summary summarizes the changes between two specs
type Summary struct {
	Diff             bool            `json:"diff"`
	PathSummary      *SummaryDetails `json:"paths,omitempty"`
	SchemaSummary    *SummaryDetails `json:"schemas,omitempty"`
	ParameterSummary *SummaryDetails `json:"parameters,omitempty"`
	ResponsesSummary *SummaryDetails `json:"responses,omitempty"`
}

// SummaryDetails summarizes the path changes between two specs
type SummaryDetails struct {
	Added    int `json:"added,omitempty"`
	Deleted  int `json:"deleted,omitempty"`
	Modified int `json:"modified,omitempty"`
}
