package repositories

import (
	"database/sql"
	"fmt"
	"smarthome-back/dtos"
	"smarthome-back/enumerations"
	"smarthome-back/models"

	"github.com/go-sql-driver/mysql"
)

type PermissionRepository interface {
	IsPermissionAlreadyExist(realEstateId, deviceId int, userEmail string) bool
	AddInactivePermission(realEstateId, deviceId int, userEmail string) error
	SetActivePermission(email string) error
	GetPermissionByRealEstate(realEstateId int) []dtos.PermissionDTO
	DeletePermission(realEstateId int, permission dtos.PermissionDTO) error
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

type PermissionDevice struct {
	Id              int
	Name            string
	Type            enumerations.DeviceType
	RealEstate      int
	IsOnline        bool
	StatusTimeStamp mysql.NullTime
	LastValue       float32
}

func (res *PermissionRepositoryImpl) GetPermissionByRealEstate(realEstateId int) []dtos.PermissionDTO {
	query := "SELECT p.*, r.*, d.* FROM permission p JOIN realEstate r ON p.RealEstateId = r.Id JOIN device d ON p.DeviceId = d.Id WHERE p.RealEstateId = ? AND p.isActive=true AND p.isDeleted=false"
	rows, err := res.db.Query(query, realEstateId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var permissionsDTO []dtos.PermissionDTO

	for rows.Next() {
		var (
			permission models.Permission
			realEstate models.RealEstate
			device     PermissionDevice

			permissionDTO dtos.PermissionDTO
		)

		if err := rows.Scan(
			&permission.Id, &permission.RealEstateId, &permission.DeviceId,
			&permission.UserEmail, &permission.IsActive, &permission.IsDeleted,
			&realEstate.Id, &realEstate.Name, &realEstate.Type, &realEstate.Address, &realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors, &realEstate.Picture, &realEstate.State, &realEstate.User, &realEstate.DiscardReason,
			&device.Id, &device.Name, &device.Type, &device.RealEstate, &device.IsOnline, &device.StatusTimeStamp, &device.LastValue,
		); err != nil {
			fmt.Println("Error: ", err.Error())
			return []dtos.PermissionDTO{}
		}
		permissionDTO.RealEstate = realEstate.Name
		permissionDTO.UserEmail = permission.UserEmail
		permissionDTO.Device = device.Name
		permissionDTO.DeviceId = device.Id
		permissionsDTO = append(permissionsDTO, permissionDTO)
	}

	return permissionsDTO
}

func (res *PermissionRepositoryImpl) DeletePermission(realEstateId int, permission dtos.PermissionDTO) error {
	updateStatement := "UPDATE permission SET IsDeleted=true WHERE UserEmail=? and isActive=true and RealEstateId=? and DeviceId=?"

	_, err := res.db.Exec(updateStatement, permission.UserEmail, realEstateId, permission.DeviceId)
	return err
}
