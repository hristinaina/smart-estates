package services

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	// SendVerificationMail(c *gin.Context)
	IsValidToken(tokeString string) bool
	CreateVarificationMail(name, surname, token string)
	GenerateToken(email, name, surname string, expiration time.Time) (string, error)
}

type MailServiceImpl struct{}

func NewMailService() MailService {
	return &MailServiceImpl{}
}

func (ms *MailServiceImpl) CreateVarificationMail(name, surname, token string) {
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
	subject := "You're almost done! Activate your account now"
	to := mail.NewEmail(name+" "+surname, "anastasijas557@gmail.com")
	plainTextContent := fmt.Sprintf("Click the following link to activate your account: %s", "http://localhost:3000/activate?token="+token)
	htmlContent := fmt.Sprintf(`<strong>Click the following link to activate your account:</strong> <a href="%s">%s</a>`, "http://localhost:3000/activate?token="+token, "http://localhost:3000/activate?token="+token)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// todo promeni ovo da se uzima iz .env
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
	}
}

type Claims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	jwt.StandardClaims
}

var secretKey = []byte("JHAS43532fsandjaskndewui217362ebwdsa")

func (ms *MailServiceImpl) GenerateToken(email, name, surname string, expiration time.Time) (string, error) {
	claims := &Claims{
		Email:   email,
		Name:    name,
		Surname: surname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (ms *MailServiceImpl) IsValidToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	_, ok := token.Claims.(*Claims)
	if !ok {
		return false
	}

	return true
}
