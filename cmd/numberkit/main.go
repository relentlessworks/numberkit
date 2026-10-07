package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/relentlessworks/numberkit/internal/api"
	"github.com/relentlessworks/numberkit/internal/config"
)

func main() {
	cfg := config.Load()

	handler := api.NewHandler(cfg.NoAuth)
	mux := http.NewServeMux()
	handler.Register(mux)

	log.Printf("numberkit starting on %s (no-auth=%v)", cfg.Addr, cfg.NoAuth)
	fmt.Printf("numberkit — agentic-first number formatting and conversion service\n")
	fmt.Printf("Listening on %s\n", cfg.Addr)
	fmt.Printf("GET /help for the operating manual\n")

	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
