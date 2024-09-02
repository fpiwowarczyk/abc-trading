package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fpiwowarczyk/abc-trading/internal/config"
	"github.com/fpiwowarczyk/abc-trading/internal/transactions"
)

// NewHandler creates a new http.Handler with all routes and middlewares needed for abc-trading service.
func NewHandler(
	logger *log.Logger,
	cfg *config.Config,
	transactionsStore transactions.Store,
) http.Handler {

	mux := http.NewServeMux()
	addRoutes(mux, transactionsStore)

	var handler http.Handler = mux
	handler = logRequest(logger, handler)

	return handler
}

func addRoutes(mux *http.ServeMux, transactionsStore transactions.Store) {
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("POST /add_batch/", handleAddBatch(transactionsStore))
	mux.Handle("GET /stats/", handleStats(transactionsStore))
}

// handleAddBatch allows to add bulk consecutive data points for specific symbol.
//
// Example request:
// curl -X POST -d '{"symbol":"AAPL","values":[1,2,3,4,5,6,7,8,9,0]}' http://localhost:8080/add_batch/
func handleAddBatch(transactionsStore transactions.Store) http.HandlerFunc {
	const maxBatchSize = 10_000

	type request struct {
		Symbol string    `json:"symbol"`
		Values []float64 `json:"values"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		req, err := decode[request](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("decode: %v", err), http.StatusBadRequest)
			return
		}

		if len(req.Values) > maxBatchSize {
			http.Error(w, fmt.Sprintf("batch size is too big, max is %d", maxBatchSize), http.StatusBadRequest)
			return
		}

		if err := transactionsStore.AddBatch(req.Symbol, req.Values); err != nil {
			http.Error(w, fmt.Sprintf("error adding batch values: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// handleStats allows to get statistics for specific symbol based on 10^k last data points.
//
// Example request:
// curl -X GET  http://localhost:8080/stats/?symbol=AAPL&k=3
func handleStats(transactionsStore transactions.Store) http.HandlerFunc {
	type request struct {
		Symbol string `json:"symbol"`
		K      int    `json:"k"`
	}

	type response struct {
		Min  float64 `json:"min"`
		Max  float64 `json:"max"`
		Last float64 `json:"last"`
		Avg  float64 `json:"avg"`
		Var  float64 `json:"var"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		symbolName := params.Get("symbol")
		if symbolName == "" {
			http.Error(w, "symbol is required", http.StatusBadRequest)
			return
		}

		k, err := strconv.Atoi(params.Get("k"))
		if err != nil {
			http.Error(w, "k is required", http.StatusBadRequest)
			return
		}

		if k < 1 || k > 8 {
			http.Error(w, "k should be in the range from 1 to 8, inclusive", http.StatusBadRequest)
			return
		}

		symbol, err := transactionsStore.Get(symbolName)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get: %v", err), http.StatusInternalServerError)
			return
		}

		res := response{
			Last: symbol.LastPoint,
			Min:  symbol.Buckets[k-1].Min,
			Max:  symbol.Buckets[k-1].Max,
			Avg:  symbol.Buckets[k-1].Avg,
			Var:  symbol.Buckets[k-1].Varian,
		}

		log.Printf("%+v \n", symbol.Buckets[k-1])

		if err := encode(w, r, http.StatusOK, res); err != nil {
			http.Error(w, fmt.Sprintf("error encoding: %v", err), http.StatusInternalServerError)
			return
		}

	}
}
