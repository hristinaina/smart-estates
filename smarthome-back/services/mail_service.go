package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	SendVerificationMail(c *gin.Context)
	IsValidToken(tokeString string) bool
}

type MailServiceImpl struct{}

func NewMailService() MailService {
	return &MailServiceImpl{}
}

// request body
type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Picture  string `json:"picture" binding:"required"`
	Role     int    `json:"role" binding:"required"`
}

func (ms *MailServiceImpl) SendVerificationMail(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
	}
	expiration := time.Now().Add(time.Hour * 24)
	token, _ := generateToken(input.Email, input.Name, input.Surname, expiration)

	createVarificationMail(token)
}

func createVarificationMail(token string) {
	from := mail.NewEmail("SMART HOME SUPPORT", "savic.sv7.2020@uns.ac.rs")
	subject := "You're almost done! Activate your account now"
	to := mail.NewEmail("Example User", "anastasijas557@gmail.com") // todo ovo izmeni da bude dinamicki
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

func generateToken(email, name, surname string, expiration time.Time) (string, error) {
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
