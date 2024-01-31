package device_simulator

import (
	"fmt"
	"simulation/config"
	"simulation/models"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicEVData = "ev/data/"
)

type CarSimulator struct {
	maxCapacity     float64
	currentCapacity float64
	startCapacity   float64
	active          bool
}

func initCar() CarSimulator {
	return CarSimulator{
		active: false,
	}
}

func createCarSimulator() CarSimulator {
	return CarSimulator{
		//todo random generisati podatke za ostale attribute
		active: true,
	}
}

type EVChargerSimulator struct {
	client                mqtt.Client
	device                models.EVCharger
	connections           map[int]CarSimulator
	maxChargingPercentage float64
}

func NewEVChargerSimulator(client mqtt.Client, device models.Device) *EVChargerSimulator {
	ev, err := config.GetEVCharger(device.ID)
	if err != nil {
		return nil
	}
	return &EVChargerSimulator{
		client:                client,
		device:                ev,
		connections:           make(map[int]CarSimulator),
		maxChargingPercentage: 90, //todo ne mora biti na beku, front moze kad udje na stranicu iz baze uzeti posljenu izmjenu (ako je nema onda je defaultna 90)
	}
}

func (ls *EVChargerSimulator) ConnectEVCharger() {
	go SendHeartBeat(ls.client, ls.device.Device.ID, ls.device.Device.Name)
	go ls.StartConnections()
}

func (ev *EVChargerSimulator) StartConnections() {
	for i := 0; i < ev.device.Connections; i++ {
		ev.connections[i] = initCar()
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		//todo svakih 10 sekundi provjeravaj ima li slobodan prikljucak i 50/50 kreiraj nit za to auto
		select {
		case <-ticker.C:
			for _, car := range ev.connections {
				if !car.active {
					fmt.Println("Poziv funkcije koja razmislja da li startovati novo auto")
				}
			}
		}
	}
}
