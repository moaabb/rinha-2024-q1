package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultMaxConnLifetime = time.Hour
const defaultMaxConnIdleTime = time.Second * 300

func Connect(dsn string, logger *log.Logger, poolsize int) (*pgxpool.Pool, error) {

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Println("[ERROR] could not connect to db", err)
		return nil, err
	}

	cfg.MaxConns = int32(poolsize)
	cfg.MaxConnLifetime = defaultMaxConnLifetime
	cfg.MaxConnIdleTime = defaultMaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		logger.Println("[ERROR] could not connect to db", err)
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("Connected to db!")

	return pool, nil
}
