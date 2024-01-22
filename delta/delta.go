package delta

import (
	"github.com/tufin/oasdiff/diff"
)

const coefficient = 0.5

func Get(asymmetric bool, diffReport *diff.Diff) float64 {
	if diffReport.Empty() {
		return 0
	}

	deltaEndpoints := getEndpointsDelta(asymmetric, diffReport.EndpointsDiff)

	return weightedAverage(deltaEndpoints)
}

func ratio(asymmetric bool, added int, deleted int, modifiedDelta float64, all int) float64 {
	if asymmetric {
		added = 0
	}

	return (float64(added+deleted) + modifiedDelta) / float64(all)
}

type WeightedDelta struct {
	delta  float64
	weight int
}

func weightedAverage(weightedDeltas ...WeightedDelta) float64 {
	dividend := 0.0
	divisor := 0
	for _, weightedDelta := range weightedDeltas {
		dividend += weightedDelta.delta
		divisor += weightedDelta.weight
	}
	return dividend / float64(divisor)
}
