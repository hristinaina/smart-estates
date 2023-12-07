package models

import (
	"github.com/go-sql-driver/mysql"
	"smarthome-back/enumerations"
)

// Device = Base device class.
// Note: solar panel is the only device that doesn't have additional attributes (for now!!!)
type Device struct {
	Id              int
	Name            string
	Type            enumerations.DeviceType
	RealEstate      int
	IsOnline        bool
	StatusTimeStamp mysql.NullTime
}

func NewDevice(name string, deviceType enumerations.DeviceType, estate int) Device {
	return Device{
		Name:       name,
		Type:       deviceType,
		RealEstate: estate,
		IsOnline:   false,
	}
}
