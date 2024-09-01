package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fpiwowarczyk/abc-trading/internal/config"
	"github.com/fpiwowarczyk/abc-trading/internal/handler"
	"github.com/fpiwowarczyk/abc-trading/internal/transactions"
)

func main() {

	inMemStore := transactions.NewCustomInMemStore()

	cfg := config.NewConfig()

	srv := handler.NewHandler(
		log.New(log.Writer(), "", log.LstdFlags),
		cfg,
		inMemStore,
	)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080), // TODO move it to config
		Handler: srv,
	}

	// TODO Maybe move it to goroutine for grafeful shutdown
	log.Printf("Server is listening on %s", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
