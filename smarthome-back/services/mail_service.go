package services

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	SendVerificationMail()
}

type mailServiceImpl struct{}

func NewMailService() MailService {
	return &mailServiceImpl{}
}

func (ms *mailServiceImpl) SendVerificationMail() {
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
	subject := "Sign up verification"
	to := mail.NewEmail("Example User", "anastasijas557@gmail.com") // todo ovo izmeni da bude dinamicki
	plainTextContent := fmt.Sprintf("Click the following link to activate your account: %s", "http://localhost:3000/login")
	htmlContent := fmt.Sprintf(`<strong>Click the following link to activate your account:</strong> <a href="%s">%s</a>`, "http://localhost:3000/login", "http://localhost:3000/login")
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// todo promeni ovo da se uzima iz .env
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
