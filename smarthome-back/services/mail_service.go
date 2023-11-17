package services

import (
	"database/sql"
	"fmt"
	"gopkg.in/mail.v2"
)

type MailService interface {
	Send(toAddress, subject, content string) error
	DiscardRealEstate(toAddress string) error
}

type MailServiceImpl struct {
	db *sql.DB
}

func NewMailService(db *sql.DB) MailService {
	return &MailServiceImpl{db: db}
}

func (ms *MailServiceImpl) Send(toAddress, subject, content string) error {
	message := mail.NewMessage()
	message.SetHeader("From", "kvucic6@gmail.com")
	message.SetHeader("To", toAddress)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", content)

	dialer := mail.NewDialer("smtp.gmail.com", 587, "kvucic6@gmail.com", "")
	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error2: ", err)
		return err
	}

	return nil
}

func (ms *MailServiceImpl) DiscardRealEstate(toAddress string) error {
	//senderName := "Smart Home Support Team"
	subject := "New Real Estate State"
	content := fmt.Sprintf("Hi, %s, we are very sorry, but your new real estate request has been rejected."+
		"Real Estate Name: %s. Reason for rejection: %s. Stay safe! Smart Home Support Team",
		"Katarina", "Kuca na moru", "Nemas para")

	err := ms.Send(toAddress, subject, content)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}
