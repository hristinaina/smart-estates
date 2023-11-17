package services

import (
	"database/sql"
	"fmt"
	"gopkg.in/gomail.v2"
)

type MailService interface {
	Send(toAddress, subject, content string) error
	DiscardRealEstate(toAddress string) error
	ApproveRealEstate(toAddress string) error
}

type MailServiceImpl struct {
	db *sql.DB
}

func NewMailService(db *sql.DB) MailService {
	return &MailServiceImpl{db: db}
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

func (ms *MailServiceImpl) ApproveRealEstate(toAddress string) error {
	subject := "New Real Estate State"
	content := fmt.Sprintf("<h1>Hi %s,</h1> <br/> We have some good news. Your real estate request has been approved!"+
		"<br/>Real Estate Name: %s. <br/> Smart Home Support Team", "Katarina", "Kuca na moru")

	err := ms.Send(toAddress, subject, content)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

func (ms *MailServiceImpl) DiscardRealEstate(toAddress string) error {
	subject := "New Real Estate State"
	content := fmt.Sprintf("<h1>Hi %s,</h1> <br/> We are very sorry, but your new real estate request has been rejected."+
		"<br/>Real Estate Name: %s. <br/>Reason for rejection: %s.<br/> Stay safe!<br/> Smart Home Support Team",
		"Katarina", "Kuca na moru", "Nemas para")

	err := ms.Send(toAddress, subject, content)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}
