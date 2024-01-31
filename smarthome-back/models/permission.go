package models

type Permission struct {
	Id           int
	RealEstateId int
	DeviceId     int
	UserEmail    string
	IsActive     bool
	IsDeleted    bool
}
