package inside

import models "smarthome-back/models/devices"

// WashingMachine inherits ConsumptionDevice declared as Device attribute
type WashingMachine struct {
	Device   models.ConsumptionDevice
	Mode     []Mode
	ModeName string
}

type ScheduledMode struct {
	Id        int
	DeviceId  int
	StartTime string
	ModeId    int
}

type Mode struct {
	Id          int
	Name        string
	Duration    int
	Temperature string
}

func NewWashingMachine(device models.ConsumptionDevice, mode []Mode) WashingMachine {
	return WashingMachine{
		Device: device,
		Mode:   mode,
	}
}

func (ac WashingMachine) ToDevice() models.Device {
	return models.Device{
		Id:         ac.Device.Device.Id,
		Name:       ac.Device.Device.Name,
		Type:       ac.Device.Device.Type,
		RealEstate: ac.Device.Device.RealEstate,
		IsOnline:   ac.Device.Device.IsOnline,
	}
}
