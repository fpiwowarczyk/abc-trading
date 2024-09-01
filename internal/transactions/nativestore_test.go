package transactions_test

import (
	"testing"

	"github.com/fpiwowarczyk/abc-trading/internal/transactions"
	"github.com/stretchr/testify/assert"
)

func Test_NativeStore(t *testing.T) {

	t.Run("AddBatch", func(t *testing.T) {
		t.Run("should add values to store", func(t *testing.T) {
			// given
			s := transactions.NewNativeInMemStore()

			// when
			err := s.AddBatch("AAPL", []float64{1, 2, 3})

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("should return last k values", func(t *testing.T) {
			// given
			s := transactions.NewNativeInMemStore()
			err := s.AddBatch("AAPL", []float64{1, 2, 3})
			assert.NoError(t, err)

			// when
			values, err := s.Get("AAPL", 1)

			// then
			assert.NoError(t, err)
			assert.Equal(t, []float64{1, 2, 3}, values)
		})

		t.Run("should return newest values when there is more values than specified k", func(t *testing.T) {
			// given
			s := transactions.NewNativeInMemStore()
			err := s.AddBatch("AAPL", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
			assert.NoError(t, err)

			// when
			values, err := s.Get("AAPL", 1)

			// then
			assert.NoError(t, err)
			assert.Equal(t, []float64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, values)
		})

		t.Run("should return nil if symbol does not exist", func(t *testing.T) {
			// given
			s := transactions.NewNativeInMemStore()

			// when
			values, err := s.Get("AAPL", 2)

			// then
			assert.NoError(t, err)
			assert.Nil(t, values)
		})
	})
}
