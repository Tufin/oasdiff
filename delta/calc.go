package delta

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

func Get(diffReport *diff.Diff, base, revision *openapi3.T) float64 {
	if diffReport.Empty() {
		return 0
	}

	d := getEndpointDelta(diffReport.EndpointsDiff, base.Paths, revision.Paths)

	return d
}

func getEndpointDelta(diff *diff.EndpointsDiff, paths1, paths2 *openapi3.Paths) float64 {
	if diff.Empty() {
		return 0
	}

	endpoints1 := countEndpoints(paths1)
	endpoints2 := countEndpoints(paths2)

	added := len(diff.Added)
	deleted := len(diff.Deleted)

	return average(devide(added, endpoints2), devide(deleted, endpoints1))

}

func countEndpoints(paths *openapi3.Paths) int {
	count := 0
	for _, pathItem := range paths.Map() {
		count = count + len(pathItem.Operations())
	}
	return count
}

func devide(a, b int) float64 {
	return float64(a) / float64(b)
}

func average(a, b float64) float64 {
	return (a + b) / 2
}
