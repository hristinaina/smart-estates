package dto

import (
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices"
	m "smarthome-back/models/devices/outside"
)

type DeviceDTO struct {
	Id               int
	Name             string
	Type             enumerations.DeviceType
	RealEstate       int
	IsOnline         bool
	PowerSupply      enumerations.PowerSupplyType
	PowerConsumption float64
	MinTemperature   int
	MaxTemperature   int
	ChargingPower    float64
	Connections      uint
	Size             float64
	UserId           int
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

func (dto *DeviceDTO) ToLamp() m.Lamp {
	return m.Lamp{
		ConsumptionDevice: models.ConsumptionDevice{
			Device: models.Device{
				Id:         dto.Id,
				Name:       dto.Name,
				Type:       dto.Type,
				Picture:    dto.Picture,
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

func (dto *DeviceDTO) ToString() string {
	return "[Device] = Name: " + dto.Name
}
