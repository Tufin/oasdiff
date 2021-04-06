package diff

// Endpoints is a list of endpoints
type Endpoints []Endpoint

// Len implements the sort.Interface interface
func (endpoints Endpoints) Len() int {
	return len(endpoints)
}

// Less implements the sort.Interface interface
func (endpoints Endpoints) Less(i, j int) bool {
	if endpoints[i].Path == endpoints[j].Path {
		return endpoints[i].Method < endpoints[j].Method
	}
	return endpoints[i].Path < endpoints[j].Path
}

// Swap implements the sort.Interface interface
func (endpoints Endpoints) Swap(i, j int) {
	endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
}
