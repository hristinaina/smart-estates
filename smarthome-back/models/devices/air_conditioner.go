package models

// AirConditioner inherits ConsumptionDevice declared as Device attribute
type AirConditioner struct {
	Device         ConsumptionDevice
	MinTemperature float32
	MaxTemperature float32
	Mode           string
	SpecialMode    []SpecialMode
}

func NewAirConditioner(device ConsumptionDevice, minTemp, maxTemp float32, mode string, sc []SpecialMode) AirConditioner {
	return AirConditioner{
		Device:         device,
		MinTemperature: minTemp,
		MaxTemperature: maxTemp,
		Mode:           mode,
		SpecialMode:    sc,
	}
}

func (ac AirConditioner) ToDevice() Device {
	return Device{
		Id:         ac.Device.Device.Id,
		Name:       ac.Device.Device.Name,
		Type:       ac.Device.Device.Type,
		RealEstate: ac.Device.Device.RealEstate,
		IsOnline:   ac.Device.Device.IsOnline,
	}
}
