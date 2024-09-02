package transactions

import (
	"errors"

	"github.com/fpiwowarczyk/abc-trading/internal/concurrentmap"
	"github.com/fpiwowarczyk/abc-trading/internal/symbol"
)

var _ Store = &InMemStore{}

// ImplStore is using my own implementation of concurrent map. It should be thread safe.
type InMemStore struct {
	symbols *concurrentmap.ConcurrentMap[string, *symbol.Data]
	// MaxK is the number of buckets to store in the symbol each bucket is 10 times bigger than the previous one.
	MaxK int
}

// NewInMemStore creates a new instance of InMemStore.
func NewInMemStore(maxK int) *InMemStore {
	return &InMemStore{
		symbols: concurrentmap.New[string, *symbol.Data](),
		MaxK:    maxK,
	}
}

// AddBatch adds a new batch of values to the symbol. If the symbol does not exist it will be created.
func (i *InMemStore) AddBatch(symbolName string, values []float64) error {
	s, exist := i.symbols.Get(symbolName)
	if !exist {
		i.symbols.Set(symbolName, symbol.New(values, i.MaxK))
		return nil
	}

	i.symbols.Set(symbolName, s.Update(values))

	return nil
}

// Get returns the symbol data for the given symbol name.
func (s *InMemStore) Get(symbolName string) (*symbol.Data, error) {
	sym, ok := s.symbols.Get(symbolName)
	if !ok {
		return nil, errors.New("symbol not found")
	}

	return sym, nil
}
