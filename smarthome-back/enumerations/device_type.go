package enumerations

type DeviceType int

const (
	AmbientSensor DeviceType = iota
	AirConditioner
	WashingMachine
	Lamp
	VehicleGate
	Sprinkler
	SolarPanel
	BatteryStorage
	EVCharger
)
