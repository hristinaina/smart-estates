package models

import "smarthome-back/enumerations"

// HomeBattery inherits Device declared as Device attribute
type HomeBattery struct {
	Device Device
	Size   float64
}

func NewHomeBattery(device Device, size float64) HomeBattery {
	return HomeBattery{
		Device: device,
		Size:   size,
	}
}

func NewHomeBatteryParam(name string, deviceType enumerations.DeviceType, img string, estate int, size float64) HomeBattery {
	return HomeBattery{
		Device: NewDevice(name, deviceType, img, estate),
		Size:   size,
	}
}
