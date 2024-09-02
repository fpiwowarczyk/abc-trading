package symbol

import (
	"math"
	"slices"

	"github.com/fpiwowarczyk/abc-trading/internal/calculations"
)

type Data struct {
	LastPoint float64
	Buckets   []*Bucket

	Points []float64

	MaxK int
}

// New creates a new Data strucutre that represends data stored for single symbol.
// During this operation data is divided into buckets, which are used to store last n points.
// Number of buckets is specified by MaxK.
func New(added []float64, maxK int) *Data {
	var buckets []*Bucket
	for k := range maxK {
		buckets = append(buckets, &Bucket{
			Size:   int(math.Pow(10, float64(k+1))),
			Points: make([]float64, 0),
		})
	}

	var lastBucket *Bucket
	var statsCalculatedForWholeSet bool

	for _, bucket := range buckets {
		if statsCalculatedForWholeSet {
			bucket.SetStats(
				lastBucket.Min, lastBucket.Max,
				lastBucket.Avg, lastBucket.Varian,
				lastBucket.Sum, lastBucket.SumSq,
				lastBucket.Points)

			continue
		}

		if bucket.CantFitIntoBucket(len(added)) {
			bucket.Points = added[len(added)-bucket.Size:]
		} else {
			bucket.Points = added
			statsCalculatedForWholeSet = true
		}

		bucket.Min, bucket.Max, bucket.Sum, bucket.SumSq = calculations.MinMaxSumSumSq(bucket.Points)
		bucket.Avg, bucket.Varian = calculations.RollingAvgVar(
			0.0, 0.0,
			0.0, 0.0,
			bucket.Sum, bucket.SumSq,
			0.0, 0.0,
			len(bucket.Points))

		lastBucket = bucket
	}

	return &Data{
		LastPoint: added[len(added)-1],
		Buckets:   buckets,
		Points:    added,
		MaxK:      maxK,
	}
}

// Update appends new data points to the existing data set.
// It also updates statistics for each bucket.
// If number of data points exceeds then maximal number of points that can be stored for single symbol,
// then the oldest points are removed.
func (d *Data) Update(added []float64) *Data {
	d.LastPoint = added[len(added)-1]
	d.Points = append(d.Points, added...)

	if len(d.Points) > int(math.Pow(10, float64(d.MaxK))) {
		d.Points = d.Points[len(d.Points)-int(math.Pow(10, float64(d.MaxK))):]
	}

	var lastBucket *Bucket
	var calculatedForWholeSet bool

	for _, bucket := range d.Buckets {
		if calculatedForWholeSet {
			bucket.SetStats(
				lastBucket.Min, lastBucket.Max,
				lastBucket.Avg, lastBucket.Varian,
				lastBucket.Sum, lastBucket.SumSq,
				lastBucket.Points)

			continue
		}

		var removed []float64
		if bucket.CantFitIntoBucket(len(added)) {
			removed = bucket.Points[:len(bucket.Points)+len(added)-bucket.Size]
			bucket.Points = d.Points[len(d.Points)-bucket.Size:]
		} else {
			calculatedForWholeSet = true
			bucket.Points = append(bucket.Points, added...)
		}

		addedMin, addedMax, addedSum, addedSumSq := calculations.MinMaxSumSumSq(added)
		_, _, outSum, outSumSq := calculations.MinMaxSumSumSq(removed)

		bucket.Avg, bucket.Varian = calculations.RollingAvgVar(bucket.Avg, bucket.Varian,
			bucket.Sum, bucket.SumSq,
			addedSum, addedSumSq,
			outSum, outSumSq,
			len(bucket.Points))

		if slices.Contains(removed, bucket.Min) || slices.Contains(removed, bucket.Max) {
			bucket.Min, bucket.Max, bucket.Sum, bucket.SumSq = calculations.MinMaxSumSumSq(bucket.Points)
		} else {
			bucket.Min = math.Min(bucket.Min, addedMin)
			bucket.Max = math.Max(bucket.Max, addedMax)
			bucket.Sum += addedSum - outSum
			bucket.SumSq += addedSumSq - outSumSq
		}

		lastBucket = bucket
	}

	return d
}
