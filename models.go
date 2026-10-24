package main

import (
	"database/sql"
)

type LocationPayload struct {
	VehicleID string  `json:"vehicle_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Timestamp int64   `json:"timestamp" binding:"required"`
}

type GeofenceEvent struct {
	VehicleID string `json:"vehicle_id"`
	Event     string `json:"event"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Timestamp int64 `json:"timestamp"`
}

type VehicleRepository struct {
	DB *sql.DB
}

func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{DB: db}
}

func (r *VehicleRepository) Save(p *LocationPayload) error {
	query := "INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) VALUES ($1, $2, $3, $4)"
	_, err := r.DB.Exec(query, p.VehicleID, p.Latitude, p.Longitude, p.Timestamp)
	return err
}

func (r *VehicleRepository) GetLatest(vehicleID string) (*LocationPayload, error) {
	query := "SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 ORDER BY timestamp DESC LIMIT 1"
	var p LocationPayload
	err := r.DB.QueryRow(query, vehicleID).Scan(&p.VehicleID, &p.Latitude, &p.Longitude, &p.Timestamp)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *VehicleRepository) GetHistory(vehicleID string, start, end int64) ([]LocationPayload, error) {
	query := "SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3 ORDER BY timestamp ASC"
	rows, err := r.DB.Query(query, vehicleID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []LocationPayload
	for rows.Next() {
		var p LocationPayload
		if err := rows.Scan(&p.VehicleID, &p.Latitude, &p.Longitude, &p.Timestamp); err == nil {
			history = append(history, p)
		}
	}
	return history, nil
}