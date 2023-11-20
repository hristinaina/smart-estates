package main

import (
	_ "database/sql"
	_ "fmt"
	"simulation/devices"
)

func main() {
	//TODO zovni preuzimanje liste uredjaja
	client := devices.CreateConnection()
	//for (device in devices){
	//	if (device.Type == "") {
	go devices.ConnectLamp(client)
	//	} else if (device.Type == "") {
	//		go devices.ConnectLamp(client)
	//	}
	//}
	select {}
}
