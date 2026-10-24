package main

import (
	"database/sql"
	"encoding/json"
	"log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"vehicles-tracker-go/helper"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	_ = godotenv.Load()

	dbConnStr := helper.GetEnvRequired("DB_URL")
	mqttBroker := helper.GetEnvRequired("MQTT_BROKER")
	rabbitURL := helper.GetEnvRequired("RABBITMQ_URL")
	serverPort := helper.GetEnvRequired("SERVER_PORT")

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed connection to database: %v", err)
	}
	defer db.Close()

	repo := NewVehicleRepository(db)
	log.Println("Database ready!.")

	rbConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed connection to RabbitMQ: %v", err)
	}
	defer rbConn.Close()

	rbCh, err := rbConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel RabbitMQ: %v", err)
	}
	defer rbCh.Close()

	_ = rbCh.ExchangeDeclare("fleet.events", "direct", true, false, false, false, nil)
	q, _ := rbCh.QueueDeclare("geofence_alerts", true, false, false, false, nil)
	_ = rbCh.QueueBind(q.Name, "geofence.entry", "fleet.events", false, nil)

	go startRabbitMQWorker(rbConn)

	mqttOpts := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID("backend_mqtt_subscriber")
	mqttClient := mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed connection to MQTT: %v", token.Error())
	}
	defer mqttClient.Disconnect(250)

	mqttClient.Subscribe("/fleet/vehicle/+/location", 0, func(client mqtt.Client, msg mqtt.Message) {
		var payload LocationPayload
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Printf("Invalid Format JSON MQTT: %v", err)
			return
		}

		if err := repo.Save(&payload); err != nil {
			log.Printf("Failed saving into DB: %v", err)
			return
		}
		log.Printf("[MQTT] Location saved for vehicle: %s", payload.VehicleID)

		matched, target := CheckGeofence(payload.Latitude, payload.Longitude)
		if matched {
			eventMsg := GeofenceEvent{
				VehicleID: payload.VehicleID,
				Event:     "geofence_entry",
				Timestamp: payload.Timestamp,
			}
			eventMsg.Location.Latitude = payload.Latitude
			eventMsg.Location.Longitude = payload.Longitude

			body, _ := json.Marshal(eventMsg)
			_ = rbCh.Publish("fleet.events", "geofence.entry", false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
			log.Printf("[GEOFENCE] Vehicle %s enter [%s] (%s)!", 
				payload.VehicleID, target.ID, target.Name)
		}
	})
	log.Println("MQTT Subscriber active listening: /fleet/vehicle/+/location")

	r := gin.Default()

	RegisterRoutes(r, repo)

	log.Printf("REST API is running at port %s", serverPort)
	r.Run(serverPort)
}

func startRabbitMQWorker(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Worker failed to open channel: %v", err)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume("geofence_alerts", "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Worker failed to consume queue: %v", err)
		return
	}

	for d := range msgs {
		log.Printf("📥 [RABBITMQ WORKER] Receiving Geofence Alert: %s", string(d.Body))
	}
}