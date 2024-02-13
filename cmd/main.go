package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/config"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/db"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/db/transactiondb"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/handlers"
)

type H map[string]string

func main() {

	cfg := config.LoadConfig()

	logger := log.New(os.Stdout, "rinha-app", log.Ldate|log.Ltime)

	conn, err := db.Connect(cfg.Dsn, logger)
	if err != nil {
		log.Fatal(err)
	}

	repo := transactiondb.NewTransactionRepository(conn, logger)
	rh := handlers.NewRinhaHandler(logger, repo)

	m := http.NewServeMux()

	m.HandleFunc("GET /clientes/{id}/extrato", rh.GetAccountStatementByPartyId)
	m.HandleFunc("POST /clientes/{id}/transacoes", rh.CreateTransaction)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: m,
	}

	log.Println("Server Listening on port", cfg.Port)
	srv.ListenAndServe()
}
