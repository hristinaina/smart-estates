package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"smarthome-back/models"
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
	GenerateToken(email, name, surname string, expiration time.Time) (string, error)
	CreateAdminLoginRequest(name, surname, email, password string)
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

func readFromEnvFile() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv("SENDGRID_API_KEY")
}

func (ms *MailServiceImpl) CreateVarificationMail(email, name, surname, token string) {
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
	subject := "You're almost done! Activate your account now"
	to := mail.NewEmail(name+" "+surname, email)

	plainTextContent := fmt.Sprintf("Click the following link to activate your account: %s", "http://localhost:3000/activate?token="+token)
	htmlContent := fmt.Sprintf(`<strong>Click the following link to activate your account:</strong> <a href="%s">%s</a>`, "http://localhost:3000/activate?token="+token, "http://localhost:3000/activate?token="+token)
	m := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	m.SetTemplateID("d-aa0a609c711d4d97ba4dbd99a943bd3f")
	m.Personalizations[0].SetDynamicTemplateData("user_name", name)
	m.Personalizations[0].SetDynamicTemplateData("link", "http://localhost:3000/activate?token="+token)

	client := sendgrid.NewSendClient(readFromEnvFile())
	response, err := client.Send(m)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
	}
}

type Claims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	jwt.RegisteredClaims
}

var secretKey = []byte("JHAS43532fsandjaskndewui217362ebwdsa")

func (ms *MailServiceImpl) GenerateToken(email, name, surname string, expiration time.Time) (string, error) {
	claims := &Claims{
		Email:   email,
		Name:    name,
		Surname: surname,
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
