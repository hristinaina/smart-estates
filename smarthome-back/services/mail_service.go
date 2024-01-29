package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"smarthome-back/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	IsValidToken(tokeString string) bool
	CreateVarificationMail(email, name, surname, token string)
	GenerateToken(email string, expiration time.Time) (string, error)
	CreateAdminLoginRequest(name, surname, email, password string)
	SendVerifyEmail(email, token string)
	DiscardRealEstate(estate models.RealEstate, email, name string)
	ApproveRealEstate(estate models.RealEstate, email, name string)
}

type MailServiceImpl struct {
	db *sql.DB
}

func NewMailService(db *sql.DB) MailService {
	return &MailServiceImpl{db: db}
}

func readFromEnvFile() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv("SENDGRID_API_KEY")
}

func readSenderEmailFromJson() string {
	file, err := os.Open("properties.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	// Parsiranje JSON fajla u mapu
	var config map[string]interface{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	email, ok := config["senderEmail"].(string)
	if !ok {
		return ""
	}

	return email
}

// Funkcija za proveru da li domen podržava HTML
func isDomainSupportingHTML(domain string) bool {
	domainsSupportingHTML := map[string]bool{
		"gmail.com":   true,
		"yahoo.com":   true,
		"outlook.com": true,
		"hotmail.com": true,
		"aol.com":     true,
	}
	return domainsSupportingHTML[domain]
}

func (ms *MailServiceImpl) CreateVarificationMail(email, name, surname, token string) {
	from := mail.NewEmail("SMART HOME SUPPORT", readSenderEmailFromJson())
	subject := "You're almost done! Activate your account now"
	to := mail.NewEmail(name+" "+surname, email)
	plainTextContent := fmt.Sprintf("Click the following link to activate your account: %s", "http://localhost:3000/activate?token="+token)
	htmlContent := fmt.Sprintf(`<strong>Click the following link to activate your account:</strong> <a href="%s">%s</a>`, "http://localhost:3000/activate?token="+token, "http://localhost:3000/activate?token="+token)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	if isDomainSupportingHTML(strings.Split(email, "@")[1]) {
		message.SetTemplateID("d-aa0a609c711d4d97ba4dbd99a943bd3f")
		message.Personalizations[0].SetDynamicTemplateData("user_name", name)
		message.Personalizations[0].SetDynamicTemplateData("link", "http://localhost:3000/activate?token="+token)
	}

	client := sendgrid.NewSendClient(readFromEnvFile())
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var secretKey = []byte("JHAS43532fsandjaskndewui217362ebwdsa")

func (ms *MailServiceImpl) GenerateToken(email string, expiration time.Time) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
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

func (ms *MailServiceImpl) CreateAdminLoginRequest(name, surname, email, password string) {
	from := mail.NewEmail("SMART HOME SUPPORT", readSenderEmailFromJson())
	subject := "Join our team!"
	to := mail.NewEmail(name+" "+surname, email)
	plainTextContent := fmt.Sprintf("Your login credentials to our system are:  E=mail: %s  Password: %s", email, password)
	htmlContent := fmt.Sprintf(`<strong>Your login credentials to our system are:</strong> <br/> E-mail: <strong>%s</strong><br/> Password: <strong>%s</strong>`, email, password)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	if isDomainSupportingHTML(strings.Split(email, "@")[1]) {
		message.SetTemplateID("d-3651b37ad1f94cd2b5a175329b474fba")
		message.Personalizations[0].SetDynamicTemplateData("admin_name", name)
		message.Personalizations[0].SetDynamicTemplateData("admin_email", email)
		message.Personalizations[0].SetDynamicTemplateData("admin_password", password)
		message.Personalizations[0].SetDynamicTemplateData("link", "http://localhost:3000")
	}

	client := sendgrid.NewSendClient(readFromEnvFile())
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
	}
}

func (ms *MailServiceImpl) SendVerifyEmail(email, token string) {
	from := mail.NewEmail("SMART HOME SUPPORT", readSenderEmailFromJson())
	subject := "Verify email!"
	to := mail.NewEmail("", email)
	plainTextContent := fmt.Sprintf("Click the following link to verify your email: %s", "http://localhost:3000/reset-password?token="+token)
	htmlContent := fmt.Sprintf(`<strong>Click the following link to verify your email and reset your password:</strong> <a href="%s">%s</a>`, "http://localhost:3000/reset-password?token="+token, "http://localhost:3000/reset-password?token="+token)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	if isDomainSupportingHTML(strings.Split(email, "@")[1]) {
		message.SetTemplateID("d-8b4418f13eba4d86982b3ef6fdff7080")
		message.Personalizations[0].SetDynamicTemplateData("link", "http://localhost:3000/reset-password?token="+token)
	}

	client := sendgrid.NewSendClient(readFromEnvFile())
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
	}
}

func (ms *MailServiceImpl) ApproveRealEstate(estate models.RealEstate, email, name string) {
	from := mail.NewEmail("SMART HOME SUPPORT", readSenderEmailFromJson())
	subject := "Approved real estate"
	to := mail.NewEmail(name, email)
	plainTextContent := fmt.Sprintf("Hi %s! We have some good news. Your real estate request has been approved! Real Estate Name: %s.", name, estate.Name)
	htmlContent := fmt.Sprintf("Hi %s! We have some good news. Your real estate request has been approved! Real Estate Name: %s.", name, estate.Name)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	if isDomainSupportingHTML(strings.Split(email, "@")[1]) {
		message.SetTemplateID("d-be3dd5200af74140a641fc12b9d2f710")
		message.Personalizations[0].SetDynamicTemplateData("user_name", name)
		message.Personalizations[0].SetDynamicTemplateData("real_estate_name", estate.Name)
	}

	client := sendgrid.NewSendClient(readFromEnvFile())
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
	}
}

func (ms *MailServiceImpl) DiscardRealEstate(estate models.RealEstate, email, name string) {
	from := mail.NewEmail("SMART HOME SUPPORT", readSenderEmailFromJson())
	subject := "Discarded real estate"
	to := mail.NewEmail(name, email)
	plainTextContent := fmt.Sprintf("Hi %s, We are very sorry, but Your new real estate request has been rejected. Real Estate Name: %s. Reason for rejection: %s.", name, estate.Name, estate.DiscardReason)
	htmlContent := fmt.Sprintf("Hi %s, We are very sorry, but Your new real estate request has been rejected. Real Estate Name: %s. Reason for rejection: %s.", name, estate.Name, estate.DiscardReason)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	if isDomainSupportingHTML(strings.Split(email, "@")[1]) {
		message.SetTemplateID("d-62230654a73b491ab95606c15f96734f")
		message.Personalizations[0].SetDynamicTemplateData("user_name", name)
		message.Personalizations[0].SetDynamicTemplateData("real_estate_name", estate.Name)
		message.Personalizations[0].SetDynamicTemplateData("reason", estate.DiscardReason)
	}

	client := sendgrid.NewSendClient(readFromEnvFile())
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
	}
}
