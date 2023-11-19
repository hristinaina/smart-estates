package models

// AirConditioner inherits ConsumptionDevice declared as Device attribute
type AirConditioner struct {
	Device         ConsumptionDevice
	MinTemperature int
	MaxTemperature int
}

func NewAirConditioner(device ConsumptionDevice, minTemp int, maxTemp int) AirConditioner {
	return AirConditioner{
		Device:         device,
		MinTemperature: minTemp,
		MaxTemperature: maxTemp,
	}
}

func (ac AirConditioner) ToDevice() Device {
	return Device{
		Id:         ac.Device.Device.Id,
		Name:       ac.Device.Device.Name,
		Type:       ac.Device.Device.Type,
		Picture:    ac.Device.Device.Picture,
		RealEstate: ac.Device.Device.RealEstate,
		IsOnline:   ac.Device.Device.IsOnline,
	}
}
