package transactions

import (
	"errors"

	"github.com/fpiwowarczyk/abc-trading/internal/concurrentmap"
	"github.com/fpiwowarczyk/abc-trading/internal/symbol"
)

var _ Store = &InMemStore{}

// ImplStore is using my own implementation of concurrent map. It can be compared with NativeInMemStore.
type InMemStore struct {
	symbols *concurrentmap.ConcurrentMap[string, *symbol.Data]
	MaxK    int
}

func NewInMemStore(maxK int) *InMemStore {
	return &InMemStore{
		symbols: concurrentmap.New[string, *symbol.Data](),
		MaxK:    maxK,
	}
}

func (i *InMemStore) AddBatch(symbolName string, values []float64) error {
	s, exist := i.symbols.Get(symbolName)
	if !exist {
		i.symbols.Set(symbolName, symbol.New(values, i.MaxK))
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
