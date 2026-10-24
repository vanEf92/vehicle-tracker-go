package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type LocationPayload struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	// broker := helper.GetEnvRequired("MQTT_BROKER")
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("mock_vehicle_publisher")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Publisher connection failure: %v", token.Error())
	}
	defer client.Disconnect(250)

	// Data mock kendaraan
	vehicleID := "B1234ABC"

	// Data mock lokasi blok m untuk kategori sangat jauh dengan geofence target
	// lat := -6.24457
	// lon := 106.80062

	// Data mock lokasi ancol untuk kategori sangat dekat dengan geofence target
	lat := -6.129121
	lon := 106.833350
	

	log.Printf("Starting mock publisher for vehicle %s every 2 seconds...", vehicleID)

	for {
		payload := LocationPayload{
			VehicleID: vehicleID,
			Latitude:  lat,
			Longitude: lon,
			Timestamp: time.Now().Unix(),
		}

		data, _ := json.Marshal(payload)
		topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)

		token := client.Publish(topic, 0, false, data)
		token.Wait()

		log.Printf("Sent to %s: %s", topic, string(data))

		// Penambahan latitude untuk simulasi pergerakan kendaraan
		lat += 0.0001

		time.Sleep(2 * time.Second)
	}
}