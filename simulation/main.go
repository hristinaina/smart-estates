package main

import (
	_ "database/sql"
	"fmt"
	_ "fmt"
	"simulation/config"
	"simulation/device_simulator"
)

func main() {
	client := config.CreateConnection()
	controller := config.NewApiClient(client) //todo mozda mu nece trebati klijent
	devices, err := controller.GetAllDevices()
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range devices {
		switch d.Type {
		case 0:
			fmt.Printf("Handling logic for type=0 device - ID: %d, Name: %s\n", d.ID, d.Name)
			go device_simulator.ConnectLamp(client)
		case 1:
			fmt.Printf("Handling logic for type=1 device - ID: %d, Name: %s\n", d.ID, d.Name)
			// Your logic for type=1 device goes here
		//todo add logic for other cases/device types
		default:
			fmt.Printf("Unknown device type - ID: %d, Name: %s\n", d.ID, d.Name)
			// Default logic or error handling for unknown device types
		}
	}

	//todo ovdje dodati da slusa za dodavanje novog uredjaja i prebaciti case u novu funkciju i isto je pozavti odavde

	select {}
}
