package enumerations

type Role int

const (
	ADMIN Role = iota
	USER
	SUPERADMIN
)

func IntToRole(value int) Role {
	switch value {
	case int(ADMIN):
		return ADMIN
	case int(SUPERADMIN):
		return SUPERADMIN
	default:
		return USER
	}
}
