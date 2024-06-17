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

func (rer *RealEstateRepository) selectQueryForUserId(id int) ([]models.RealEstate, error) {
	query := "SELECT * FROM realestate WHERE USERID = ?"
	rows, err := rer.db.Query(query, id)

	if IsError(err) {
		return nil, err
	}
	defer rows.Close()

	realEstates, err := ScanRows(rows)

	return realEstates, err
}

func (rer *RealEstateRepository) selectQueryForRealEstateId(id int) (models.RealEstate, error) {
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

	return models.RealEstate{}, nil
}

func (rer *RealEstateRepository) GetByUserId(id int) ([]models.RealEstate, error) {
	cacheKey := fmt.Sprintf("real_estate_user_%d", id)

	var realEstates []models.RealEstate
	if found, err := rer.cacheService.GetFromCache(cacheKey, &realEstates); found {
		return realEstates, err
	}

	realEstates, _ = rer.selectQueryForUserId(id)

	if err := rer.cacheService.SetToCache(cacheKey, realEstates); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return realEstates, nil
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

	estates, err := rer.selectQueryForRealEstateId(id)

	if err := rer.cacheService.SetToCache(cacheKey, estates); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return estates, err
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
	estates, err := rer.selectQueryForRealEstateId(id)
	userId := estates.User

	// TODO: finish this
	query := "DELETE FROM realestate WHERE id = ?"
	_, err = rer.db.Exec(query, id)

	if IsError(err) {
		return err
	}

	// refresh cache
	realEstates, _ := rer.selectQueryForUserId(userId)

	cacheKey := fmt.Sprintf("real_estate_user_%d", userId)
	if err := rer.cacheService.SetToCache(cacheKey, realEstates); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	// remove real estate from cache
	cacheKey = fmt.Sprintf("real_estate_%d", id)
	err = rer.cacheService.Cache.Delete(cacheKey)
	if err != nil {
		fmt.Printf("Failed to delete key from cache: %v", err)
	} else {
		fmt.Println("Key deleted successfully")
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

	realEstates, _ := rer.selectQueryForUserId(estate.User)

	cacheKey := fmt.Sprintf("real_estate_user_%d", estate.User)
	if err := rer.cacheService.SetToCache(cacheKey, realEstates); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
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

	realEstates, _ := rer.selectQueryForUserId(realEstate.User)

	cacheKey := fmt.Sprintf("real_estate_user_%d", realEstate.User)
	if err := rer.cacheService.SetToCache(cacheKey, realEstates); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	cacheKey = fmt.Sprintf("real_estate_%d", realEstate.Id)
	if err := rer.cacheService.SetToCache(cacheKey, realEstate); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
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
