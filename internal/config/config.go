package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB_URL     string
	JWT_SECRET string
	Port       string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := &Config{
		DB_URL:     os.Getenv("DB_URL"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		Port:       os.Getenv("Port"),
	}
	if cfg.DB_URL == "" || cfg.JWT_SECRET == "" || cfg.Port == "" {
		log.Fatal("Missing required environment variables")
	}
	return cfg
}
