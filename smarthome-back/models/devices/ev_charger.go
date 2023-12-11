package models

import "smarthome-back/enumerations"

// EVCharger inherits Device declared as Device attribute
type EVCharger struct {
	Device        Device
	ChargingPower float64
	Connections   uint
}

func NewEVCharger(device Device, power float64, connections uint) EVCharger {
	return EVCharger{
		Device:        device,
		ChargingPower: power,
		Connections:   connections,
	}
}

func NewEVChargerParam(name string, deviceType enumerations.DeviceType, estate int,
	power float64, connections uint) EVCharger {
	return EVCharger{
		Device:        NewDevice(name, deviceType, estate),
		ChargingPower: power,
		Connections:   connections,
	}
}

func (ac EVCharger) ToDevice() Device {
	return Device{
		Id:         ac.Device.Id,
		Name:       ac.Device.Name,
		Type:       ac.Device.Type,
		RealEstate: ac.Device.RealEstate,
		IsOnline:   ac.Device.IsOnline,
	}
}
