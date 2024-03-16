package config

import (
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	Port     string
	Dsn      string
	PoolSize int
}

func LoadConfig() *AppConfig {
	pool, _ := strconv.Atoi(GetEnvOrDie("POOL_SIZE"))

	return &AppConfig{
		Port:     GetEnvOrDie("PORT"),
		Dsn:      GetEnvOrDie("DB_URL"),
		PoolSize: pool,
	}
}

func GetEnvOrDie(key string) string {
	v := os.Getenv(key)

	if v == "" {
		log.Fatal("Failed to load env")
	}

	return v
}
