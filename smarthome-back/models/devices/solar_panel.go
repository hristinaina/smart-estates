package models

import "smarthome-back/enumerations"

// SolarPanel inherits Device declared as Device attribute
type SolarPanel struct {
	Device      Device
	SurfaceArea float64
	Efficiency  float64
}

func NewSolarPanel(device Device, surfaceArea float64, efficiency float64) SolarPanel {
	return SolarPanel{
		Device:      device,
		SurfaceArea: surfaceArea,
		Efficiency:  efficiency,
	}
}

func NewSolarPanelParam(name string, deviceType enumerations.DeviceType, estate int, surfaceArea float64, efficiency float64) SolarPanel {
	return SolarPanel{
		Device:      NewDevice(name, deviceType, estate),
		SurfaceArea: surfaceArea,
		Efficiency:  efficiency,
	}
}

func (ac SolarPanel) ToDevice() Device {
	return Device{
		Id:         ac.Device.Id,
		Name:       ac.Device.Name,
		Type:       ac.Device.Type,
		RealEstate: ac.Device.RealEstate,
		IsOnline:   ac.Device.IsOnline,
	}
}
