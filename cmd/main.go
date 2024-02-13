package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/db"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/db/transactiondb"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/handlers"
)

type H map[string]string

func main() {

	logger := log.New(os.Stdout, "rinha-app", log.Ldate|log.Ltime)

	conn, err := db.Connect("postgres://moab:supersecure@localhost:5432/rinhadb", logger)
	if err != nil {
		log.Fatal(err)
	}

	repo := transactiondb.NewTransactionRepository(conn, logger)
	rh := handlers.NewRinhaHandler(logger, repo)

	m := http.NewServeMux()

	m.HandleFunc("GET /clientes/{id}/extrato", rh.GetAccountStatementByPartyId)
	m.HandleFunc("POST /clientes/{id}/transacoes", rh.CreateTransaction)

	srv := http.Server{
		Addr:    ":8080",
		Handler: m,
	}

	log.Println("Server Listening on", srv.Addr)
	srv.ListenAndServe()
}
