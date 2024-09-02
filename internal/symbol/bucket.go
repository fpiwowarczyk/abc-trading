package symbol

// Bucket represents a part of the data set, which is last n points.
// How many points is specified by the Size field.
type Bucket struct {
	Size int

	Min    float64
	Max    float64
	Avg    float64
	Varian float64

	Points []float64
	Sum    float64
	SumSq  float64
}

func (b *Bucket) cantFitIntoBucket(addedPointsSize int) bool {
	return len(b.Points)+addedPointsSize > b.Size
}

func (b *Bucket) setStats(min, max, avg, varian, sum, sumSq float64, points []float64) {
	b.Min = min
	b.Max = max
	b.Avg = avg
	b.Varian = varian
	b.Sum = sum
	b.SumSq = sumSq
	b.Points = points
}
