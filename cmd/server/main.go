package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dionisvl/my.web3.auth/internal/auth"
	"github.com/dionisvl/my.web3.auth/internal/config"
	"github.com/dionisvl/my.web3.auth/internal/handlers"
	"github.com/dionisvl/my.web3.auth/internal/wallet"
)

func main() {
	cfg := config.Load()

	authSvc := auth.New(cfg.SessionKey)
	walletSvc := wallet.New(cfg.EthNetwork)

	h, err := handlers.New(authSvc, walletSvc)
	if err != nil {
		log.Fatalf("failed to init handlers: %v", err)
	}

	mux := http.NewServeMux()
	if err := h.Register(mux); err != nil {
		log.Fatalf("failed to register routes: %v", err)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("web3auth listening on :%s (network=%s)", cfg.Port, cfg.EthNetwork)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
