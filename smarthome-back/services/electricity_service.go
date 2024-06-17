package services

import (
	"context"
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type ElectricityService interface {
	GetElectricityForSelectedTime(queryType string, selectedTime string, inputType string, selectedOptions []string, batteryId string) interface{}
	GetElectricityForSelectedDate(queryType string, startDate, endDate string, inputType string, selectedOptions []string, batteryId string) interface{}
	GetRatioForSelectedDate(startDate, endDate string, inputType string, selectedOptions []string, batteryId string) interface{}
	GetRatioForSelectedTime(selectedTime string, inputType string, selectedOptions []string, batteryId string) interface{}
}

type ElectricityServiceImpl struct {
	db                *sql.DB
	influxDb          influxdb2.Client
	realEstateService RealEstateService
}

func NewElectricityService(db *sql.DB, influxDb influxdb2.Client, cacheService *cache.CacheService) ElectricityService {
	return &ElectricityServiceImpl{db: db, influxDb: influxDb, realEstateService: NewRealEstateService(db, cacheService)}
}

func (uc *ElectricityServiceImpl) GetRatioForSelectedDate(startDate, endDate string, inputType string, selectedOptions []string, batteryId string) interface{} {
	resultsC := uc.GetElectricityForSelectedDate("consumption", startDate, endDate, inputType, selectedOptions, batteryId).(map[string]map[time.Time]float64)
	resultsP := uc.GetElectricityForSelectedDate("solar_panel", startDate, endDate, inputType, selectedOptions, batteryId).(map[string]map[time.Time]float64)
	results := calculateRatio(resultsC, resultsP)
	return results
}

func (uc *ElectricityServiceImpl) GetRatioForSelectedTime(selectedTime string, inputType string, selectedOptions []string, batteryId string) interface{} {
	resultsC := uc.GetElectricityForSelectedTime("consumption", selectedTime, inputType, selectedOptions, batteryId).(map[string]map[time.Time]float64)
	resultsP := uc.GetElectricityForSelectedTime("solar_panel", selectedTime, inputType, selectedOptions, batteryId).(map[string]map[time.Time]float64)
	results := calculateRatio(resultsC, resultsP)
	return results
}

func (s *ElectricityServiceImpl) GetElectricityForSelectedTime(queryType string, selectedTime string, inputType string, selectedOptions []string, batteryId string) interface{} {
	if inputType == "rs" {
		startDate, endDate := calculateDates(selectedTime)
		return s.getElectricityForRealEstates(queryType, startDate, endDate, selectedOptions, batteryId, true)

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
			startDate, endDate := calculateDates(selectedTime)
			realEstatesMap := s.getElectricityForRealEstates(queryType, startDate, endDate, estateIds, "", true)
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

func (s *ElectricityServiceImpl) GetElectricityForSelectedDate(queryType, startDate, endDate string, inputType string, selectedOptions []string, batteryId string) interface{} {
	if inputType == "rs" {
		return s.getElectricityForRealEstates(queryType, startDate, endDate, selectedOptions, batteryId, false)

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
			realEstatesMap := s.getElectricityForRealEstates(queryType, startDate, endDate, estateIds, "", false)
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

func (s *ElectricityServiceImpl) getElectricityForRealEstates(queryType string, startDate, endDate string, selectedOptions []string, batteryId string, isTime bool) map[string]map[time.Time]float64 {
	var results = make(map[string]map[time.Time]float64)

	for _, estateId := range selectedOptions {
		estateId, _ := strconv.Atoi(estateId)
		estate, _ := s.realEstateService.Get(estateId)
		sample := "1h"
		if !isTime {
			sample = "12h"
		}
		if batteryId == "" {
			query := fmt.Sprintf(`from(bucket:"bucket") 
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "%s" and r["_field"] == "electricity" and r["estate_id"] == "%d")
		|> aggregateWindow(every: %s, fn: sum)`, startDate, endDate, queryType, estateId, sample)
			tempMap := s.processingQuery(query, startDate, endDate)
			if len(tempMap) == 0 {
				continue
			}
			results[estate.Name] = tempMap
		} else {
			query := fmt.Sprintf(`from(bucket:"bucket") 
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "%s" and r["_field"] == "electricity" and r["estate_id"] == "%d" and r["battery_id"] == "%s")
		|> aggregateWindow(every: %s, fn: sum)`, startDate, endDate, queryType, estateId, batteryId)
			tempMap := s.processingQuery(query, startDate, endDate)
			if len(tempMap) == 0 {
				continue
			}
			results[estate.Name] = tempMap
		}
	}
	if len(results) == 0 {
		return nil
	}
	return results
}

func calculateDates(selectedTime string) (string, string) {
	endDate := time.Now().UTC()
	var startDate time.Time

	switch selectedTime {
	case "-24h":
		startDate = endDate.Add(-24 * time.Hour)
	case "-6h":
		startDate = endDate.Add(-6 * time.Hour)
	case "-12h":
		startDate = endDate.Add(-12 * time.Hour)
	case "-7d":
		startDate = endDate.Add(-7 * 24 * time.Hour)
	case "-30d":
		startDate = endDate.Add(-30 * 24 * time.Hour)
	default:
		// Handle unsupported selectedTime or provide a default behavior
		fmt.Println("Unsupported selectedTime:", selectedTime)
		return "", ""
	}

	// Format the dates as strings
	startDateStr := startDate.Format("2006-01-02T15:04:05Z")
	endDateStr := endDate.Format("2006-01-02T15:04:05Z")

	return startDateStr, endDateStr
}

func calculateRatio(resultC, resultP map[string]map[time.Time]float64) map[string]map[time.Time]float64 {
	aggregatedMap := make(map[string]map[time.Time]float64)

	// this will go through every city/rs in mapC and aggregate if in mapP or just place value from mapC
	for city, innerMapC := range resultC {
		innerMapP, _ := resultP[city]
		aggregatedMap[city] = make(map[time.Time]float64)

		// Iterate over timestamps in innerMapC
		for timestampC, valueC := range innerMapC {
			// Subtract value from resultC
			aggregatedMap[city][timestampC] -= valueC
		}

		// Iterate over timestamps in innerMapP
		for timestampP, valueP := range innerMapP {
			// Add value from resultP
			aggregatedMap[city][timestampP] += valueP
		}
	}

	// case: city/rs doesn't exist in mapC but exists in mapP
	for city, innerMapP := range resultP {
		_, ok := resultC[city]
		if ok {
			// Don't do anything if the city is present in resultC
			continue
		}
		aggregatedMap[city] = make(map[time.Time]float64)

		// Iterate over timestamps in innerMapP
		for timestampP, valueP := range innerMapP {
			// Subtract value from resultC
			aggregatedMap[city][timestampP] += valueP
		}
	}

	return aggregatedMap
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
	layout := "2006-01-02T15:04:05Z"

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

	if daysDiff > 2 {
		return "-7d"
	} else {
		return ""
	}
}

func (s *ElectricityServiceImpl) processingQuery(query string, startDate, endDate string) map[time.Time]float64 {
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

		selectedTime := getDateDifference(startDate, endDate)
		// Check if selectedTime is "-7d" and aggregate values by day
		if selectedTime == "-7d" {
			parsedTime = time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, parsedTime.Location())
		}

		resultPoints[parsedTime] += value
	}

	return resultPoints
}
