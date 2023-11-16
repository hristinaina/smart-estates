package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"smarthome-back/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Middleware struct {
	repo repositories.UserRepository
}

func NewMiddleware(db *sql.DB) Middleware {
	return Middleware{repo: repositories.NewUserRepository(db)}
}

func (mw Middleware) RequireAuth(c *gin.Context) {
	// get the cookie off request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// decode/validate it

	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// find the user with token sub
		sub, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		subInt := int(sub)
		user, err := mw.repo.GetUserById(subInt)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// attach to request
		c.Set("user", user)

		// continue
		c.Next()

		fmt.Println(claims["sub"], claims["exp"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
