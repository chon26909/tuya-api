package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL   string
	ClientID  string
	Secret    string
	DeviceID  string
	MQTTHost  string
	MQTTPort  uint16
	MQTTUser  string
	MQTTPass  string
	MQTTTopic string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	return &Config{
		BaseURL:   getEnv("TUYA_BASE_URL", ""),
		ClientID:  getEnv("TUYA_CLIENT_ID", ""),
		Secret:    getEnv("TUYA_SECRET", ""),
		DeviceID:  getEnv("TUYA_DEVICE_ID", ""),
		MQTTHost:  getEnv("MQTT_HOST", ""),
		MQTTPort:  parseUint16(getEnv("MQTT_PORT", "0")),
		MQTTUser:  getEnv("MQTT_USER", ""),
		MQTTPass:  getEnv("MQTT_PASS", ""),
		MQTTTopic: getEnv("MQTT_TOPIC", ""),
	}
}

func parseUint16(s string) uint16 {
	var v uint16
	_, err := fmt.Sscanf(s, "%d", &v)
	if err != nil {
		return 0
	}
	return v
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
