package main

import (
	"log"
	"time"

	"github.com/chon26909/tuya-api/config"
	"github.com/chon26909/tuya-api/repository"
	"github.com/chon26909/tuya-api/service"
)

func main() {
	cfg := config.LoadConfig()
	tuyaRepo := repository.NewTuyaRepository(cfg)
	mqttRepo := repository.NewMqttRepository(cfg)
	svc := service.NewTuyaService(tuyaRepo, mqttRepo, cfg)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		if err := svc.PublishDeviceMetricsToMQTT(); err != nil {
			log.Println("publish error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("published to mqtt topics: power/home1/voltage, power/home1/current, power/home1/power")
		<-ticker.C
	}
}
