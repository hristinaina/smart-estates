package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"simulation/models"
	"strconv"
)

const api = "http://localhost:8081/api"

// GetAllDevices performs a GET request and returns list of devices
func GetAllDevices() ([]models.Device, error) {
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

// Get performs a GET request and returns device data based on device id
func Get(id int) (models.Device, error) {
	url := api + "/devices/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return models.Device{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.Device{}, fmt.Errorf("error reading response body: %v", err)
	}

	var device models.Device
	err = json.Unmarshal(body, &device)
	if err != nil {
		return models.Device{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return device, nil
}

// Get performs a GET request and returns device data based on device id
func GetSP(id int) (models.SolarPanel, error) {
	url := api + "/sp/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return models.SolarPanel{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.SolarPanel{}, fmt.Errorf("error reading response body: %v", err)
	}

	var device models.SolarPanel
	err = json.Unmarshal(body, &device)
	if err != nil {
		return models.SolarPanel{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return device, nil
}
