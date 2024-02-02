package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
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

// Get performs a GET request and returns device data based on device id
func GetEVCharger(id int) (models.EVCharger, error) {
	url := api + "/ev/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return models.EVCharger{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.EVCharger{}, fmt.Errorf("error reading response body: %v", err)
	}

	var device models.EVCharger
	err = json.Unmarshal(body, &device)
	if err != nil {
		return models.EVCharger{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return device, nil
}

// Get performs a GET request and returns device data based on device id
func GetLastPercentage(id int) (float64, error) {
	url := api + "/ev/lastPercentage/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return 0.9, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0.9, fmt.Errorf("unexpected status code: %v", response.Status)
	}

	var lastPercentage float64
	err = json.NewDecoder(response.Body).Decode(&lastPercentage)
	if err != nil {
		return 0.9, fmt.Errorf("error decoding JSON: %v", err)
	}

	return math.Round(lastPercentage*100) / 100, nil
}

func GetAC(id int) (models.AirConditioner, error) {
	url := api + "/ac/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return models.AirConditioner{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return models.AirConditioner{}, fmt.Errorf("error reading response body: %v", err)
	}

	var device models.AirConditioner
	err = json.Unmarshal(body, &device)
	if err != nil {
		fmt.Println(err)
		return models.AirConditioner{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return device, nil
}

func GetAmbientSensor(id int) (models.AmbientSensor, error) {
	url := api + "/ambient/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return models.AmbientSensor{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.AmbientSensor{}, fmt.Errorf("error reading response body: %v", err)
	}

	var device models.AmbientSensor
	err = json.Unmarshal(body, &device)
	if err != nil {
		return models.AmbientSensor{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return device, nil
}

func GetWashingMachine(id int) (models.WashingMachine, error) {
	url := api + "/wm/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return models.WashingMachine{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.WashingMachine{}, fmt.Errorf("error reading response body: %v", err)
	}

	var device models.WashingMachine
	err = json.Unmarshal(body, &device)
	if err != nil {
		return models.WashingMachine{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return device, nil
}

func GetWashingMachineScheduledMode(id int) ([]models.ScheduledMode, error) {
	url := api + "/wm/schedule/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return []models.ScheduledMode{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []models.ScheduledMode{}, fmt.Errorf("error reading response body: %v", err)
	}

	var modes []models.ScheduledMode
	err = json.Unmarshal(body, &modes)
	if err != nil {
		return []models.ScheduledMode{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return modes, nil
}

func GetSolarRadiation(latitude float64, longitude float64) (models.OpenMeteoResponse, error) {
	apiUrl := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", latitude))
	params.Set("longitude", fmt.Sprintf("%f", longitude))
	params.Set("hourly", "direct_normal_irradiance")
	params.Set("forecast_days", "1")
	fullURL := fmt.Sprintf("%s?%s", apiUrl, params.Encode())
	fmt.Println(fullURL)

	response, err := http.Get(fullURL)
	if err != nil {
		return models.OpenMeteoResponse{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		// Read the response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return models.OpenMeteoResponse{}, fmt.Errorf("Error reading response body: %v", err)
		}

		var openMeteoResponse models.OpenMeteoResponse
		err = json.Unmarshal(body, &openMeteoResponse)
		if err != nil {
			return models.OpenMeteoResponse{}, fmt.Errorf("Error unmarshaling JSON response: %v", err)
		}

		// SolarRadiation is in W/m^2
		return openMeteoResponse, nil
	} else {
		fmt.Println("Open-Meteo API Request Failed. Status Code:", response.StatusCode)
		return models.OpenMeteoResponse{}, fmt.Errorf("Open-Meteo API Request Failed. Status Code: %v", response.StatusCode)
	}

}

func GetConsumptionDevice(id int) (models.ConsumptionDevice, error) {
	url := api + "/devices/consumption-device/" + strconv.Itoa(id)

	response, err := http.Get(url)
	if err != nil {
		return models.ConsumptionDevice{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.ConsumptionDevice{}, fmt.Errorf("error reading response body: %v", err)
	}

	var device models.ConsumptionDevice
	err = json.Unmarshal(body, &device)
	if err != nil {
		return models.ConsumptionDevice{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return device, nil

}
