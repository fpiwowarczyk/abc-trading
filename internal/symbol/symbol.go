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
}

func New(added []float64) *Data {
	buckets := []*Bucket{
		{ // k = 1
			Size: 10,
		},
		{ // k = 2
			Size: 100,
		},
		{ // k = 3
			Size: 1_000,
		},
		{ // k = 4
			Size: 10_000,
		},
		{ // k = 5
			Size: 100_000,
		},
		{ // k = 6
			Size: 1_000_000,
		},
		{ // k = 7
			Size: 10_000_000,
		},
		{ // k = 8
			Size: 100_000_000,
		},
	}

	var lastBucket *Bucket
	var statsCalculatedForWholeSet bool // by default false

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

func (d *Data) Update(added []float64) *Data {
	d.LastPoint = added[len(added)-1]
	d.Points = append(d.Points, added...)

	if len(d.Points) > 100_000_000 {
		d.Points = d.Points[len(d.Points)-100_000_000:]
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
