package calculations

import (
	"math"
)

// RollingAvgAngVariance main purpose is to calculate new variance but part of this calculations is to calculate new average.
func RollingAvgVar(oldAvg, oldVar float64,
	oldSetSum, oldSetSumSq,
	newSetSum, newSetSumSq,
	removedSetSum, removedSetSumSq float64,
	newSetSize int) (avg float64, varian float64) {
	if newSetSize == 0 {
		return 0.0, 0.0
	}

	oldSumAfterRem := oldSetSum - removedSetSum
	oldSumSqAfterRem := oldSetSumSq - removedSetSumSq

	totalSum := oldSumAfterRem + newSetSum
	totalSumSq := oldSumSqAfterRem + newSetSumSq

	avg = totalSum / float64(newSetSize)

	varian = (totalSumSq - math.Pow(totalSum, 2)/float64(newSetSize)) / float64(newSetSize)
	return avg, varian
}

func MinMaxSumSumSq(v []float64) (min, max, sum, sumSq float64) {
	if len(v) == 0 {
		return math.NaN(), math.NaN(), 0.0, 0.0
	}
	min, max = v[0], v[0]
	for _, val := range v {
		sum += val
		sumSq += math.Pow(val, 2)

		min = math.Min(min, val)
		max = math.Max(max, val)
	}

	return min, max, sum, sumSq
}
