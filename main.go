package main

import (
	"log"

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

	if err := h.PrintDeviceStatus(); err != nil {
		log.Fatal(err)
	}
}
