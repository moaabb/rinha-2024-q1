package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/data"
)

type H map[string]string

func main() {
	m := http.NewServeMux()

	m.HandleFunc("GET /clinetes/{id}/extrato", GetRoot)
	m.HandleFunc("POST /clientes/{id}/transacoes", PostRoot)

	srv := http.Server{
		Addr:    ":8080",
		Handler: m,
	}

	log.Println("Server Listening on", srv.Addr)
	srv.ListenAndServe()
}

func GetRoot(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(&H{
		"message": "hello from root",
	})
}

func PostRoot(w http.ResponseWriter, r *http.Request) {

	data.Response(w, http.StatusOK, &data.H{
		"message": "hello from Post",
	})
}
