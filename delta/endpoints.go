package delta

import (
	"github.com/tufin/oasdiff/diff"
)

func getEndpointsDelta(asymmetric bool, d *diff.EndpointsDiff) float64 {
	if d.Empty() {
		return 0
	}

	added := len(d.Added)
	deleted := len(d.Deleted)
	modified := len(d.Modified)
	unchanged := len(d.Unchanged)
	all := added + deleted + modified + unchanged

	modifiedDelta := coefficient * getModifiedEndpointsDelta(asymmetric, d.Modified)

	return ratio(asymmetric, added, deleted, modifiedDelta, all)
}

func getModifiedEndpointsDelta(asymmetric bool, d diff.ModifiedEndpoints) float64 {
	weightedDeltas := make([]*WeightedDelta, len(d))
	i := 0
	for _, methodDiff := range d {
		weightedDeltas[i] = NewWeightedDelta(getModifiedEndpointDelta(asymmetric, methodDiff), 1)
		i++
	}
	return weightedAverage(weightedDeltas)
}

func getModifiedEndpointDelta(asymmetric bool, d *diff.MethodDiff) float64 {
	if d.Empty() {
		return 0
	}

	// TODO: consider additional elements of MethodDiff
	paramsDelta := getParametersDelta(asymmetric, d.ParametersDiff)
	responsesDelta := getResponsesDelta(asymmetric, d.ResponsesDiff)

	return weightedAverage([]*WeightedDelta{paramsDelta, responsesDelta})
}
