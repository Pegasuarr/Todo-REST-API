package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db        string `yaml:"db"`
	Port      string `yaml:"port"`
	JwtSecret string `yaml:"jwt_secret"`
}

func LoadConfig() (*Config, error) {
	var err error = godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var config *Config = &Config{
		Db:        os.Getenv("DB_URL"),
		Port:      os.Getenv("PORT"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
	return config, nil
}
