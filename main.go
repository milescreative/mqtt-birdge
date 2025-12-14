package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func main() {
	// Get MQTT broker from environment
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		log.Fatal("MQTT_BROKER environment variable is required")
	}

	// Setup MQTT client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("pi-http-mqtt-bridge")
	opts.SetConnectTimeout(5 * time.Second)

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer mqttClient.Disconnect(250)

	log.Printf("Connected to MQTT broker: %s", broker)

	// Setup HTTP server
	http.HandleFunc("/publish", publishHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	topic := r.URL.Query().Get("mqtt_topic")
	message := r.URL.Query().Get("mqtt_message")

	if topic == "" || message == "" {
		http.Error(w, "Both mqtt_topic and mqtt_message query parameters are required", http.StatusBadRequest)
		return
	}

	// Publish to MQTT
	token := mqttClient.Publish(topic, 0, false, message)
	if token.Wait() && token.Error() != nil {
		log.Printf("Failed to publish message: %v", token.Error())
		http.Error(w, fmt.Sprintf("Failed to publish: %v", token.Error()), http.StatusInternalServerError)
		return
	}

	log.Printf("Published to topic '%s': %s", topic, message)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message published to topic: %s\n", topic)
}
