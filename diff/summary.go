package diff

type Summary struct {
	Diff          bool           `json:"diff"`
	PathSummary   *PathSummary   `json:"endpoints,omitempty"`
	SchemaSummary *SchemaSummary `json:"schemas,omitempty"`
}

type PathSummary struct {
	Added    int `json:"added,omitempty"`
	Deleted  int `json:"deleted,omitempty"`
	Modified int `json:"modified,omitempty"`
}

type SchemaSummary struct {
	Added    int `json:"added,omitempty"`
	Deleted  int `json:"deleted,omitempty"`
	Modified int `json:"modified,omitempty"`
}
