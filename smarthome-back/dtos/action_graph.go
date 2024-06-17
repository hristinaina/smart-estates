package dtos

type ActionGraphRequest struct {
	DeviceId  int
	UserEmail string
	StartDate string
	EndDate   string
}

type ActionGraphResponse struct {
	Labels []string
	Values []interface{}
}
