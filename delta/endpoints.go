package delta

import (
	"github.com/tufin/oasdiff/diff"
)

func getEndpointsDelta(asymmetric bool, d *diff.EndpointsDiff) WeightedDelta {
	if d.Empty() {
		return WeightedDelta{}
	}

	added := len(d.Added)
	deleted := len(d.Deleted)
	modified := len(d.Modified)
	unchanged := len(d.Unchanged)
	all := added + deleted + modified + unchanged

	modifiedDelta := coefficient * getModifiedEndpointsDelta(asymmetric, d.Modified)

	return WeightedDelta{
		delta:  ratio(asymmetric, added, deleted, modifiedDelta, all),
		weight: 1,
	}
}

func getModifiedEndpointsDelta(asymmetric bool, d diff.ModifiedEndpoints) float64 {
	result := 0.0
	for _, methodDiff := range d {
		result += getModifiedEndpointDelta(asymmetric, methodDiff)
	}
	return result
}

func getModifiedEndpointDelta(asymmetric bool, d *diff.MethodDiff) float64 {
	if d.Empty() {
		return 0
	}

	// TODO: consider additional elements of MethodDiff
	paramsDelta := getParametersDelta(asymmetric, d.ParametersDiff)
	responsesDelta := getResponsesDelta(asymmetric, d.ResponsesDiff)

	return weightedAverage(paramsDelta, responsesDelta)
}
