package repository

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chon26909/tuya-api/config"
)

type TuyaRepository struct {
	cfg *config.Config
}

type DeviceStatusResult struct {
	Result []struct {
		Code  string      `json:"code"`
		Value interface{} `json:"value"`
	} `json:"result"`
	Success bool   `json:"success"`
	T       int64  `json:"t"`
	Tid     string `json:"tid"`
}

type DeviceMetrics struct {
	Switch    bool    `json:"switch_1"`
	Countdown int     `json:"countdown_1"`
	AddEle    int     `json:"add_ele"`
	Current   float64 `json:"cur_current"`
	Power     float64 `json:"cur_power"`
	Voltage   float64 `json:"cur_voltage"`
}

type PowerSensorPayload struct {
	Voltage     float64 `json:"voltage"`
	CurrentMA   float64 `json:"current"`
	PowerMW     float64 `json:"power"`
	LoadVoltage float64 `json:"load_voltage"`
	ShuntMV     float64 `json:"shunt_mV"`
}

func NewTuyaRepository(cfg *config.Config) *TuyaRepository {
	return &TuyaRepository{cfg: cfg}
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func (r *TuyaRepository) sign(method, path, body, accessToken, t string) string {
	contentHash := sha256Hex(body)
	stringToSign := method + "\n" + contentHash + "\n" + "\n" + path
	raw := r.cfg.ClientID + accessToken + t + stringToSign
	mac := hmac.New(sha256.New, []byte(r.cfg.Secret))
	mac.Write([]byte(raw))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

func (r *TuyaRepository) doReq(method, path, body, accessToken string) ([]byte, error) {
	t := fmt.Sprintf("%d", time.Now().UnixMilli())
	sig := r.sign(method, path, body, accessToken, t)

	req, err := http.NewRequest(method, r.cfg.BaseURL+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("client_id", r.cfg.ClientID)
	req.Header.Set("sign", sig)
	req.Header.Set("t", t)
	req.Header.Set("sign_method", "HMAC-SHA256")
	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("access_token", accessToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status=%d body=%s", resp.StatusCode, string(b))
	}
	return b, nil
}

func (r *TuyaRepository) GetToken() (string, error) {
	path := "/v1.0/token?grant_type=1"
	b, err := r.doReq("GET", path, "", "")
	if err != nil {
		return "", err
	}
	var res struct {
		Success bool `json:"success"`
		Result  struct {
			AccessToken string `json:"access_token"`
		} `json:"result"`
		Msg string `json:"msg"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return "", err
	}
	if !res.Success {
		return "", fmt.Errorf("tuya token error: %s body=%s", res.Msg, string(b))
	}
	return res.Result.AccessToken, nil
}

func (r *TuyaRepository) GetDeviceStatus(token string) ([]byte, error) {
	path := "/v1.0/iot-03/devices/" + r.cfg.DeviceID + "/status"
	return r.doReq("GET", path, "", token)
}

func (r *TuyaRepository) GetDeviceMetrics() (*DeviceMetrics, error) {
	token, err := r.GetToken()
	if err != nil {
		return nil, err
	}

	raw, err := r.GetDeviceStatus(token)
	if err != nil {
		return nil, err
	}

	var status DeviceStatusResult
	if err := json.Unmarshal(raw, &status); err != nil {
		return nil, err
	}
	if !status.Success {
		return nil, fmt.Errorf("tuya status error: %s", string(raw))
	}

	metrics := &DeviceMetrics{}
	for _, item := range status.Result {
		switch item.Code {
		case "switch_1":
			if v, ok := item.Value.(bool); ok {
				metrics.Switch = v
			}
		case "countdown_1":
			metrics.Countdown = toInt(item.Value)
		case "add_ele":
			metrics.AddEle = toInt(item.Value)
		case "cur_current":
			metrics.Current = toFloat(item.Value)
		case "cur_power":
			metrics.Power = toFloat(item.Value)
		case "cur_voltage":
			metrics.Voltage = toFloat(item.Value)
		}
	}

	return metrics, nil
}

func (r *TuyaRepository) GetPowerSensorPayload() (*PowerSensorPayload, error) {
	metrics, err := r.GetDeviceMetrics()
	if err != nil {
		return nil, err
	}

	return &PowerSensorPayload{
		Voltage:     metrics.Voltage / 100,
		CurrentMA:   metrics.Current / 1000,
		PowerMW:     metrics.Power / 100,
		LoadVoltage: metrics.Voltage / 100,
		ShuntMV:     0,
	}, nil
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case float32:
		return int(x)
	case int:
		return x
	case int64:
		return int(x)
	case int32:
		return int(x)
	default:
		return 0
	}
}

func toFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case int32:
		return float64(x)
	default:
		return 0
	}
}
