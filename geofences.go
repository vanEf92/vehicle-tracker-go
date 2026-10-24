package main

import (
	"math"
)

type GeofenceTarget struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

var GeofenceTargets = []GeofenceTarget{
	{
		ID:        "ZONE-01",
		Name:      "Grand Indonesia",
		Latitude:  -6.1956278,
		Longitude: 106.8218750,
		Radius:    50.0,
	},
	{
		ID:        "ZONE-02",
		Name:      "Monumen Nasional",
		Latitude:  -6.175388,
		Longitude: 106.827139,
		Radius:    50.0,
	},
	{
		ID:        "ZONE-03",
		Name:      "Ancol",
		Latitude:  -6.129121,
		Longitude: 106.833350,
		Radius:    50.0,
	},
}

func CheckGeofence(lat, lon float64) (bool, *GeofenceTarget) {
	for _, target := range GeofenceTargets {
		dist := calculateDistance(lat, lon, target.Latitude, target.Longitude)
		if dist <= target.Radius {
			return true, &target
		}
	}
	return false, nil
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}