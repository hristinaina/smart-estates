package enumerations

import (
	"fmt"
	"strings"
)

type Days int

const (
	MONDAY Days = iota
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY
	SUNDAY
)

var daysStringMap = map[string]Days{
	"MONDAY":    MONDAY,
	"TUESDAY":   TUESDAY,
	"WEDNESDAY": WEDNESDAY,
	"THURSDAY":  THURSDAY,
	"FRIDAY":    FRIDAY,
	"SATURDAY":  SATURDAY,
	"SUNDAY":    SUNDAY,
}

func (d Days) String() string {
	switch d {
	case MONDAY:
		return "Monday"
	case TUESDAY:
		return "Tuesday"
	case WEDNESDAY:
		return "Wednesday"
	case THURSDAY:
		return "Thursday"
	case FRIDAY:
		return "Friday"
	case SATURDAY:
		return "Saturday"
	case SUNDAY:
		return "Sunday"
	default:
		return fmt.Sprintf("Unknown day: %d", d)
	}
}

func ConvertStringsToEnumValues(days string) ([]Days, error) {
	var enumValues []Days
	if days != "" {
		if strings.HasSuffix(days, ",") {
			days = days[:len(days)-1]
		}
		selectedDays := strings.Split(days, ",")
		for _, dayStr := range selectedDays {
			day, exists := daysStringMap[strings.ToUpper(dayStr)]
			if !exists {
				return nil, fmt.Errorf("invalid day string: %s", dayStr)
			}
			enumValues = append(enumValues, day)
		}

		return enumValues, nil
	} else {
		return []Days{}, nil
	}
}
