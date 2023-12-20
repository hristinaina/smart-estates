package dtos

import "smarthome-back/enumerations"

type ConsumptionDeviceDto struct {
	PowerSupply      enumerations.PowerSupplyType
	PowerConsumption float64
}
