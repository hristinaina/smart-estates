package enumerations

import "fmt"

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
