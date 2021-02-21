package diff

// Summary summarizes the changes between two specs
type Summary struct {
	Diff          bool           `json:"diff"`
	PathSummary   *PathSummary   `json:"paths,omitempty"`
	SchemaSummary *SchemaSummary `json:"schemas,omitempty"`
}

// PathSummary summarizes the path changes between two specs
type PathSummary struct {
	Added    int `json:"added,omitempty"`
	Deleted  int `json:"deleted,omitempty"`
	Modified int `json:"modified,omitempty"`
}

// SchemaSummary summarizes the schema changes between two specs
type SchemaSummary struct {
	Added    int `json:"added,omitempty"`
	Deleted  int `json:"deleted,omitempty"`
	Modified int `json:"modified,omitempty"`
}
