package transactions

import (
	"sync"
)

var _ Store = &NativeInMemStore{}

// NativeInMemStore is a simple in-memory store that stores all data using native golang data structures.
// It is used mostry for comparison with other stores.
type NativeInMemStore struct {
	fields sync.Map
}

func NewNativeInMemStore() *NativeInMemStore {
	return &NativeInMemStore{
		fields: sync.Map{},
	}
}

func (s *NativeInMemStore) AddBatch(symbol string, values []float64) error {
	v, exist := s.fields.Load(symbol)
	if !exist {
		s.fields.Store(symbol, values)
		return nil
	}

	v = append(v.([]float64), values...)
	s.fields.Store(symbol, v)

	return nil
}

func (s *NativeInMemStore) Get(symbol string, k int) ([]float64, error) {
	v, ok := s.fields.Load(symbol)
	if !ok {
		return nil, nil
	}

	size := getSizeFromK(k)

	values := v.([]float64)
	if size > len(values) {
		return values, nil
	}

	return lastNValues(values, size), nil
}
