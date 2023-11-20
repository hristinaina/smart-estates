package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Greška prilikom parsiranja tokena: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token nije validan")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Greška prilikom konverzije claims strukture")
	}

	return claims, nil
}
