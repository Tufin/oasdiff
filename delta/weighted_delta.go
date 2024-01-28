package delta

type WeightedDelta struct {
	delta  float64
	weight int
}

func NewWeightedDelta(delta float64, weight int) *WeightedDelta {
	return &WeightedDelta{
		delta:  delta,
		weight: weight,
	}
}

func weightedAverage(weightedDeltas []*WeightedDelta) float64 {
	dividend := 0.0
	divisor := 0
	for _, weightedDelta := range weightedDeltas {
		dividend += weightedDelta.delta * float64(weightedDelta.weight)
		divisor += weightedDelta.weight
	}
	if dividend == 0 {
		return 0
	}
	return dividend / float64(divisor)
}
