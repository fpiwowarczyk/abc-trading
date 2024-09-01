package transactions

import "math"

func lastNValues(values []float64, n int) []float64 {
	return values[len(values)-n:]
}

func getSizeFromK(k int) int {
	return int(math.Pow(10, float64(k)))
}
