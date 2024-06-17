package energetic

import (
	"smarthome-back/enumerations"
	"smarthome-back/models/devices"
)

// EVCharger inherits Device declared as Device attribute
type EVCharger struct {
	Device        models.Device
	ChargingPower float64
	Connections   uint
}

func NewEVCharger(device models.Device, power float64, connections uint) EVCharger {
	return EVCharger{
		Device:        device,
		ChargingPower: power,
		Connections:   connections,
	}
}

func NewEVChargerParam(name string, deviceType enumerations.DeviceType, estate int,
	power float64, connections uint) EVCharger {
	return EVCharger{
		Device:        models.NewDevice(name, deviceType, estate),
		ChargingPower: power,
		Connections:   connections,
	}
}

func (ac EVCharger) ToDevice() models.Device {
	return models.Device{
		Id:         ac.Device.Id,
		Name:       ac.Device.Name,
		Type:       ac.Device.Type,
		RealEstate: ac.Device.RealEstate,
		IsOnline:   ac.Device.IsOnline,
	}
}
