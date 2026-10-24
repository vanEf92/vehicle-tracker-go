package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"vehicles-tracker-go/helper"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Starting DB Migrations...")

	query := `
		CREATE TABLE IF NOT EXISTS vehicle_locations (
			id SERIAL PRIMARY KEY,
			vehicle_id VARCHAR(50) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			timestamp BIGINT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_vehicle_timestamp ON vehicle_locations(vehicle_id, timestamp);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("DB Migrations Done.")
	return nil
}

func main() {
	_ = godotenv.Load()

	dbConnStr := helper.GetEnvRequired("DB_URL_MIGRATE")

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed connection to database: %v", err)
	}
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		log.Fatalf("Failed to database migrate: %v", err)
	}
}