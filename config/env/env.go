package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	// Logging error if .env file not found or cannot be loaded
	if err!=nil {
		fmt.Println("Error loading .env file")
	}
}

func GetString(key string, fallback string) string{
	value,ok:=os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}