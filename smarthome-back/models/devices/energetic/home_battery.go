package energetic

import (
	"smarthome-back/enumerations"
	"smarthome-back/models/devices"
)

// HomeBattery inherits Device declared as Device attribute
type HomeBattery struct {
	Device       models.Device
	Size         float64
	CurrentValue float64
}

func NewHomeBattery(device models.Device, size float64) HomeBattery {
	return HomeBattery{
		Device:       device,
		Size:         size,
		CurrentValue: 0,
	}
}

func NewHomeBatteryParam(name string, deviceType enumerations.DeviceType, estate int, size float64) HomeBattery {
	return HomeBattery{
		Device:       models.NewDevice(name, deviceType, estate),
		Size:         size,
		CurrentValue: 0,
	}
}

func (ac HomeBattery) ToDevice() models.Device {
	return models.Device{
		Id:         ac.Device.Id,
		Name:       ac.Device.Name,
		Type:       ac.Device.Type,
		RealEstate: ac.Device.RealEstate,
		IsOnline:   ac.Device.IsOnline,
	}
}
