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

	envVars := []string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
	}

	// Explicitly bind environment variables and check for errors
	if err := bindEnvVariables(envVars); err != nil {
		log.Fatalf("Error binding environment variable: %v, it might be empty", err)
	}
}

func bindEnvVariables(vars []string) error {
	for _, v := range vars {
		if err := viper.BindEnv(v); err != nil {
			return err
		}
	}
	return nil
}
