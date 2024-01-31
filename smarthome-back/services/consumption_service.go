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
	if inputType == "rs" {
		return s.getConsumptionForRealEstates(selectedTime, selectedOptions)

	} else { // input type is "city"
		// Initialize a map to store aggregated values for each city (there can be multiple cities in selectedOptions)
		var results = make(map[string]map[time.Time]float64) //[city][timestamp]value

		for _, city := range selectedOptions {
			estates, _ := s.realEstateService.GetByCity(city)
			estateIds := make([]string, len(estates))
			for i, estate := range estates {
				estateIds[i] = strconv.Itoa(estate.Id)
			}
			// [estate.Name][timestamp]value
			realEstatesMap := s.getConsumptionForRealEstates(selectedTime, estateIds)
			cityAggregatedValues := aggregateResults(realEstatesMap)
			// Store the aggregated values for the city in the results map
			if len(cityAggregatedValues) == 0 {
				continue
			}
			results[city] = cityAggregatedValues
		}

		if len(results) == 0 {
			return nil
		}
		return results
	}
}

func (s *ConsumptionServiceImpl) GetConsumptionForSelectedDate(startDate, endDate string, estateId int) interface{} {
	query := fmt.Sprintf(`from(bucket:"bucket") 
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r._measurement == "consumption" and r["_field"] == "electricity" and r["estate_id"] == "%d")
	|> yield(name: "sum")`, startDate, endDate, estateId)
	selectedTime := getDateDifference(startDate, endDate)
	return s.processingQuery(query, selectedTime)
}

func (s *ConsumptionServiceImpl) getConsumptionForRealEstates(selectedTime string, selectedOptions []string) map[string]map[time.Time]float64 {
	var results = make(map[string]map[time.Time]float64)

	for _, estateId := range selectedOptions {
		estateId, _ := strconv.Atoi(estateId)
		estate, _ := s.realEstateService.Get(estateId)
		query := fmt.Sprintf(`from(bucket:"bucket") 
			|> range(start: %s, stop: now())
			|> filter(fn: (r) => r._measurement == "consumption" and r["_field"] == "electricity" and r["estate_id"] == "%d")
			|> yield(name: "sum")`, selectedTime, estateId)

		tempMap := s.processingQuery(query, selectedTime)
		if len(tempMap) == 0 {
			continue
		}
		results[estate.Name] = tempMap
	}
	if len(results) == 0 {
		return nil
	}
	return results
}

func aggregateResults(results map[string]map[time.Time]float64) map[time.Time]float64 {
	aggregatedMap := make(map[time.Time]float64)

	// Iterate over the outer map (estate names)
	for _, innerMap := range results {
		// Iterate over the inner map (timestamps and values)
		for timestamp, value := range innerMap {
			// Accumulate values for the same timestamp
			aggregatedMap[timestamp] += value
		}
	}

	return aggregatedMap
}

func getDateDifference(startDate, endDate string) string {
	layout := "2006-01-02"

	// Parse the start and end dates
	start, err := time.Parse(layout, startDate)
	if err != nil {
		fmt.Printf("Error parsing start date '%s': %v\n", startDate, err)
		return ""
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		fmt.Printf("Error parsing end date '%s': %v\n", endDate, err)
		return ""
	}

	// Calculate the difference in days
	daysDiff := int(end.Sub(start).Hours() / 24)

	if daysDiff > 3 {
		return "-7d"
	} else {
		return ""
	}
}

func (s *ConsumptionServiceImpl) processingQuery(query string, selectedTime string) map[time.Time]float64 {
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

		// Check if selectedTime is "-7d" and aggregate values by day
		if selectedTime == "-7d" || selectedTime == "-30d" {
			parsedTime = time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, parsedTime.Location())
		}

		resultPoints[parsedTime] += value
	}

	return resultPoints
}
