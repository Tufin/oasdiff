package delta

import (
	"github.com/tufin/oasdiff/diff"
)

func getParametersDelta(asymmetric bool, d *diff.ParametersDiffByLocation) *WeightedDelta {
	if d.Empty() {
		return &WeightedDelta{}
	}

	added := d.Added.Len()
	deleted := d.Deleted.Len()
	modified := d.Modified.Len()
	unchanged := d.Unchanged.Len()
	all := added + deleted + modified + unchanged

	modifiedDelta := coefficient * getModifiedParametersDelta(asymmetric, d.Modified)

	return NewWeightedDelta(
		ratio(asymmetric, added, deleted, modifiedDelta, all),
		all,
	)
}

func getModifiedParametersDelta(asymmetric bool, d diff.ParamDiffByLocation) float64 {
	weightedDeltas := make([]*WeightedDelta, len(d))
	i := 0
	for _, paramsDiff := range d {
		for _, parameterDiff := range paramsDiff {
			weightedDeltas[i] = NewWeightedDelta(getModifiedParameterDelta(asymmetric, parameterDiff), 1)
			i++
		}
	}
	return weightedAverage(weightedDeltas)
}

func getModifiedParameterDelta(asymmetric bool, d *diff.ParameterDiff) float64 {
	if d.Empty() {
		return 0.0
	}

	// TODO: consider additional elements of ParameterDiff
	schemaDelta := getSchemaDelta(asymmetric, d.SchemaDiff)

	return schemaDelta
}
