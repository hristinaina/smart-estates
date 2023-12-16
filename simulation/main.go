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
	devices, err := config.GetAllDevices()
	if err != nil {
		fmt.Println("Error happened:", err)
	} else {
		fmt.Println("Devices successfully loaded!")
	}

	for _, d := range devices {
		device_simulator.StartSimulation(client, d)
	}

	// listen if new device has been added
	config.SubscribeToTopic(client, config.TopicNewDevice+"+", device_simulator.HandleNewDevice)
	openMeteoResponse, err := config.GetTemp()
	fmt.Println(openMeteoResponse)
	fmt.Println(err)

	select {}
}
