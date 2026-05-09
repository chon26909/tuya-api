package service

import (
	"github.com/chon26909/tuya-api/config"
	"github.com/chon26909/tuya-api/repository"
)

type TuyaService struct {
	tuyaRepo *repository.TuyaRepository
	mqttRepo *repository.MqttRepository
	cfg      *config.Config
}

func NewTuyaService(tuyaRepo *repository.TuyaRepository, mqttRepo *repository.MqttRepository, cfg *config.Config) *TuyaService {
	return &TuyaService{tuyaRepo: tuyaRepo, mqttRepo: mqttRepo, cfg: cfg}
}

func (s *TuyaService) GetDeviceStatus() ([]byte, error) {
	token, err := s.tuyaRepo.GetToken()
	if err != nil {
		return nil, err
	}
	raw, err := s.tuyaRepo.GetDeviceStatus(token)
	if err != nil {
		return nil, err
	}

	if s.cfg != nil {
		if err := s.PublishDeviceMetricsToMQTT(); err != nil {
			return raw, err
		}
	}

	return raw, nil
}

func (s *TuyaService) PublishDeviceMetricsToMQTT() error {
	metrics, err := s.tuyaRepo.GetPowerSensorPayload()
	if err != nil {
		return err
	}

	if err := s.mqttRepo.Publish("power/home1/voltage", metrics.Voltage); err != nil {
		return err
	}
	if err := s.mqttRepo.Publish("power/home1/current", metrics.CurrentMA); err != nil {
		return err
	}
	return s.mqttRepo.Publish("power/home1/power", metrics.PowerMW)
}
