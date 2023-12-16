package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"simulation/models"
)

func GetTemp() (models.OpenMeteoTempResponse, error) {
	apiUrl := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", 45.27))
	params.Set("longitude", fmt.Sprintf("%f", 19.83))
	params.Set("current", "temperature_2m,relative_humidity_2m")
	params.Set("houtly", "temperature_2m,relative_humidity_2m")
	params.Set("timezone", "Europe/Berlin")
	fullURL := fmt.Sprintf("%s?%s", apiUrl, params.Encode())
	fmt.Println(fullURL)

	response, err := http.Get(fullURL)
	if err != nil {
		return models.OpenMeteoTempResponse{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		// Read the response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return models.OpenMeteoTempResponse{}, fmt.Errorf("Error reading response body: %v", err)
		}

		var openMeteoResponse models.OpenMeteoTempResponse
		err = json.Unmarshal(body, &openMeteoResponse)
		if err != nil {
			return models.OpenMeteoTempResponse{}, fmt.Errorf("Error unmarshaling JSON response: %v", err)
		}

		return openMeteoResponse, nil

	} else {
		fmt.Println("Open-Meteo API Request Failed. Status Code:", response.StatusCode)
		return models.OpenMeteoTempResponse{}, fmt.Errorf("Open-Meteo API Request Failed. Status Code: %v", response.StatusCode)
	}
}
