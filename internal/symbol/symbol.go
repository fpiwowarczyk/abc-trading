package symbol

import (
	"log"
	"math"
	"slices"

	"github.com/fpiwowarczyk/abc-trading/internal/calculations"
)

type Data struct {
	Last    float64
	Buckets []*Bucket

	Points []float64
}

type Bucket struct {
	Size   int
	Min    float64
	Max    float64
	Avg    float64
	Varian float64
	Points []float64

	Sum   float64
	SumSq float64
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

	var lastMin, lastMax, lastSum, lastSumSq, lastAvg, lastVarian float64 // by default 0.0
	var calculatedForWholeSet bool                                        // by default false

	for _, b := range buckets {
		if calculatedForWholeSet {
			b.Min, b.Max, b.Sum, b.SumSq, b.Avg, b.Varian, b.Points = lastMin, lastMax, lastSum, lastSumSq, lastAvg, lastVarian, added
			log.Println("Bucket", b.Size, "Min", b.Min, "Max", b.Max, "Avg", b.Avg, "Varian", b.Varian, "Points", b.Points)
			continue
		}

		var currentWindow []float64
		if b.Size < len(added) {
			currentWindow = added[len(added)-b.Size:]
		} else {
			currentWindow = added
			calculatedForWholeSet = true
		}

		b.Points = currentWindow
		b.Min, b.Max, b.Sum, b.SumSq = calculations.MinMaxSumSumSq(currentWindow)
		b.Avg, b.Varian = calculations.RollingAvgAndVar(
			0.0, 0.0,
			0.0, 0.0,
			b.Sum, b.SumSq,
			0.0, 0.0,
			len(currentWindow))

		lastMin, lastMax, lastSum, lastSumSq, lastAvg, lastVarian = b.Min, b.Max, b.Sum, b.SumSq, b.Avg, b.Varian

		log.Println("Bucket", b.Size, "Min", b.Min, "Max", b.Max, "Avg", b.Avg, "Varian", b.Varian, "Points", b.Points)

	}

	return &Data{
		Last:    added[len(added)-1],
		Buckets: buckets,
		Points:  added,
	}
}

func (d *Data) Update(added []float64) *Data {
	d.Last = added[len(added)-1]
	d.Points = append(d.Points, added...)

	if len(d.Points) > 100_000_000 {
		d.Points = d.Points[len(d.Points)-100_000_000:]
	}

	var lastMin, lastMax, lastSum, lastSumSq, lastAvg, lastVarian float64 // by default 0.0
	var calculatedForWholeSet bool                                        // by default false

	for _, b := range d.Buckets {
		if calculatedForWholeSet {
			b.Min, b.Max, b.Sum, b.SumSq, b.Avg, b.Varian = lastMin, lastMax, lastSum, lastSumSq, lastAvg, lastVarian
			continue
		}

		var out []float64
		if len(b.Points)+len(added) > b.Size {
			out = b.Points[:len(b.Points)+len(added)-b.Size]
		}

		addedMin, addedMax, addedSum, addedSumSq := calculations.MinMaxSumSumSq(added)
		_, _, outSum, outSumSq := calculations.MinMaxSumSumSq(out)

		b.Avg, b.Varian = calculations.RollingAvgAndVar(b.Avg, b.Varian,
			b.Sum, b.SumSq,
			addedSum, addedSumSq,
			outSum, outSumSq,
			len(b.Points)+len(added)-len(out))

		if len(b.Points)+len(added) > b.Size {
			b.Points = d.Points[len(d.Points)-b.Size:]
		} else {
			b.Points = append(b.Points, added...)
		}

		if slices.Contains(out, b.Min) || slices.Contains(out, b.Max) {
			b.Min, b.Max, b.Sum, b.SumSq = calculations.MinMaxSumSumSq(b.Points)
		} else {
			b.Min = math.Min(b.Min, addedMin)
			b.Max = math.Max(b.Max, addedMax)
			b.Sum += addedSum - outSum
			b.SumSq += addedSumSq - outSumSq
		}

		log.Println("Bucket", b.Size, "Min", b.Min, "Max", b.Max, "Avg", b.Avg, "Varian", b.Varian, "Points", b.Points)

	}

	return d
}
