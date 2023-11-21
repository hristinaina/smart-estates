package services

import (
	"database/sql"
	"errors"
	"fmt"
	"smarthome-back/enumerations"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"strings"
)

type RealEstateService interface {
	GetAll() ([]models.RealEstate, error)
	GetByUserId(id int) ([]models.RealEstate, error)
	Get(id int) (models.RealEstate, error)
	GetPending() ([]models.RealEstate, error)
	// ChangeState if state == 0 it is accepted, in opposite it is declined
	ChangeState(id int, state int, reason string) (models.RealEstate, error)
	Add(estate models.RealEstate) (models.RealEstate, error)
}

type RealEstateServiceImpl struct {
	db          *sql.DB
	repository  repositories.RealEstateRepository
	mailService MailService
}

func NewRealEstateService(db *sql.DB) RealEstateService {
	return &RealEstateServiceImpl{db: db, mailService: NewMailService(db), repository: *repositories.NewRealEstateRepository(db)}
}

func (res *RealEstateServiceImpl) GetAll() ([]models.RealEstate, error) {
	return res.repository.GetAll()
}

func (res *RealEstateServiceImpl) GetByUserId(userId int) ([]models.RealEstate, error) {
	return res.repository.GetByUserId(userId)
}

func (res *RealEstateServiceImpl) Get(id int) (models.RealEstate, error) {
	return res.repository.Get(id)
}

func (res *RealEstateServiceImpl) GetPending() ([]models.RealEstate, error) {
	return res.repository.GetPending()
}

func (res *RealEstateServiceImpl) ChangeState(id int, state int, reason string) (models.RealEstate, error) {
	realEstate, err := res.Get(id)
	if CheckIfError(err) {
		return models.RealEstate{}, err
	}

	if realEstate.State != enumerations.PENDING {
		return models.RealEstate{}, errors.New("only pending real estates can be accepted/declined")
	}

	if CheckIfError(err) {
		return models.RealEstate{}, err
	}
	if state == 0 {
		realEstate.State = enumerations.ACCEPTED
		err = res.mailService.ApproveRealEstate(realEstate)
	} else {
		realEstate.State = enumerations.DECLINED
		realEstate.DiscardReason = reason
		err = res.mailService.DiscardRealEstate(realEstate)
	}
	realEstate, err = res.repository.UpdateState(realEstate)

	return realEstate, err
}

func (res *RealEstateServiceImpl) Add(estate models.RealEstate) (models.RealEstate, error) {
	estate.Id = res.generateId()
	estate.Name = strings.Trim(estate.Name, " \t\n\r")
	estate.Address = strings.Trim(estate.Address, " \t\n\r")
	estate.City = strings.Trim(estate.City, " \t\n\r")

	// TODO: add some validation for pictures
	if estate.Name != "" && estate.Address != "" && estate.City != "" && estate.SquareFootage > 0.0 && estate.NumberOfFloors > 0 {
		estate, err := res.repository.Add(estate)
		return estate, err
	}
	return models.RealEstate{}, errors.New("invalid input")
}

func (res *RealEstateServiceImpl) generateId() int {
	id := 0
	estates, err := res.GetAll()

	if err != nil {
		return -1
	}

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
