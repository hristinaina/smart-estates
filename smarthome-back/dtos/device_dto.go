package dtos

import (
	"encoding/json"
	"fmt"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices"
	"smarthome-back/models/devices/energetic"
	"smarthome-back/models/devices/inside"
	m "smarthome-back/models/devices/outside"
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

func (dto *DeviceDTO) ToAirConditioner() inside.AirConditioner {
	return inside.AirConditioner{
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

func (dto *DeviceDTO) ToEVCharger() energetic.EVCharger {
	return energetic.EVCharger{
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

func (dto *DeviceDTO) ToHomeBattery() energetic.HomeBattery {
	return energetic.HomeBattery{
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

func (dto *DeviceDTO) ToSolarPanel() energetic.SolarPanel {
	return energetic.SolarPanel{
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

func (dto *DeviceDTO) ToLamp() m.Lamp {
	return m.Lamp{
		ConsumptionDevice: models.ConsumptionDevice{
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
		// this is default when the lamp is created
		IsOn:           false,
		LightningLevel: enumerations.OFF,
	}
}

func (dto *DeviceDTO) ToVehicleGate() m.VehicleGate {
	return m.VehicleGate{
		ConsumptionDevice: models.ConsumptionDevice{
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
		IsOpen: false,
		Mode:   enumerations.Private,
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

func (dto *DeviceDTO) ToAmbientSensor() models.ConsumptionDevice {
	return models.ConsumptionDevice{
		Device: models.Device{
			Id:         dto.Id,
			Name:       dto.Name,
			Type:       dto.Type,
			RealEstate: dto.RealEstate,
			IsOnline:   dto.IsOnline,
		},
		PowerSupply:      dto.PowerSupply,
		PowerConsumption: dto.PowerConsumption,
	}
}

func (dto *DeviceDTO) ToSpecialMode() []inside.SpecialMode {
	var specialModesDTO []inside.SpecialModeDTO
	err := json.Unmarshal([]byte(dto.SpecialMode), &specialModesDTO)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	var result []inside.SpecialMode

	for _, mode := range specialModesDTO {

		selectedDays := strings.Join(mode.SelectedDays, ",")

		sm := inside.NewSpecialMode(mode.Start, mode.End, mode.SelectedMode, mode.Temperature, selectedDays)

		result = append(result, sm)
	}

	fmt.Println(result)
	return result
}

func (dto *DeviceDTO) ToSprinkler() m.Sprinkler {
	return m.Sprinkler{
		ConsumptionDevice: models.ConsumptionDevice{
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
		IsOn: false,
	}
}

func (dto *DeviceDTO) ToString() string {
	return "[Device] = Name: " + dto.Name
}
