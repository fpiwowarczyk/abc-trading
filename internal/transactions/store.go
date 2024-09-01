package transactions

// Store in an interface that is used to store and retrieve data for abc-tranding service.
type Store interface {
	AddBatch(symbol string, values []float64) error
	Get(symbol string, k int) ([]float64, error)
}
