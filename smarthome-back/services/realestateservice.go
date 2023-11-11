package services

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"smarthome-back/models"
)

type RealEstateService interface {
	GetAll(userId int) []models.RealEstate
	Get(id int) (models.RealEstate, error)
	GetPending() []models.RealEstate
}

type RealEstateServiceImpl struct {
	db       *gorm.DB
	database *sql.DB
}

func NewRealEstateService(db *gorm.DB, database *sql.DB) RealEstateService {
	return &RealEstateServiceImpl{db: db, database: database}
}

func (res *RealEstateServiceImpl) GetAll(userId int) []models.RealEstate {
	query := "SELECT * FROM realestate WHERE USERID = ?"
	rows, err := res.database.Query(query, userId)

	if err != nil {
		fmt.Print("Error: ", err.Error())
		return []models.RealEstate{}
	}
	defer rows.Close()

	var realEstates []models.RealEstate
	for rows.Next() {
		var (
			realEstate models.RealEstate
		)
		if err := rows.Scan(&realEstate.Id, &realEstate.Type, &realEstate.Address,
			&realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors,
			&realEstate.Picture, &realEstate.State, &realEstate.User); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.RealEstate{}
		}
		realEstates = append(realEstates, realEstate)
		fmt.Println(realEstate)
	}

	return realEstates
}

func (res *RealEstateServiceImpl) Get(id int) (models.RealEstate, error) {
	query := "SELECT * FROM realestate WHERE ID = ?"
	rows, err := res.database.Query(query, id)

	if err != nil {
		fmt.Print("Error: ", err.Error())
		return models.RealEstate{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			realEstate models.RealEstate
		)
		if err := rows.Scan(&realEstate.Id, &realEstate.Type, &realEstate.Address,
			&realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors,
			&realEstate.Picture, &realEstate.State, &realEstate.User); err != nil {
			fmt.Println("Error: ", err.Error())
			return models.RealEstate{}, err
		}
		return realEstate, nil
	}
	return models.RealEstate{}, err
}

func (res *RealEstateServiceImpl) GetPending() []models.RealEstate {
	query := "SELECT * FROM realestate WHERE STATE = 0"
	rows, err := res.database.Query(query)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil
	}
	var realEstates []models.RealEstate
	for rows.Next() {
		var (
			realEstate models.RealEstate
		)
		if err := rows.Scan(&realEstate.Id, &realEstate.Type, &realEstate.Address, &realEstate.City,
			&realEstate.SquareFootage, &realEstate.NumberOfFloors, &realEstate.Picture, &realEstate.State,
			&realEstate.User); err != nil {
			fmt.Println("Error: ", err.Error())
			return nil
		}
		realEstates = append(realEstates, realEstate)
	}
	return realEstates

}
