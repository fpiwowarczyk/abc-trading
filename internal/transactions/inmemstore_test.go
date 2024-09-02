package transactions_test

import (
	"testing"

	"github.com/fpiwowarczyk/abc-trading/internal/transactions"
	"github.com/stretchr/testify/assert"
)

func Test_CustomStore(t *testing.T) {
	t.Run("AddBatch", func(t *testing.T) {
		store := transactions.NewInMemStore()
		err := store.AddBatch("test", []float64{1, 2, 3, 4, 5})
		assert.Nil(t, err)

		stats, err := store.Get("test")
		assert.Nil(t, err)

		assert.Equal(t, 5.0, stats.LastPoint)
		assert.Len(t, stats.Buckets, 8)
		assert.Len(t, stats.Points, 5)
		assert.Equal(t, []float64{1, 2, 3, 4, 5}, stats.Buckets[0].Points)
		assert.Equal(t, 1.0, stats.Buckets[0].Min)
		assert.Equal(t, 5.0, stats.Buckets[0].Max)
		assert.Equal(t, 15.0, stats.Buckets[0].Sum)
		assert.Equal(t, 55.0, stats.Buckets[0].SumSq)
		assert.Equal(t, 3.0, stats.Buckets[0].Avg)
		assert.Equal(t, 2.0, stats.Buckets[0].Varian)
	})

}
