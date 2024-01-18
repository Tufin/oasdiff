package delta

import (
	"github.com/tufin/oasdiff/diff"
)

const coefficient = 0.5

func Get(asymmetric bool, diffReport *diff.Diff) float64 {
	if diffReport.Empty() {
		return 0
	}

	delta := getEndpointsDelta(asymmetric, diffReport.EndpointsDiff)

	return delta
}

func ratio(asymmetric bool, added int, deleted int, modifiedDelta float64, all int) float64 {
	if asymmetric {
		added = 0
	}

	return (float64(added+deleted) + modifiedDelta) / float64(all)
}
