package dto

import (
	"encoding/json"
	"fmt"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices"
	"strings"
)

type DeviceDTO struct {
	Id               int
	Name             string
	Type             enumerations.DeviceType
	RealEstate       int
	IsOnline         bool
	PowerSupply      enumerations.PowerSupplyType
	PowerConsumption float64
	MinTemperature   float32
	MaxTemperature   float32
	Mode             string
	SpecialMode      string
	ChargingPower    float64
	Connections      uint
	Size             float64
	UserId           int
	SurfaceArea      float64
	Efficiency       float64
	IsOn             bool
	NumberOfPanels   int
}

// Object conversion has been localized here:

func (dto *DeviceDTO) ToAirConditioner() models.AirConditioner {
	return models.AirConditioner{
		//todo change according to code from front.
		Device: models.ConsumptionDevice{
			Device: models.Device{
				Id:         dto.Id,
				Name:       dto.Name,
				Type:       dto.Type,
				RealEstate: dto.RealEstate,
				IsOnline:   dto.IsOnline,
			},
			PowerSupply:      dto.PowerSupply,
			PowerConsumption: dto.PowerConsumption,
		},
		MinTemperature: dto.MinTemperature,
		MaxTemperature: dto.MaxTemperature,
		Mode:           dto.Mode,
		SpecialMode:    dto.ToSpecialMode(),
	}
}

func (dto *DeviceDTO) ToEVCharger() models.EVCharger {
	return models.EVCharger{
		//todo change according to code from front.
		Device: models.Device{
			Id:         dto.Id,
			Name:       dto.Name,
			Type:       dto.Type,
			RealEstate: dto.RealEstate,
			IsOnline:   dto.IsOnline,
		},
		ChargingPower: dto.ChargingPower,
		Connections:   dto.Connections,
	}
}

func (dto *DeviceDTO) ToHomeBattery() models.HomeBattery {
	return models.HomeBattery{
		//todo change according to code from front.
		Device: models.Device{
			Id:         dto.Id,
			Name:       dto.Name,
			Type:       dto.Type,
			RealEstate: dto.RealEstate,
			IsOnline:   dto.IsOnline,
		},
		Size: dto.Size,
	}
}

func (dto *DeviceDTO) ToSolarPanel() models.SolarPanel {
	return models.SolarPanel{
		Device: models.Device{
			Id:         dto.Id,
			Name:       dto.Name,
			Type:       dto.Type,
			RealEstate: dto.RealEstate,
			IsOnline:   dto.IsOnline,
		},
		SurfaceArea:    dto.SurfaceArea,
		Efficiency:     dto.Efficiency,
		NumberOfPanels: dto.NumberOfPanels,
	}
}

func (dto *DeviceDTO) ToDevice() models.Device {
	return models.Device{
		//todo change according to code from front.
		Id:         dto.Id,
		Name:       dto.Name,
		Type:       dto.Type,
		RealEstate: dto.RealEstate,
		IsOnline:   dto.IsOnline,
	}
}

func (dto *DeviceDTO) ToSpecialMode() []models.SpecialMode {
	var specialModesDTO []models.SpecialModeDTO
	err := json.Unmarshal([]byte(dto.SpecialMode), &specialModesDTO)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	var result []models.SpecialMode

	for _, mode := range specialModesDTO {

		selectedDays := strings.Join(mode.SelectedDays, ",")

		sm := models.NewSpecialMode(mode.Start, mode.End, mode.SelectedMode, mode.Temperature, selectedDays)

		result = append(result, sm)
	}

	fmt.Println(result)
	return result
}
