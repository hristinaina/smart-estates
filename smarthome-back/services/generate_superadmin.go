package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"smarthome-back/cache"
	"smarthome-back/enumerations"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type GenerateSuperadmin struct {
	repo repositories.UserRepository
}

func NewGenerateSuperAdmin(db *sql.DB, cacheService cache.CacheService) GenerateSuperadmin {
	return GenerateSuperadmin{repo: repositories.NewUserRepository(db, &cacheService)}
}

func (gs GenerateSuperadmin) GenerateSuperadmin() {
	// check if superadmin already exist
	_, err := gs.repo.GetUserByEmail("admin")
	if err != nil {

		// save password in file
		generatePassword := gs.GenerateRandomPassword(30)

		err = writePasswordToFile(generatePassword)
		if err != nil {
			panic("Failed to write initial password to file")
		}

		// save super admin account
		newSuperadmin := models.User{Email: "admin", Password: gs.HashPassword(generatePassword), Role: enumerations.SUPERADMIN, IsLogin: false}
		gs.repo.SaveUser(newSuperadmin)
	}
}

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	specialChars     = "!@#$%^*()-=_+[]{}|;:,.?"
)

func (gs GenerateSuperadmin) GenerateRandomPassword(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	allChars := lowercaseLetters + uppercaseLetters + digits + specialChars

	password := make([]byte, length)
	password[0] = lowercaseLetters[random.Intn(len(lowercaseLetters))]
	password[1] = uppercaseLetters[random.Intn(len(uppercaseLetters))]
	password[2] = digits[random.Intn(len(digits))]
	password[3] = specialChars[random.Intn(len(specialChars))]

	for i := 4; i < length; i++ {
		password[i] = allChars[random.Intn(len(allChars))]
	}

	// password shuffle
	random.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}

func (gs GenerateSuperadmin) HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic("Failed to hash password")
	}
	return string(hashedPassword)
}

type Superadmin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func writePasswordToFile(password string) error {
	admin := []Superadmin{{Email: "admin", Password: password}}

	jsonData, err := json.MarshalIndent(admin, "", "    ")
	if err != nil {
		fmt.Println("Error while converting in JSON:", err)
		return err
	}

	err = os.WriteFile("admin.json", jsonData, 0400)
	if err != nil {
		fmt.Println("Write file error:", err)
		return err
	}

	fmt.Println("\nFOR SUPERADMIN EMAIL AND PASSWORD LOOK smart-home/smarthome-back/admin.json FILE\n")
	return nil
}
