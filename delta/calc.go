package delta

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func Get(asymmetric bool, diffReport *diff.Diff, base, revision *openapi3.T) float64 {
	if diffReport.Empty() {
		return 0
	}

	d := getEndpointDelta(asymmetric, diffReport.EndpointsDiff, base.Paths, revision.Paths)

	return d
}

func getEndpointDelta(asymmetric bool, diff *diff.EndpointsDiff, paths1, paths2 *openapi3.Paths) float64 {
	if diff.Empty() {
		return 0
	}

	endpoints1 := countEndpoints(paths1)
	endpoints2 := countEndpoints(paths2)

	added := len(diff.Added)
	deleted := len(diff.Deleted)

	if asymmetric {
		return utils.Devide(deleted, endpoints1)
	}
	return utils.Average(utils.Devide(added, endpoints2), utils.Devide(deleted, endpoints1))
}

func countEndpoints(paths *openapi3.Paths) int {
	count := 0
	for _, pathItem := range paths.Map() {
		count = count + len(pathItem.Operations())
	}
	return count
}
