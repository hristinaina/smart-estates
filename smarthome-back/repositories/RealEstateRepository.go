package repositories

import (
	"database/sql"
	"fmt"
	"smarthome-back/models"
)

type RealEstateRepository struct {
	db *sql.DB
}

func NewRealEstateRepository(db *sql.DB) *RealEstateRepository {
	return &RealEstateRepository{db: db}
}

func (rer *RealEstateRepository) GetAll() ([]models.RealEstate, error) {
	query := "SELECT * FROM realestate"
	rows, err := rer.db.Query(query)
	if CheckIfError(err) {
		return []models.RealEstate{}, err
	}
	defer rows.Close()
	return ScanRows(rows)
}

func CheckIfError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return true
	}
	return false
}

// ScanRows mapping returned value from db to model
func ScanRows(rows *sql.Rows) ([]models.RealEstate, error) {
	var realEstates []models.RealEstate
	for rows.Next() {
		var (
			realEstate models.RealEstate
		)

		if err := rows.Scan(&realEstate.Id, &realEstate.Name, &realEstate.Type, &realEstate.Address,
			&realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors,
			&realEstate.Picture, &realEstate.State, &realEstate.User, &realEstate.DiscardReason); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.RealEstate{}, err
		}
		realEstates = append(realEstates, realEstate)
	}

	return realEstates, nil
}
