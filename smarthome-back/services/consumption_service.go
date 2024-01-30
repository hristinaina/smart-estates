package services

import (
	"context"
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"strconv"
	"time"
)

type ConsumptionService interface {
	GetConsumptionForSelectedTime(selectedTime string, inputType string, selectedOptions []string) interface{}
	GetConsumptionForSelectedDate(startDate, endDate string, estateId int) interface{}
}

type ConsumptionServiceImpl struct {
	db                *sql.DB
	influxDb          influxdb2.Client
	realEstateService RealEstateService
}

func NewConsumptionService(db *sql.DB, influxDb influxdb2.Client) ConsumptionService {
	return &ConsumptionServiceImpl{db: db, influxDb: influxDb, realEstateService: NewRealEstateService(db)}
}

func (s *ConsumptionServiceImpl) GetConsumptionForSelectedTime(selectedTime string, inputType string, selectedOptions []string) interface{} {
	var results = make(map[string]map[time.Time]float64)

	if inputType == "rs" {
		for _, estateId := range selectedOptions {
			estateId, _ := strconv.Atoi(estateId)
			estate, _ := s.realEstateService.Get(estateId)
			query := fmt.Sprintf(`from(bucket:"bucket") 
			|> range(start: %s, stop: now())
			|> filter(fn: (r) => r._measurement == "consumption" and r["_field"] == "electricity" and r["estate_id"] == "%d")
			|> yield(name: "sum")`, selectedTime, estateId)

			tempMap := s.processingQuery(query)
			if len(tempMap) == 0 {
				continue
			}
			results[estate.Name] = tempMap
		}
	} else if inputType == "city" {
		//todo
		fmt.Println("not implemented")
	}
	if len(results) == 0 {
		return nil
	}
	return results
}

func (s *ConsumptionServiceImpl) GetConsumptionForSelectedDate(startDate, endDate string, estateId int) interface{} {
	query := fmt.Sprintf(`from(bucket:"bucket") 
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r._measurement == "consumption" and r["_field"] == "electricity" and r["estate_id"] == "%d")
	|> yield(name: "sum")`, startDate, endDate, estateId)

	return s.processingQuery(query)
}

func (s *ConsumptionServiceImpl) processingQuery(query string) map[time.Time]float64 {
	Org := "Smart Home"
	queryAPI := s.influxDb.QueryAPI(Org)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing InfluxDB query:", err)
		return nil
	}

	var tempPoints map[string]float64
	tempPoints = make(map[string]float64)

	if err == nil {
		// Iterate over query response
		for result.Next() {
			if result.Record().Value() != nil {
				timeStr := result.Record().Time().Format("2006-01-02 15:04")
				tempPoints[timeStr] = tempPoints[timeStr] + result.Record().ValueByKey("_value").(float64)
			}
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	var resultPoints map[time.Time]float64
	resultPoints = make(map[time.Time]float64)

	for timeStr, value := range tempPoints {
		layout := "2006-01-02 15:04"

		parsedTime, err := time.Parse(layout, timeStr)
		if err != nil {
			fmt.Printf("Error parsing time '%s': %v\n", timeStr, err)
			continue
		}

		resultPoints[parsedTime] = value
	}

	return resultPoints
}
