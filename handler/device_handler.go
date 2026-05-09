package handler

import (
	"github.com/chon26909/tuya-api/service"
)

type DeviceHandler struct {
	service *service.TuyaService
}

func NewDeviceHandler(service *service.TuyaService) *DeviceHandler {
	return &DeviceHandler{service: service}
}

func (h *DeviceHandler) PrintDeviceStatus() error {
	_, err := h.service.GetDeviceStatus()
	if err != nil {
		return err
	}
	// fmt.Println(string(b))
	return nil
}
