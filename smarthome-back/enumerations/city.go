package enumerations

type City string

const (
	NoviSad   City = "Novi Sad"
	Zrenjanin City = "Zrenjanin"
	Beograd   City = "Beograd"
)

func AllCities() []string {
	return []string{
		string(NoviSad),
		string(Zrenjanin),
		string(Beograd),
	}
}
