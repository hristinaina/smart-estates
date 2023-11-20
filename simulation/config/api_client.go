package config

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"net/http"
	"simulation/models"
)

const api = "http://localhost:8081/api"

type ApiClient struct {
	client mqtt.Client
}

func NewApiClient(client mqtt.Client) *ApiClient {
	return &ApiClient{client: client}
}

// GetAllDevices performs a GET request and returns list of device_simulator
func (c *ApiClient) GetAllDevices() ([]models.Device, error) {
	url := api + "/devices/"

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var devices []models.Device
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return devices, nil
}
