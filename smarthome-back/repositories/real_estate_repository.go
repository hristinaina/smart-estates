package repositories

import (
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"smarthome-back/enumerations"
	"smarthome-back/models"
)

type RealEstateRepository struct {
	db           *sql.DB
	cacheService *cache.CacheService
}

func NewRealEstateRepository(db *sql.DB, cacheService *cache.CacheService) *RealEstateRepository {
	return &RealEstateRepository{db: db, cacheService: cacheService}
}

func (rer *RealEstateRepository) GetAll() ([]models.RealEstate, error) {
	query := "SELECT * FROM realestate"
	rows, err := rer.db.Query(query)
	if IsError(err) {
		return nil, err
	}
	defer rows.Close()
	return ScanRows(rows)
}

func (rer *RealEstateRepository) GetByUserId(id int) ([]models.RealEstate, error) {
	cacheKey := fmt.Sprintf("real_estate_user_%d", id)

	var realEstates []models.RealEstate
	if found, err := rer.cacheService.GetFromCache(cacheKey, &realEstates); found {
		return realEstates, err
	}

	query := "SELECT * FROM realestate WHERE USERID = ?"
	rows, err := rer.db.Query(query, id)

	if IsError(err) {
		return nil, err
	}
	defer rows.Close()

	realEstates, _ = ScanRows(rows)

	if err := rer.cacheService.SetToCache(cacheKey, realEstates); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return ScanRows(rows)
}

func (rer *RealEstateRepository) GetByCity(city string) ([]models.RealEstate, error) {
	query := "SELECT * FROM realestate WHERE city = ?"
	rows, err := rer.db.Query(query, city)

	if IsError(err) {
		return nil, err
	}
	defer rows.Close()
	return ScanRows(rows)
}

func (rer *RealEstateRepository) Get(id int) (models.RealEstate, error) {
	cacheKey := fmt.Sprintf("real_estate_%d", id)

	var realEstate models.RealEstate
	if found, err := rer.cacheService.GetFromCache(cacheKey, &realEstate); found {
		return realEstate, err
	}

	query := "SELECT * FROM realestate WHERE ID = ?"
	rows, err := rer.db.Query(query, id)
	if IsError(err) {
		return models.RealEstate{}, nil
	}
	defer rows.Close()
	estates, err := ScanRows(rows)
	if estates != nil {
		return estates[0], err
	}

	if err := rer.cacheService.SetToCache(cacheKey, estates); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return models.RealEstate{}, err
}

func (rer *RealEstateRepository) GetPending() ([]models.RealEstate, error) {
	query := "SELECT * FROM realestate WHERE STATE = 0"
	rows, err := rer.db.Query(query)
	if IsError(err) {
		return nil, nil
	}
	return ScanRows(rows)
}

func (rer *RealEstateRepository) Delete(id int) error {
	// TODO: finish this
	query := "DELETE FROM realestate WHERE id = ?"
	_, err := rer.db.Exec(query, id)

	if IsError(err) {
		return err
	}
	return nil
}

func (rer *RealEstateRepository) Add(estate models.RealEstate) (models.RealEstate, error) {
	query := "INSERT INTO realestate (Id, Name, Type, Address, City, SquareFootage, NumberOfFloors, Picture, State, UserId, DiscardReason)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := rer.db.Exec(query, estate.Id, estate.Name, estate.Type, estate.Address, estate.City, estate.SquareFootage,
		estate.NumberOfFloors, estate.Picture, estate.State, estate.User, "")
	if IsError(err) {
		return models.RealEstate{}, err
	}
	return estate, nil
}

func (rer *RealEstateRepository) UpdateState(realEstate models.RealEstate) (models.RealEstate, error) {
	query := "UPDATE realestate SET State = ?, DiscardReason = ? WHERE Id = ?"
	_, err := rer.db.Exec(query, realEstate.State, realEstate.DiscardReason, realEstate.Id)
	if IsError(err) {
		return models.RealEstate{}, err
	}
	if realEstate.State == enumerations.ACCEPTED {
		realEstate.State = enumerations.ACCEPTED
	} else {
		realEstate.State = enumerations.DECLINED
	}
	return realEstate, err
}

func IsError(err error) bool {
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
