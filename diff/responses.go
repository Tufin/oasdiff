package diff

type Responses struct {
	Added    map[string]ResponseNames `json:"added,omitempty"`    // key is param location (path, query, header or cookie)
	Deleted  map[string]ResponseNames `json:"deleted,omitempty"`  // key is param location
	Modified map[string]ResponseDiffs `json:"modified,omitempty"` // key is param location
}

// ResponseNames is a set of parameter names
type ResponseNames map[string]struct{}

// ResponseDiffs is map of parameter names to their respective diffs
type ResponseDiffs map[string]ResponseDiff

type ResponseDiff struct {
}
