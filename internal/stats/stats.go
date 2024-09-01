package stats

import "math"

func FindMax(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	max := values[0]
	for _, v := range values {
		if max != math.Max(max, v) {
			max = v
		}
	}
	return max
}

func FindMin(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	min := values[0]
	for _, v := range values {
		if min != math.Min(min, v) {
			min = v
		}
	}
	return min
}

func FindLast(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	return values[len(values)-1]
}

func FindAvg(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func FindVar(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}

	avg := FindAvg(values)
	sum := 0.0
	for _, v := range values {
		sum += (v - avg) * (v - avg)
	}
	return sum / float64(len(values))
}
