package transactions

import (
	"errors"

	"github.com/fpiwowarczyk/abc-trading/internal/concurrentmap"
	"github.com/fpiwowarczyk/abc-trading/internal/symbol"
)

var _ Store = &InMemStore{}

const (
	MaxBatchSize = 10_000
)

// ImplStore is using my own implementation of concurrent map. It can be compared with NativeInMemStore.
type InMemStore struct {
	symbols *concurrentmap.ConcurrentMap[string, *symbol.Data]
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		symbols: concurrentmap.New[string, *symbol.Data](),
	}
}

func (i *InMemStore) AddBatch(symbolName string, values []float64) error {
	s, exist := i.symbols.Get(symbolName)
	if !exist {
		i.symbols.Set(symbolName, symbol.New(values, 8))
		return nil
	}

	i.symbols.Set(symbolName, s.Update(values))

	return nil
}

func (s *InMemStore) Get(symbolName string) (*symbol.Data, error) {
	sym, ok := s.symbols.Get(symbolName)
	if !ok {
		return nil, errors.New("symbol not found")
	}

	return sym, nil
}
