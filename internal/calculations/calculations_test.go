package calculations_test

import (
	"math"
	"testing"

	"github.com/fpiwowarczyk/abc-trading/internal/calculations"
	"github.com/stretchr/testify/assert"
)

func Test_MinMaxSumSumSq(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		input         []float64
		expectedMin   float64
		expectedMax   float64
		expectedSum   float64
		expectedSumSq float64
		expectNaN     bool
	}{
		{
			name:          "empty input",
			input:         []float64{},
			expectedMin:   math.NaN(),
			expectedMax:   math.NaN(),
			expectedSum:   0.0,
			expectedSumSq: 0.0,
			expectNaN:     true,
		},
		{
			name:          "single element",
			input:         []float64{1.0},
			expectedMin:   1.0,
			expectedMax:   1.0,
			expectedSum:   1.0,
			expectedSumSq: 1.0,
		},
		{
			name:          "multiple elements",
			input:         []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			expectedMin:   1.0,
			expectedMax:   5.0,
			expectedSum:   15.0,
			expectedSumSq: 55.0,
		},
		{
			name:          "negative elements",
			input:         []float64{-1.0, -2.0, -3.0, -4.0, -5.0},
			expectedMin:   -5.0,
			expectedMax:   -1.0,
			expectedSum:   -15.0,
			expectedSumSq: 55.0,
		},
		{
			name:          "mixed elements",
			input:         []float64{-1.0, 2.0, -3.0, 4.0, -5.0},
			expectedMin:   -5.0,
			expectedMax:   4.0,
			expectedSum:   -3.0,
			expectedSumSq: 55.0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			min, max, sum, sumSq := calculations.MinMaxSumSumSq(tc.input)

			if tc.expectNaN {
				assert.True(t, math.IsNaN(min))
				assert.True(t, math.IsNaN(max))
			} else {
				assert.Equal(t, tc.expectedMin, min)
				assert.Equal(t, tc.expectedMax, max)
				assert.Equal(t, tc.expectedSum, sum)
				assert.Equal(t, tc.expectedSumSq, sumSq)
			}
		})

	}
}

func Test_RollingAvgVar(t *testing.T) {

	testCases := []struct {
		name            string
		oldAvg          float64
		oldVar          float64
		oldSetSum       float64
		oldSetSumSq     float64
		newSetSum       float64
		newSetSumSq     float64
		removedSetSum   float64
		removedSetSumSq float64
		newSetSize      int
		expectedAvg     float64
		expectedVarian  float64
	}{
		{
			name:            "empty input",
			oldAvg:          0.0,
			oldVar:          0.0,
			oldSetSum:       0.0,
			oldSetSumSq:     0.0,
			newSetSum:       0.0,
			newSetSumSq:     0.0,
			removedSetSum:   0.0,
			removedSetSumSq: 0.0,
			newSetSize:      0,
			expectedAvg:     0.0,
			expectedVarian:  0.0,
		},
		{
			name:            "single element",
			oldAvg:          0.0,
			oldVar:          0.0,
			oldSetSum:       0.0,
			oldSetSumSq:     0.0,
			newSetSum:       1.0,
			newSetSumSq:     1.0,
			removedSetSum:   0.0,
			removedSetSumSq: 0.0,
			newSetSize:      1,
			expectedAvg:     1.0,
			expectedVarian:  0.0,
		},
		{
			name:            "multiple elements",
			oldAvg:          0.0,
			oldVar:          0.0,
			oldSetSum:       0.0,
			oldSetSumSq:     0.0,
			newSetSum:       15.0,
			newSetSumSq:     55.0,
			removedSetSum:   0.0,
			removedSetSumSq: 0.0,
			newSetSize:      5,
			expectedAvg:     3.0,
			expectedVarian:  2.0,
		},
		{
			name:            "negative elements",
			oldAvg:          0.0,
			oldVar:          0.0,
			oldSetSum:       0.0,
			oldSetSumSq:     0.0,
			newSetSum:       -15.0,
			newSetSumSq:     55.0,
			removedSetSum:   0.0,
			removedSetSumSq: 0.0,
			newSetSize:      5,
			expectedAvg:     -3.0,
			expectedVarian:  2.0,
		},
		{
			name:            "mixed elements",
			oldAvg:          0.0,
			oldVar:          0.0,
			oldSetSum:       0.0,
			oldSetSumSq:     0.0,
			newSetSum:       -3.0,
			newSetSumSq:     55.0,
			removedSetSum:   0.0,
			removedSetSumSq: 0.0,
			newSetSize:      5,
			expectedAvg:     -0.6,
			expectedVarian:  10.64,
		},
		{
			name:            "remove elements",
			oldAvg:          3.0,
			oldVar:          2.0,
			oldSetSum:       15.0,
			oldSetSumSq:     55.0,
			newSetSum:       0.0,
			newSetSumSq:     0.0,
			removedSetSum:   15.0,
			removedSetSumSq: 55.0,
			newSetSize:      0,
			expectedAvg:     0.0,
			expectedVarian:  0.0,
		},
		{
			name:            "remove elements and add new",
			oldAvg:          3.0,
			oldVar:          2.0,
			oldSetSum:       15.0,
			oldSetSumSq:     55.0,
			newSetSum:       15.0,
			newSetSumSq:     55.0,
			removedSetSum:   15.0,
			removedSetSumSq: 55.0,
			newSetSize:      5,
			expectedAvg:     3.0,
			expectedVarian:  2.0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			avg, varian := calculations.RollingAvgVar(
				tc.oldAvg, tc.oldVar,
				tc.oldSetSum, tc.oldSetSumSq,
				tc.newSetSum, tc.newSetSumSq,
				tc.removedSetSum, tc.removedSetSumSq,
				tc.newSetSize)

			assert.Equal(t, tc.expectedAvg, avg)
			assert.Equal(t, tc.expectedVarian, varian)
		})
	}

}
