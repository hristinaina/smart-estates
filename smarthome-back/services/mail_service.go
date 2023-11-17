package services

import (
	"database/sql"
	"fmt"
	"gopkg.in/gomail.v2"
	"smarthome-back/models"
)

type MailService interface {
	Send(toAddress, subject, content string) error
	DiscardRealEstate(estate models.RealEstate) error
	ApproveRealEstate(estate models.RealEstate) error
}

type MailServiceImpl struct {
	db      *sql.DB
	service UserService
}

func NewMailService(db *sql.DB) MailService {
	return &MailServiceImpl{db: db, service: NewUserService(db)}
}

func (ms *MailServiceImpl) Send(to, subject, body string) error {
	appPass, err := NewConfigService().GetAppPassword("config/config.json")

	if err != nil {
		return err
	}

	var emailConfig = map[string]string{
		"host":     "smtp.gmail.com",
		"port":     "587",
		"username": "kacorinav@gmail.com",
		"password": appPass,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig["username"])
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(emailConfig["host"], 587, emailConfig["username"], emailConfig["password"])

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

func (ms *MailServiceImpl) ApproveRealEstate(estate models.RealEstate) error {
	// TODO : use this after user is implemented
	// user := ms.service.GetUser(estate.User)
	// toAddress := user.Email
	toAddress := "kacorinav@gmail.com"
	subject := "New Real Estate State"
	// TODO : change parameter in content
	content := fmt.Sprintf("<h1>Hi %s,</h1> <br/> We have some good news. Your real estate request has been approved!"+
		"<br/>Real Estate Name: %s. <br/> Smart Home Support Team", "Katarina", estate.Name)

	err := ms.Send(toAddress, subject, content)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

func (ms *MailServiceImpl) DiscardRealEstate(estate models.RealEstate) error {
	// TODO : use this after user is implemented
	// user := ms.service.GetUser(estate.User)
	// toAddress := user.Email
	toAddress := "kacorinav@gmail.com"
	subject := "New Real Estate State"
	// TODO : change parameter in content
	content := fmt.Sprintf("<h1>Hi %s,</h1> <br/> We are very sorry, but your new real estate request has been rejected."+
		"<br/>Real Estate Name: %s. <br/>Reason for rejection: %s.<br/> Stay safe!<br/> Smart Home Support Team",
		"Katarina", estate.Name, estate.DiscardReason)

	err := ms.Send(toAddress, subject, content)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}
