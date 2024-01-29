package repositories

import (
	"database/sql"
	"fmt"
)

type PermissionRepository interface {
	IsPermissionAlreadyExist(realEstateId, deviceId int, userEmail string) bool
	AddInactivePermission(realEstateId, deviceId int, userEmail string) error
	SetActivePermission(email string) error
}

type PermissionRepositoryImpl struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

func (res *PermissionRepositoryImpl) IsPermissionAlreadyExist(realEstateId, deviceId int, userEmail string) bool {
	var count int

	query := "SELECT COUNT(*) FROM permission WHERE RealEstateId = ? AND DeviceId = ? AND UserEmail = ? AND isActive = true AND isDeleted = false"

	_, err := res.db.Exec(query, realEstateId, deviceId, userEmail)
	if err != nil {
		fmt.Println(err)
		return true
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (res *PermissionRepositoryImpl) AddInactivePermission(realEstateId, deviceId int, userEmail string) error {
	query := "INSERT INTO permission (RealEstateId, DeviceId, UserEmail, isActive, isDeleted) VALUES (?, ?, ?, ?, ?);"

	_, err := res.db.Exec(query, realEstateId, deviceId, userEmail, false, false)
	if err != nil {
		return fmt.Errorf("Failed to save user: %v", err)

	}
	return nil
}

func (res *PermissionRepositoryImpl) SetActivePermission(email string) error {
	updateStatement := "UPDATE permission SET IsActive=true WHERE UserEmail=? and isDeleted=false"

	_, err := res.db.Exec(updateStatement, email)
	return err
}
