package models

import (
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices"
)

type VehicleGate struct {
	ConsumptionDevice models.ConsumptionDevice
	IsOpen            bool
	Mode              enumerations.VehicleGateMode
}

func NewVehicleGate(device models.ConsumptionDevice) VehicleGate {
	return VehicleGate{
		ConsumptionDevice: device,
		IsOpen:            false,
		Mode:              enumerations.Private,
	}
}

func (vehicleGate VehicleGate) ToDevice() models.Device {
	return models.Device{
		Id:         vehicleGate.ConsumptionDevice.Device.Id,
		Name:       vehicleGate.ConsumptionDevice.Device.Name,
		Type:       vehicleGate.ConsumptionDevice.Device.Type,
		RealEstate: vehicleGate.ConsumptionDevice.Device.RealEstate,
		IsOnline:   vehicleGate.ConsumptionDevice.Device.IsOnline,
	}
}
