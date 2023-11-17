package devices

import "smarthome-back/enumerations"

// Device = Base device class.
// Note: solar panel is the only device that doesn't have additional attributes
type Device struct {
	Id         int
	Name       string
	Type       enumerations.DeviceType
	Picture    string // todo change this later (upload picture)
	RealEstate int
	IsOnline   bool
}

func NewDevice(name string, deviceType enumerations.DeviceType, img string, estate int) Device {
	return Device{
		Name:       name,
		Type:       deviceType,
		Picture:    img,
		RealEstate: estate,
		IsOnline:   false,
	}
}
