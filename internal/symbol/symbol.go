package symbol

import (
	"math"
	"slices"

	"github.com/fpiwowarczyk/abc-trading/internal/calculations"
)

const (
	// MaxK is a exponent of a power of 10, which is used to calculate the maximal number of data
	// that can be stored for single symbol.
	MaxK = 8
)

type Data struct {
	LastPoint float64
	Buckets   []*Bucket

	Points []float64
}

// New creates a new Data strucutre that represends data stored for single symbol.
// During this operation data is divided into buckets, which are used to store last n points.
// Number of buckets is specified by MaxK.
func New(added []float64) *Data {
	var buckets []*Bucket
	for k := range MaxK {
		buckets = append(buckets, &Bucket{
			Size:   int(math.Pow(10, float64(k+1))),
			Points: make([]float64, 0),
		})
	}

	var lastBucket *Bucket
	var statsCalculatedForWholeSet bool

	for _, b := range buckets {
		if statsCalculatedForWholeSet {
			b.setStats(
				lastBucket.Min, lastBucket.Max,
				lastBucket.Avg, lastBucket.Varian,
				lastBucket.Sum, lastBucket.SumSq,
				lastBucket.Points)

			continue
		}

		if b.cantFitIntoBucket(len(added)) {
			b.Points = added[len(added)-b.Size:]
		} else {
			b.Points = added
			statsCalculatedForWholeSet = true
		}

		b.Min, b.Max, b.Sum, b.SumSq = calculations.MinMaxSumSumSq(b.Points)
		b.Avg, b.Varian = calculations.RollingAvgVar(
			0.0, 0.0,
			0.0, 0.0,
			b.Sum, b.SumSq,
			0.0, 0.0,
			len(b.Points))

		lastBucket = b
	}

	return &Data{
		LastPoint: added[len(added)-1],
		Buckets:   buckets,
		Points:    added,
	}
}

// Update appends new data points to the existing data set.
// It also updates statistics for each bucket.
// If number of data points exceeds then maximal number of points that can be stored for single symbol,
// then the oldest points are removed.
func (d *Data) Update(added []float64) *Data {
	d.LastPoint = added[len(added)-1]
	d.Points = append(d.Points, added...)

	if len(d.Points) > int(math.Pow(10, MaxK)) {
		d.Points = d.Points[len(d.Points)-int(math.Pow(10, MaxK)):]
	}

	var lastBucket *Bucket
	var calculatedForWholeSet bool

	for _, b := range d.Buckets {
		if calculatedForWholeSet {
			b.setStats(
				lastBucket.Min, lastBucket.Max,
				lastBucket.Avg, lastBucket.Varian,
				lastBucket.Sum, lastBucket.SumSq,
				lastBucket.Points)

			continue
		}

		var out []float64
		if b.cantFitIntoBucket(len(added)) {
			out = b.Points[:len(b.Points)+len(added)-b.Size]
			b.Points = d.Points[len(d.Points)-b.Size:]
		} else {
			calculatedForWholeSet = true
			b.Points = append(b.Points, added...)
		}

		addedMin, addedMax, addedSum, addedSumSq := calculations.MinMaxSumSumSq(added)
		_, _, outSum, outSumSq := calculations.MinMaxSumSumSq(out)

		b.Avg, b.Varian = calculations.RollingAvgVar(b.Avg, b.Varian,
			b.Sum, b.SumSq,
			addedSum, addedSumSq,
			outSum, outSumSq,
			len(b.Points))

		if slices.Contains(out, b.Min) || slices.Contains(out, b.Max) {
			b.Min, b.Max, b.Sum, b.SumSq = calculations.MinMaxSumSumSq(b.Points)
		} else {
			b.Min = math.Min(b.Min, addedMin)
			b.Max = math.Max(b.Max, addedMax)
			b.Sum += addedSum - outSum
			b.SumSq += addedSumSq - outSumSq
		}

		lastBucket = b
	}

	return d
}
