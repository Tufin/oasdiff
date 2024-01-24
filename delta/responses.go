package delta

import (
	"github.com/tufin/oasdiff/diff"
)

func getResponsesDelta(asymmetric bool, d *diff.ResponsesDiff) *WeightedDelta {
	if d.Empty() {
		return &WeightedDelta{}
	}

	added := d.Added.Len()
	deleted := d.Deleted.Len()
	modified := len(d.Modified)
	unchanged := d.Unchanged.Len()
	all := added + deleted + modified + unchanged

	// TODO: drill down into modified
	modifiedDelta := coefficient * float64(modified)

	return NewWeightedDelta(ratio(asymmetric, added, deleted, modifiedDelta, all), all)
}
