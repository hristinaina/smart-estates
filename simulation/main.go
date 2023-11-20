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
	} else {
		fmt.Println("Devices successfully loaded!")
	}

	for _, d := range devices {
		switch d.Type {
		case 0:
			fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
			go device_simulator.ConnectLamp(client, d)
		case 1:
			fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
			go device_simulator.ConnectLamp(client, d)
			//todo add logic for other cases/device types
		default:
			fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
			// Default logic or error handling for unknown device types
			go device_simulator.ConnectLamp(client, d)
		}
	}

	//todo ovdje dodati da slusa za dodavanje novog uredjaja i prebaciti case u novu funkciju i isto je pozavti odavde

	select {}
}
