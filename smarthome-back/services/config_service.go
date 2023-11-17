package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"smarthome-back/models"
)

type ConfigService interface {
	GetToken(path string) (string, error)
	GetAppPassword(path string) (string, error)
	getConfiguration(path string) (models.Configuration, error)
}

type ConfigServiceImpl struct {
}

func NewConfigService() ConfigService {
	return &ConfigServiceImpl{}
}

func (cs *ConfigServiceImpl) GetToken(path string) (string, error) {
	config, err := cs.getConfiguration(path)
	if err != nil {
		return "", err
	}
	return config.InfluxdbToken, nil
}

func (cs *ConfigServiceImpl) GetAppPassword(path string) (string, error) {
	config, err := cs.getConfiguration(path)
	if err != nil {
		return "", err
	}
	return config.AppPassword, nil
}

func (cs *ConfigServiceImpl) getConfiguration(path string) (models.Configuration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading config file: ", err)
		return models.Configuration{}, err
	}

	var config models.Configuration
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return models.Configuration{}, err
	}

	return config, nil
}
