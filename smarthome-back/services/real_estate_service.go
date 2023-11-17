package services

import (
	"database/sql"
	"fmt"
	"smarthome-back/enumerations"
	"smarthome-back/models"
)

type RealEstateService interface {
	GetAll() []models.RealEstate
	GetAllByUserId(userId int) []models.RealEstate
	Get(id int) (models.RealEstate, error)
	GetPending() []models.RealEstate
	// ChangeState if state == 0 it is accepted, in opposite it is declined
	ChangeState(id int, state int, reason string) models.RealEstate
	Add(estate models.RealEstate) models.RealEstate
}

type RealEstateServiceImpl struct {
	db          *sql.DB
	mailService MailService
}

func NewRealEstateService(db *sql.DB) RealEstateService {
	return &RealEstateServiceImpl{db: db, mailService: NewMailService(db)}
}

func (res *RealEstateServiceImpl) GetAll() []models.RealEstate {
	query := "SELECT * FROM realestate"
	rows, err := res.db.Query(query)
	if CheckIfError(err) {
		return nil
	}
	defer rows.Close()

	var realEstates []models.RealEstate
	for rows.Next() {
		var (
			realEstate models.RealEstate
		)

		// TODO: create function for this -> DRY
		if err := rows.Scan(&realEstate.Id, &realEstate.Name, &realEstate.Type, &realEstate.Address,
			&realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors,
			&realEstate.Picture, &realEstate.State, &realEstate.User, &realEstate.DiscardReason); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.RealEstate{}
		}
		realEstates = append(realEstates, realEstate)
	}

	return realEstates
}

func (res *RealEstateServiceImpl) GetAllByUserId(userId int) []models.RealEstate {
	query := "SELECT * FROM realestate WHERE USERID = ?"
	rows, err := res.db.Query(query, userId)

	if CheckIfError(err) {
		return nil
	}
	defer rows.Close()

	var realEstates []models.RealEstate
	for rows.Next() {
		var (
			realEstate models.RealEstate
		)
		if err := rows.Scan(&realEstate.Id, &realEstate.Name, &realEstate.Type, &realEstate.Address,
			&realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors,
			&realEstate.Picture, &realEstate.State, &realEstate.User, &realEstate.DiscardReason); err != nil {
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
	rows, err := res.db.Query(query, id)

	if CheckIfError(err) {
		return models.RealEstate{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var (
			realEstate models.RealEstate
		)
		if err := rows.Scan(&realEstate.Id, &realEstate.Name, &realEstate.Type, &realEstate.Address,
			&realEstate.City, &realEstate.SquareFootage, &realEstate.NumberOfFloors,
			&realEstate.Picture, &realEstate.State, &realEstate.User, &realEstate.DiscardReason); err != nil {
			fmt.Println("Error: ", err.Error())
			return models.RealEstate{}, err
		}
		return realEstate, nil
	}
	return models.RealEstate{}, err
}

func (res *RealEstateServiceImpl) GetPending() []models.RealEstate {
	query := "SELECT * FROM realestate WHERE STATE = 0"
	rows, err := res.db.Query(query)

	if CheckIfError(err) {
		return nil
	}
	var realEstates []models.RealEstate
	for rows.Next() {
		var (
			realEstate models.RealEstate
		)
		if err := rows.Scan(&realEstate.Id, &realEstate.Name, &realEstate.Type, &realEstate.Address, &realEstate.City,
			&realEstate.SquareFootage, &realEstate.NumberOfFloors, &realEstate.Picture, &realEstate.State,
			&realEstate.User, &realEstate.DiscardReason); err != nil {
			fmt.Println("Error: ", err.Error())
			return nil
		}
		realEstates = append(realEstates, realEstate)
	}
	return realEstates
}

func (res *RealEstateServiceImpl) ChangeState(id int, state int, reason string) models.RealEstate {
	realEstate, err := res.Get(id)

	if CheckIfError(err) {
		return models.RealEstate{}
	}

	if realEstate.State != enumerations.PENDING {
		return models.RealEstate{}
	}

	// delete old data
	query := "DELETE FROM realestate WHERE id = ?"
	_, err = res.db.Exec(query, id)

	if CheckIfError(err) {
		return models.RealEstate{}
	}

	// insert updated data
	if state == 0 {
		realEstate.State = enumerations.ACCEPTED
		err = res.mailService.ApproveRealEstate(realEstate)
	} else {
		realEstate.State = enumerations.DECLINED
		realEstate.DiscardReason = reason
		err = res.mailService.DiscardRealEstate(realEstate)
	}
	if CheckIfError(err) {
		return models.RealEstate{}
	}

	query = "INSERT INTO realestate (Id, Name, Type, Address, City, SquareFootage, NumberOfFloors, Picture, State, UserId, DiscardReason)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err = res.db.Exec(query, realEstate.Id, realEstate.Name, realEstate.Type, realEstate.Address, realEstate.City,
		realEstate.SquareFootage, realEstate.NumberOfFloors, realEstate.Picture, realEstate.State, realEstate.User,
		realEstate.DiscardReason)

	if CheckIfError(err) {
		return models.RealEstate{}
	}

	return realEstate

}

func (res *RealEstateServiceImpl) Add(estate models.RealEstate) models.RealEstate {
	estate.Id = res.generateId()

	// TODO: add some validation for pictures
	if estate.Address != "" && estate.City != "" && estate.SquareFootage != 0.0 && estate.NumberOfFloors != 0 {
		query := "INSERT INTO realestate (Id, Name, Type, Address, City, SquareFootage, NumberOfFloors, Picture, State, UserId, DiscardReason)" +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
		_, err := res.db.Exec(query, estate.Id, estate.Name, estate.Type, estate.Address, estate.City, estate.SquareFootage,
			estate.NumberOfFloors, estate.Picture, estate.State, estate.User, "")
		if CheckIfError(err) {
			return models.RealEstate{}
		}
		return estate
	}
	return models.RealEstate{}
}

func (res *RealEstateServiceImpl) generateId() int {
	id := 0
	estates := res.GetAll()

	for _, estate := range estates {
		if estate.Id > id {
			id = estate.Id
		}
	}
	return id + 1
}

func CheckIfError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return true
	}
	return false
}
