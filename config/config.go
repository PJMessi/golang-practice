package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APP_PORT            string
	DB_HOST             string
	DB_PORT             string
	DB_DATABASE         string
	DB_USER             string
	DB_PASSWORD         string
	JWT_SECRET          string
	JWT_EXPIRATION_TIME string
	NATS_URL            string
}

func GetAppConfig(env string) *AppConfig {
	envPath := getEnvFilePath(env)
	if doesEnvFileExist(envPath) {
		loadEnvFile(envPath)
	}

	return &AppConfig{
		APP_PORT:            os.Getenv("APP_PORT"),
		DB_HOST:             os.Getenv("DB_HOST"),
		DB_PORT:             os.Getenv("DB_PORT"),
		DB_DATABASE:         os.Getenv("DB_DATABASE"),
		DB_USER:             os.Getenv("DB_USER"),
		DB_PASSWORD:         os.Getenv("DB_PASSWORD"),
		JWT_SECRET:          os.Getenv("JWT_SECRET"),
		JWT_EXPIRATION_TIME: os.Getenv("JWT_EXPIRATION_TIME"),
		NATS_URL:            os.Getenv("NATS_URL"),
	}
}

func doesEnvFileExist(envFilePath string) bool {
	_, err := os.Stat(envFilePath)
	return err == nil
}

func getEnvFilePath(env string) string {
	if env == "" {
		return ".env"
	}

	if env == "test" {
		return "../.env"
	}

	return ""
}

func loadEnvFile(envFilePath string) {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalf("error while loading env file: %v\n", err)
	}
}
