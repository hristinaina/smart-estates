package repositories

import (
	"database/sql"
	"fmt"
	"log"
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
	GetPermitRealEstateByEmail(email string) []models.RealEstate
	GetDeviceForSharedRealEstate(email string, realEstateId int) []PermissionDevice
	GetPermissionsForDevice(deviceId int, estateId int) []string
	GetAllUsersForRealEstate(deviceId int, realEstateId int) []string
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

func (res *PermissionRepositoryImpl) GetPermitRealEstateByEmail(email string) []models.RealEstate {
	query := `SELECT DISTINCT r.*
	FROM smart_home.permission p
	JOIN smart_home.realestate r ON p.RealEstateId = r.Id
	WHERE p.UserEmail = ? AND p.isActive=true and p.isDeleted=false`

	rows, err := res.db.Query(query, email)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var realEstates []models.RealEstate

	for rows.Next() {
		var (
			realEstate models.RealEstate
		)

		if err := rows.Scan(
			&realEstate.Id, &realEstate.Name, &realEstate.Type, &realEstate.Address, &realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors, &realEstate.Picture, &realEstate.State, &realEstate.User, &realEstate.DiscardReason,
		); err != nil {
			return []models.RealEstate{}
		}
		realEstates = append(realEstates, realEstate)
	}

	return realEstates
}

func (res *PermissionRepositoryImpl) GetDeviceForSharedRealEstate(email string, realEstateId int) []PermissionDevice {
	query := `SELECT d.*
	FROM permission p
	JOIN device d ON p.DeviceId = d.Id
	WHERE p.UserEmail = ? AND p.RealEstateId = ? AND p.isActive = true AND p.isDeleted = false`

	rows, err := res.db.Query(query, email, realEstateId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var devices []PermissionDevice

	for rows.Next() {
		var (
			device PermissionDevice
		)

		if err := rows.Scan(
			&device.Id, &device.Name, &device.Type, &device.RealEstate, &device.IsOnline, &device.StatusTimeStamp, &device.LastValue,
		); err != nil {
			return []PermissionDevice{}
		}
		devices = append(devices, device)
	}

	return devices
}

func (res *PermissionRepositoryImpl) GetPermissionsForDevice(deviceId int, estateId int) []string {
	query := `SELECT DISTINCT p.userEmail
	FROM smart_home.permission p
	WHERE p.deviceId  = ? OR p.realEstateId =? and p.isDeleted=false`

	rows, err := res.db.Query(query, deviceId, estateId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var userEmails []string

	for rows.Next() {
		var userEmail string
		if err := rows.Scan(&userEmail); err != nil {
			log.Printf("Row scan error: %v", err)
			return nil
		}
		userEmails = append(userEmails, userEmail)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil
	}

	return userEmails
}

func (res *PermissionRepositoryImpl) getOwnerNameByRealEstateId(realEstateId int) (string, error) {
	query := `SELECT u.Name, u.Surname
			  FROM user u
			  JOIN realestate r ON u.Id = r.UserId
			  WHERE r.Id = ?`

	var name, surname string
	err := res.db.QueryRow(query, realEstateId).Scan(&name, &surname)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no owner found for real estate id %d", realEstateId)
		}
		return "", err
	}

	fullName := fmt.Sprintf("%s %s", name, surname)
	return fullName, nil
}

func (res *PermissionRepositoryImpl) GetAllUsersForRealEstate(deviceId int, estateId int) []string {
	query := `SELECT DISTINCT u.Name, u.Surname
			  FROM smart_home.permission p
			  JOIN user u ON p.userEmail = u.Email
			  WHERE (p.deviceId = ? AND p.realEstateId = ?) AND p.isDeleted = false`

	rows, err := res.db.Query(query, deviceId, estateId)
	if err != nil {
		log.Printf("Error querying permissions: %v", err)
		return nil
	}
	defer rows.Close()

	var users []string

	for rows.Next() {
		var name, surname string
		if err := rows.Scan(&name, &surname); err != nil {
			log.Printf("Row scan error: %v", err)
			return nil
		}
		fullName := fmt.Sprintf("%s %s", name, surname)
		users = append(users, fullName)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil
	}

	fmt.Println(users)

	owner, err := res.getOwnerNameByRealEstateId(estateId)
	if err == nil {
		users = append(users, owner)
	}

	fmt.Println(owner)

	return users
}
