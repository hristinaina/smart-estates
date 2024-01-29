package services

import (
	"database/sql"
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
	"gopkg.in/gomail.v2"
)

type MailService interface {
	IsValidToken(tokeString string) bool
	CreateVarificationMail(email, name, surname, token string)
	GenerateToken(email string, expiration time.Time) (string, error)
	CreateAdminLoginRequest(name, surname, email, password string)
	SendVerifyEmail(email, token string)
	Send(toAddress, subject, content string) error
	DiscardRealEstate(estate models.RealEstate) error
	ApproveRealEstate(estate models.RealEstate) error
	PermissionMail(email, name, owner, realEstate, token string)
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
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
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
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
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
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
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
	// TODO : use this after A&A is implemented
	// user := ms.service.GetUser(estate.User)
	// toAddress := user.Email
	toAddress := "kacorinav@gmail.com"
	subject := "New Real Estate State"
	// TODO : change parameter name in content
	content := fmt.Sprintf("<h1>Hi %s,</h1> <br/> We have some good news. Your real estate request has been approved!"+
		"<br/>Real Estate Name: %s. <br/> Smart Home Support Team", "User", estate.Name)

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
	// TODO : change parameter name in content
	content := fmt.Sprintf("<h1>Hi %s,</h1> <br/> We are very sorry, but your new real estate request has been rejected."+
		"<br/>Real Estate Name: %s. <br/>Reason for rejection: %s.<br/> Stay safe!<br/> Smart Home Support Team",
		"User", estate.Name, estate.DiscardReason)

	err := ms.Send(toAddress, subject, content)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

func (ms *MailServiceImpl) PermissionMail(email, name, owner, realEstate, token string) {
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
	subject := "Granted permissions"
	to := mail.NewEmail("", email)
	plainTextContent := fmt.Sprintf("Click the following link to verify your email: %s", "http://localhost:3000/activate-permission?token="+token)
	htmlContent := fmt.Sprintf(`<strong>Click the following link to verify your email and reset your password:</strong> <a href="%s">%s</a>`, "http://localhost:3000/activate-permission?token="+token, "http://localhost:3000/activate-permission?token="+token)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	if isDomainSupportingHTML(strings.Split(email, "@")[1]) {
		message.SetTemplateID("d-227ab0e9abe64dc2b1a847bc32ff4113")
		message.Personalizations[0].SetDynamicTemplateData("user_name", name)
		message.Personalizations[0].SetDynamicTemplateData("owner_name", owner)
		message.Personalizations[0].SetDynamicTemplateData("real_estate", realEstate)
		message.Personalizations[0].SetDynamicTemplateData("link", "http://localhost:3000/activate-permission?token="+token)
	}

	client := sendgrid.NewSendClient(readFromEnvFile())
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
	}
}
