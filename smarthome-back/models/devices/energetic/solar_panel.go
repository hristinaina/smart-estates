package energetic

import (
	"smarthome-back/enumerations"
	"smarthome-back/models/devices"
)

// SolarPanel inherits Device declared as Device attribute
type SolarPanel struct {
	Device         models.Device
	SurfaceArea    float64
	Efficiency     float64
	NumberOfPanels int
	IsOn           bool
}

func NewSolarPanel(device models.Device, surfaceArea float64, efficiency float64, panels int) SolarPanel {
	return SolarPanel{
		Device:         device,
		SurfaceArea:    surfaceArea,
		Efficiency:     efficiency,
		NumberOfPanels: panels,
		IsOn:           false,
	}
}

func NewSolarPanelParam(name string, deviceType enumerations.DeviceType, estate int, surfaceArea float64, efficiency float64, panels int) SolarPanel {
	return SolarPanel{
		Device:         models.NewDevice(name, deviceType, estate),
		SurfaceArea:    surfaceArea,
		Efficiency:     efficiency,
		NumberOfPanels: panels,
		IsOn:           false,
	}
}

func (ac SolarPanel) ToDevice() models.Device {
	return models.Device{
		Id:         ac.Device.Id,
		Name:       ac.Device.Name,
		Type:       ac.Device.Type,
		RealEstate: ac.Device.RealEstate,
		IsOnline:   ac.Device.IsOnline,
	}
}
