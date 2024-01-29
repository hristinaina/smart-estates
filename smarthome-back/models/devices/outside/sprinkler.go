package models

import models "smarthome-back/models/devices"

type Sprinkler struct {
	ConsumptionDevice models.ConsumptionDevice
	IsOn              bool
	SpecialMode       SprinklerSpecialMode
}

func NewSprinkler(device models.ConsumptionDevice, isOn bool, specialMode SprinklerSpecialMode) Sprinkler {
	return Sprinkler{
		ConsumptionDevice: device,
		IsOn:              isOn,
		SpecialMode:       specialMode,
	}
}

func (s Sprinkler) ToDevice() models.Device {
	return models.Device{
		Id:         s.ConsumptionDevice.Device.Id,
		Name:       s.ConsumptionDevice.Device.Name,
		Type:       s.ConsumptionDevice.Device.Type,
		RealEstate: s.ConsumptionDevice.Device.RealEstate,
		IsOnline:   s.ConsumptionDevice.Device.IsOnline,
	}
}
