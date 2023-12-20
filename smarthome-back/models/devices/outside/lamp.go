package models

import (
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices"
)

// Lamp inherits Device class
type Lamp struct {
	ConsumptionDevice models.ConsumptionDevice
	IsOn              bool
	LightningLevel    enumerations.LampLightningLevel
}

func NewLamp(device models.ConsumptionDevice) Lamp {
	return Lamp{
		ConsumptionDevice: device,
		IsOn:              false,
		LightningLevel:    enumerations.OFF,
	}
}

func (lamp Lamp) ToDevice() models.Device {
	return models.Device{
		Id:         lamp.ConsumptionDevice.Device.Id,
		Name:       lamp.ConsumptionDevice.Device.Name,
		Type:       lamp.ConsumptionDevice.Device.Type,
		RealEstate: lamp.ConsumptionDevice.Device.RealEstate,
		IsOnline:   lamp.ConsumptionDevice.Device.IsOnline,
	}
}
