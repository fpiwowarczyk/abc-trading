package transactions_test

import (
	"math/rand"
	"testing"

	"github.com/fpiwowarczyk/abc-trading/internal/transactions"
)

func BenchmarkNativeStore(b *testing.B) {
	nativeStore := transactions.NewNativeInMemStore()

	for range b.N {
		for range 1000 {
			nativeStore.AddBatch("symbol", []float64{rand.Float64(), rand.Float64(), rand.Float64()})
		}

		_, err := nativeStore.Get("symbol", 8)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkCustomStore(b *testing.B) {
	nativeStore := transactions.NewCustomInMemStore()

	for range b.N {
		for range 1000 {
			nativeStore.AddBatch("symbol", []float64{rand.Float64(), rand.Float64(), rand.Float64()})
		}

		_, err := nativeStore.Get("symbol", 8)
		if err != nil {
			panic(err)
		}
	}
}
