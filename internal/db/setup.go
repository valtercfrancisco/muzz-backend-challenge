package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// ConnectDB connects to the PostgreSQL database using configuration from environment variables.
func ConnectDB() (*sql.DB, error) {
	dbHost := viper.GetString("POSTGRES_HOST")
	dbPort := viper.GetString("POSTGRES_PORT")
	dbUser := viper.GetString("POSTGRES_USER")
	dbPassword := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DB")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Printf("Connection URL: %s", connectionString)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	return db, nil
}

// RunMigrations runs the database migrations.
func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///db/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations applied successfully!")
	return nil
}
