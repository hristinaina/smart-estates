package enumerations

type Role int

const (
	ADMIN Role = iota
	USER
)

func IntToRole(value int) Role {
	switch value {
	case int(ADMIN):
		return ADMIN
	case int(USER):
		return USER
	default:
		return USER
	}
}
