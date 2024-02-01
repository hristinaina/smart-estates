package device_simulator

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicEVData       = "ev/data/" //used to send updates only for front
	topicEVStart      = "ev/start/"
	topicEVEnd        = "ev/end/"
	topicEVPercentage = "ev/percentage/"
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
	percentage, err := config.GetLastPercentage(device.ID)
	if err != nil {
		fmt.Println(err)
	}
	return &EVChargerSimulator{
		client:                client,
		device:                ev,
		connections:           make(map[int]CarSimulator),
		maxChargingPercentage: percentage,
	}
}

func (ls *EVChargerSimulator) ConnectEVCharger() {
	go SendHeartBeat(ls.client, ls.device.Device.ID, ls.device.Device.Name)
	go ls.StartConnections()
	config.SubscribeToTopic(ls.client, topicEVPercentage+strconv.Itoa(ls.device.Device.ID), ls.HandlePercentageChange)
}

func (ev *EVChargerSimulator) StartConnections() {
	for i := 0; i < ev.device.Connections; i++ {
		ev.connections[i] = initCar()
	}

	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		// svakih 15 sekundi provjeravaj ima li slobodan prikljucak i 30/70 kreiraj nit/auto za taj prikljucak
		select {
		case <-ticker.C:
			for connectionId, car := range ev.connections {
				if !car.active {
					fmt.Println("Choosing whether to create new electrical car simulator or not ")
					randomNumber := rand.Intn(3)
					if randomNumber == 0 {
						car := createCarSimulator()
						ev.connections[connectionId] = car
						// send action to front and back (start of charging)
						fmt.Printf("Car created. Electrical charger: id=%d, plugId %d, currentCapacity %f \n", ev.device.Device.ID, connectionId, car.currentCapacity)
						data := map[string]interface{}{
							"PlugId":          connectionId,
							"MaxCapacity":     car.maxCapacity,
							"CurrentCapacity": car.currentCapacity,
							"Active":          true,
							"Action":          "start",
							"Email":           "auto",
						}
						jsonString, err := json.Marshal(data)
						if err != nil {
							fmt.Println("greska")
						}
						config.PublishToTopic(ev.client, topicEVStart+strconv.Itoa(ev.device.Device.ID), string(jsonString))
						go ev.simulateCarCharging(connectionId)
					}
				}
			}
		}
	}
}

func (ev *EVChargerSimulator) simulateCarCharging(connectionId int) {
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		//svakih 10s uvecaj popunjenost baterije auta i provjeri da li je stiglo do maksimuma
		case <-ticker.C:
			car := ev.connections[connectionId]
			toCharge := ev.device.ChargingPower / 60 / 6 // inace je po satu, 60 je za po minuti i 6 za 10s
			allowedMaxCapacity := float64(car.maxCapacity) * ev.maxChargingPercentage
			//send to back and front that car has reached maximum capacity allowed
			if car.currentCapacity+toCharge >= allowedMaxCapacity {
				shouldEnd := ev.handleMaxCapacityReached(car, connectionId, allowedMaxCapacity)
				if shouldEnd {
					return
				}
				//send to front updates about capacity (maximum not reached)
			} else {
				ev.handleCurrentCapacity(car, connectionId, toCharge)
			}
		}
	}
}

func (ev *EVChargerSimulator) handleMaxCapacityReached(oldCar CarSimulator, connectionId int, allowedMaxCapacity float64) bool {
	randomNumber := rand.Intn(3)
	if oldCar.currentCapacity > allowedMaxCapacity {
		allowedMaxCapacity = oldCar.currentCapacity
	}
	car := updateCarSimulator(oldCar, allowedMaxCapacity)
	ev.connections[connectionId] = car
	// although car battery has reached it's maximum capacity, the car can still stay pluged to the charger
	data := map[string]interface{}{
		"PlugId":          connectionId,
		"MaxCapacity":     car.maxCapacity,
		"CurrentCapacity": car.currentCapacity,
		"Active":          true,
		"Action":          "update",
		"Email":           "auto",
	}

	//send to back and front that car has left the station
	if randomNumber == 0 {
		fmt.Printf("Car left the station. Electrical charger: id=%d, plugId %d, currentCapacity %f \n", ev.device.Device.ID, connectionId, car.currentCapacity)
		data["Active"] = false
		data["Action"] = "end"
		jsonString, err := json.Marshal(data)
		if err != nil {
			fmt.Println("greska")
		}
		config.PublishToTopic(ev.client, topicEVEnd+strconv.Itoa(ev.device.Device.ID), string(jsonString))
		ev.connections[connectionId] = initCar()
		return true

	} else { // car is still pluged
		fmt.Printf("Car is full but not leaving. Electrical charger: id=%d, plugId %d, currentCapacity %f \n", ev.device.Device.ID, connectionId, car.currentCapacity)
		jsonString, err := json.Marshal(data)
		if err != nil {
			fmt.Println("greska")
		}
		config.PublishToTopic(ev.client, topicEVData+strconv.Itoa(ev.device.Device.ID), string(jsonString))
		return false
	}
}

func (ev *EVChargerSimulator) handleCurrentCapacity(oldCar CarSimulator, connectionId int, toCharge float64) {
	car := updateCarSimulator(oldCar, oldCar.currentCapacity+toCharge)
	ev.connections[connectionId] = car

	// send to back how much electricity has been consumed
	err := config.PublishToTopic(ev.client, config.TopicConsumption+strconv.Itoa(ev.device.Device.ID), strconv.FormatFloat(toCharge,
		'f', -1, 64))

	// send updated data to front
	data := map[string]interface{}{
		"PlugId":          connectionId,
		"MaxCapacity":     car.maxCapacity,
		"CurrentCapacity": car.currentCapacity,
		"Active":          true,
		"Action":          "update",
		"Email":           "auto",
	}
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println("greska")
	}
	config.PublishToTopic(ev.client, topicEVData+strconv.Itoa(ev.device.Device.ID), string(jsonString))
	fmt.Printf("Car capacity updated. Electrical charger: id=%d, plugId %d, currentCapacity %f \n", ev.device.Device.ID, connectionId, car.currentCapacity)
}

type Percentage struct {
	CurrentCapacity float64
	Email           string
}

func (ev *EVChargerSimulator) HandlePercentageChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}

	var p Percentage
	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(msg.Payload()), &p)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	ev.maxChargingPercentage = math.Round(p.CurrentCapacity*100) / 100
	fmt.Printf("MaxChargingPercentage. Electrical charger: id=%d, percentage %f \n", deviceId, ev.maxChargingPercentage)
}
