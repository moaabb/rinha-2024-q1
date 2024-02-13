package config

import (
	"log"
	"os"
)

type AppConfig struct {
	Port string
	Dsn  string
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		Port: GetEnvOrDie("PORT"),
		Dsn:  GetEnvOrDie("DB_URL"),
	}
}

func GetEnvOrDie(key string) string {
	v := os.Getenv(key)

	if v == "" {
		log.Fatal("Failed to load env")
	}

	return v
}
