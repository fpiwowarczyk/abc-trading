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
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	inMemStore := transactions.NewInMemStore(cfg.MaxK)

	srv := handler.NewHandler(
		log.New(log.Writer(), "", log.LstdFlags),
		cfg,
		inMemStore,
	)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: srv,
	}

	log.Printf("Server is listening on %s", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
