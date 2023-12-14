package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smarthome-back/mqtt_client"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SendAmbientValues(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Client Successfully Connected...")

	// _, p, err := ws.ReadMessage()
	// if err != nil {
	// 	fmt.Println("GRESKA PRILIKOM CITANJA PORUKE")
	// 	log.Println(err)
	// 	return
	// }

	// values := mqtt_client.GetLastOneHourValues(string(p))

	// jsonData, err := json.Marshal(values)
	// if err != nil {
	// 	fmt.Println("Error encoding JSON:", err)
	// 	return
	// }

	// if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
	// 	fmt.Println("GRESKA PRILIKOM SLANJA PORUKE")
	// 	log.Println(err)
	// 	return
	// }

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			newValue := mqtt_client.GetNewValue()

			// fmt.Println(newValue)

			jsonData, err := json.Marshal(newValue)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}

			// fmt.Println(jsonData)

			if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				fmt.Println("GRESKA PRILIKOM SLANJA PORUKE")
				log.Println(err)
				return
			}
		}
	}
}

func SetupWebSocketRoutes() {
	// http.HandleFunc("/ws", HandleWebSocket)
	http.HandleFunc("/ambient", SendAmbientValues)
}
