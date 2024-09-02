package symbol_test

import (
	"testing"

	"github.com/fpiwowarczyk/abc-trading/internal/symbol"
	"github.com/stretchr/testify/assert"
)

func Test_CreatesSymbol(t *testing.T) {
	t.Parallel()

	s := symbol.New([]float64{1, 2, 3, 4, 5}, 8)
	assert.NotNil(t, s)
	assert.Equal(t, 5.0, s.LastPoint)
	assert.Len(t, s.Buckets, 8)
	assert.Len(t, s.Points, 5)
	assert.Equal(t, []float64{1, 2, 3, 4, 5}, s.Buckets[0].Points)
	assert.Equal(t, 1.0, s.Buckets[0].Min)
	assert.Equal(t, 5.0, s.Buckets[0].Max)
	assert.Equal(t, 15.0, s.Buckets[0].Sum)
	assert.Equal(t, 55.0, s.Buckets[0].SumSq)
	assert.Equal(t, 3.0, s.Buckets[0].Avg)
	assert.Equal(t, 2.0, s.Buckets[0].Varian)
}

func Test_UpdateSymbol(t *testing.T) {
	t.Run("points fit into bucket", func(t *testing.T) {
		s := symbol.New([]float64{1, 2, 3, 4, 5}, 8)
		s.Update([]float64{6, 7, 8, 9, 10})
		assert.Equal(t, 10.0, s.LastPoint)
		assert.Equal(t, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s.Buckets[0].Points)
		assert.Equal(t, 1.0, s.Buckets[0].Min)
		assert.Equal(t, 10.0, s.Buckets[0].Max)
		assert.Equal(t, 55.0, s.Buckets[0].Sum)
		assert.Equal(t, 385.0, s.Buckets[0].SumSq)
		assert.Equal(t, 5.5, s.Buckets[0].Avg)
		assert.Equal(t, 8.25, s.Buckets[0].Varian)
	})

	t.Run("points do not fit into bucket", func(t *testing.T) {
		s := symbol.New([]float64{1, 2, 3, 4, 5}, 8)
		s.Update([]float64{6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
		assert.Equal(t, 15.0, s.LastPoint)
		assert.Equal(t, []float64{6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, s.Buckets[0].Points)
		assert.Equal(t, 6.0, s.Buckets[0].Min)
		assert.Equal(t, 15.0, s.Buckets[0].Max)
		assert.Equal(t, 105.0, s.Buckets[0].Sum)
		assert.Equal(t, 1185.0, s.Buckets[0].SumSq)
		assert.Equal(t, 10.5, s.Buckets[0].Avg)
		assert.Equal(t, 8.25, s.Buckets[0].Varian)
	})

	t.Run("removes points that past max size correctly", func(t *testing.T) {
		s := symbol.New([]float64{1, 2, 3, 4, 5}, 1)
		assert.Equal(t, 5, len(s.Points))
		assert.Equal(t, 1, len(s.Buckets))

		s.Update([]float64{6, 7, 8, 9, 10})
		assert.Equal(t, 10, len(s.Points))

		s.Update([]float64{11, 12, 13, 14, 15})
		assert.Equal(t, 10, len(s.Points))
		assert.Equal(t, []float64{6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, s.Points)

	})
}
