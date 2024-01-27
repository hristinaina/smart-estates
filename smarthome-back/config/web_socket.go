package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"net/http"
	"smarthome-back/mqtt_client"
	"smarthome-back/repositories"
	"smarthome-back/services/devices/energetic"
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

	fmt.Println("Client Successfully Connected... 1")

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

func SetupWebSocketRoutes(db *sql.DB, influxDb influxdb2.Client) {
	// http.HandleFunc("/ws", HandleWebSocket)
	http.HandleFunc("/ambient", SendAmbientValues)
	http.HandleFunc("/consumption", func(w http.ResponseWriter, r *http.Request) {
		SendConsumptionValues(db, influxDb, w, r)
	})
}

func SendConsumptionValues(db *sql.DB, influxDb influxdb2.Client, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Client Successfully Connected...")
	realEstateRepository := repositories.NewRealEstateRepository(db)
	homeBatteryService := energetic.NewHomeBatteryService(db, influxDb)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			estates, err := realEstateRepository.GetAll()
			if err != nil {
				return
			}
			for _, estate := range estates {
				value, _ := homeBatteryService.GetConsumptionFromLastMinute(estate.Id)
				fmt.Println(value)
				data := map[string]interface{}{
					"consumed":  value,
					"estateId":  estate.Id,
					"timestamp": time.Now(),
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Println("Error encoding JSON:", err)
					return
				}
				if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					fmt.Println("ERROR when trying to send socket message")
					log.Println(err)
					return
				}
			}
		}
	}
}
