package transactions

import "github.com/fpiwowarczyk/abc-trading/internal/concurrentmap"

var _ Store = &CustomInMemStore{}

// ImplStore is using my own implementation of concurrent map. It can be compared with NativeInMemStore.
type CustomInMemStore struct {
	fields *concurrentmap.ConcurrentMap[string, []float64]
}

func NewCustomInMemStore() *CustomInMemStore {
	return &CustomInMemStore{
		fields: concurrentmap.New[string, []float64](),
	}
}

func (s *CustomInMemStore) AddBatch(symbol string, values []float64) error {
	v, exist := s.fields.Get(symbol)
	if !exist {
		s.fields.Set(symbol, values)
		return nil
	}

	v = append(v, values...)
	s.fields.Set(symbol, v)

	return nil
}

func (s *CustomInMemStore) Get(symbol string, k int) ([]float64, error) {
	v, ok := s.fields.Get(symbol)
	if !ok {
		return nil, nil
	}

	size := getSizeFromK(k)

	values := v
	if size > len(values) {
		return values, nil
	}

	return lastNValues(values, size), nil
}
