package main

import (
	"fmt"
	"log"
	"time"

	"github.com/chon26909/tuya-api/config"
	"github.com/chon26909/tuya-api/handler"
	"github.com/chon26909/tuya-api/repository"
	"github.com/chon26909/tuya-api/service"
)

func main() {
	cfg := config.LoadConfig()
	tuyaRepo := repository.NewTuyaRepository(cfg)
	mqttRepo := repository.NewMqttRepository(cfg)
	svc := service.NewTuyaService(tuyaRepo, mqttRepo, cfg)
	h := handler.NewDeviceHandler(svc)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		if err := h.PrintDeviceStatus(); err != nil {
			log.Printf("print device status error: %v", err)
		} else {
			fmt.Println("device status checked")
		}

		<-ticker.C
	}
}
