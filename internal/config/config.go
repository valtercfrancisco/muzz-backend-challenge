package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

func InitConfig() {
	if err := godotenv.Load("./db-variables.env"); err != nil {
		log.Fatalf("Error loading db-variables.env file: %v", err)
	}

	viper.AutomaticEnv()

	// Explicitly bind environment variables
	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_PORT")
	viper.BindEnv("POSTGRES_USER")
	viper.BindEnv("POSTGRES_PASSWORD")
	viper.BindEnv("POSTGRES_DB")
}
