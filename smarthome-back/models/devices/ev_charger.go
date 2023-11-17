package devices

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

func NewEVChargerParam(name string, deviceType enumerations.DeviceType, img string, estate int,
	power float64, connections uint) EVCharger {
	return EVCharger{
		Device:        NewDevice(name, deviceType, img, estate),
		ChargingPower: power,
		Connections:   connections,
	}
}
