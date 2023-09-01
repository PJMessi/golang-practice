package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB_HOST     string
	DB_PORT     string
	DB_DATABASE string
	DB_USER     string
	DB_PASSWORD string
}

func GetAppConfig() *AppConfig {
	envPath := getEnvFilePath()
	if doesEnvFileExist(envPath) {
		loadEnvFile(envPath)
	}

	return &AppConfig{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_DATABASE: os.Getenv("DB_DATABASE"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	}
}

func doesEnvFileExist(envFilePath string) bool {
	_, err := os.Stat(envFilePath)
	return err == nil
}

func getEnvFilePath() string {
	return ".env"
}

func loadEnvFile(envFilePath string) {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalf("error while loading env file: %v\n", err)
	}
}
