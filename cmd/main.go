package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	cfg "github.com/moaabb/rinha-de-backend-2024-q1/internal/config"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/db"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/db/transactiondb"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/handlers"
)

var logger = log.New(os.Stdout, "rinha-app ", log.Ldate|log.Ltime)

func main() {

	cfg := cfg.LoadConfig()

	conn, err := db.Connect(cfg.Dsn, logger, cfg.PoolSize)
	if err != nil {
		log.Fatal(err)
	}

	repo := transactiondb.NewTransactionRepository(conn, logger)
	rh := handlers.NewRinhaHandler(logger, repo)

	app := fiber.New()

	app.Use(LoggingMiddleware)
	app.Post("/clientes/:id/transacoes", rh.CreateTransaction)
	app.Get("/clientes/:id/extrato", rh.GetAccountStatementByPartyId)

	logger.Println("Fiber Listening on port", cfg.Port)
	app.Listen(fmt.Sprintf("0.0.0.0:%s", cfg.Port))
}

func LoggingMiddleware(c *fiber.Ctx) error {
	logger.Println(c.Method(), c.Request().URI())
	return c.Next()
}
