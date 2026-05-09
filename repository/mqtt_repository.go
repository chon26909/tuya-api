package repository

import (
	"encoding/json"
	"fmt"

	"github.com/chon26909/tuya-api/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttRepository struct {
	cfg *config.Config
}

func NewMqttRepository(cfg *config.Config) *MqttRepository {
	return &MqttRepository{cfg: cfg}
}

func (r *MqttRepository) Publish(topic string, payload interface{}) error {
	broker := fmt.Sprintf("tcp://%s:%d", r.cfg.MQTTHost, r.cfg.MQTTPort)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(r.cfg.MQTTUser)
	opts.SetPassword(r.cfg.MQTTPass)
	opts.SetClientID("tuya-publisher")

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250)

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	pub := client.Publish(topic, 1, false, b)
	pub.Wait()
	return pub.Error()
}
