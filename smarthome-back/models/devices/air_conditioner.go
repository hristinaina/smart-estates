package devices

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
