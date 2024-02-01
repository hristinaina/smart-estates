package device_simulator

import (
	"fmt"
	"math/rand"
	"simulation/config"
	"simulation/models"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicEVData = "ev/data/"
)

type CarSimulator struct {
	maxCapacity     int
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
	rand.Seed(time.Now().UnixNano())
	maxCapacity := rand.Intn(61) + 20
	startCapacity := (rand.Float64()*(0.6) + 0.1) * float64(maxCapacity)

	return CarSimulator{
		maxCapacity:     maxCapacity,
		startCapacity:   startCapacity,
		currentCapacity: startCapacity,
		active:          true,
	}
}

func updateCarSimulator(car CarSimulator, current float64) CarSimulator {
	return CarSimulator{
		maxCapacity:     car.maxCapacity,
		startCapacity:   car.startCapacity,
		currentCapacity: current,
		active:          true,
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
		maxChargingPercentage: 0.9, //todo ne mora biti na beku, front moze kad udje na stranicu iz baze uzeti posljenu izmjenu (ako je nema onda je defaultna 90)
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

	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		// svakih 10 sekundi provjeravaj ima li slobodan prikljucak i 50/50 kreiraj nit/auto za taj prikljucak
		select {
		case <-ticker.C:
			for connectionId, car := range ev.connections {
				if !car.active {
					fmt.Println("Poziv funkcije koja razmislja da li startovati novo auto")
					randomNumber := rand.Intn(2)
					if randomNumber == 0 {
						//todo poslati frontu i beku pocetak akcije punjenja sa svim podacima auta
						go ev.simulateCarCharging(connectionId)
					}
				}
			}
		}
	}
}

func (ev *EVChargerSimulator) simulateCarCharging(connectionId int) {
	ev.connections[connectionId] = createCarSimulator()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		//svakih 10s uvecaj popunjenost baterije auta i provjeri da li je stiglo do maksimuma
		case <-ticker.C:
			car := ev.connections[connectionId]
			toCharge := ev.device.ChargingPower / 60 / 6 // inace je po satu, 60 je za po minuti i 6 za 10s
			allowedMaxCapacity := float64(car.maxCapacity) * ev.maxChargingPercentage
			if car.currentCapacity+toCharge >= allowedMaxCapacity {
				//todo javi beku i frontu da je zavrseno punjene auta i posalji id prikljucka i id punjaca ofc i naziv akcije
				ev.connections[connectionId] = initCar()
				break
			} else {
				ev.connections[connectionId] = updateCarSimulator(car, car.currentCapacity+toCharge)
				//todo poslati beku koliko je potroseno struje tj toCharge (na onaj consumption topic)
				//todo poslati frontu novu vrijednost sa svim podacima o bateriji (jer je mozda korisnik prvi put usao)
			}
		}
	}
}

//todo mijenjanje maxChargingPercentage
