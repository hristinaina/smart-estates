package models

import "smarthome-back/enumerations"

// ConsumptionDevice Here belong all devices except VEU group
type ConsumptionDevice struct {
	Device           Device
	PowerSupply      enumerations.PowerSupplyType
	PowerConsumption float64 // only if power supply is "home", otherwise it's null
}

func NewConsumptionDeviceParam(name string, deviceType enumerations.DeviceType, estate int,
	powerSupply enumerations.PowerSupplyType, consumption float64) ConsumptionDevice {
	if powerSupply == enumerations.Home {
		return ConsumptionDevice{
			Device:           NewDevice(name, deviceType, estate),
			PowerSupply:      powerSupply,
			PowerConsumption: consumption,
		}
	}
	return ConsumptionDevice{
		Device:           NewDevice(name, deviceType, estate),
		PowerSupply:      powerSupply,
		PowerConsumption: 0,
	}
}

func NewConsumptionDevice(device Device, powerSupply enumerations.PowerSupplyType, consumption float64) ConsumptionDevice {
	if powerSupply == enumerations.Home {
		return ConsumptionDevice{
			Device:           device,
			PowerSupply:      powerSupply,
			PowerConsumption: consumption,
		}
	}
	return ConsumptionDevice{
		Device:           device,
		PowerSupply:      powerSupply,
		PowerConsumption: 0,
	}
}

func (cd ConsumptionDevice) ToDevice() Device {
	return Device{
		Id:         cd.Device.Id,
		Name:       cd.Device.Name,
		Type:       cd.Device.Type,
		RealEstate: cd.Device.RealEstate,
		IsOnline:   cd.Device.IsOnline,
	}
}
