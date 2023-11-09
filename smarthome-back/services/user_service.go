package services

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"smarthome-back/models"
)

type UserService interface {
	ListUsers() []models.User
	GetUser(id string) (models.User, error)
	TestGetMethod()
}

type UserServiceImpl struct {
	db       *gorm.DB
	database *sql.DB
}

func NewUserService(db *gorm.DB, database *sql.DB) UserService {
	return &UserServiceImpl{db: db, database: database}
}

func (us *UserServiceImpl) ListUsers() []models.User {
	var users = []models.User{
		{Id: 1, Name: "Blue Train"},
		{Id: 2, Name: "Jeru"},
		{Id: 3, Name: "Sarah Vaughan and Clifford Brown"},
	}
	// var users []models.User
	// us.db.Find(&users)
	return users
}

func (us *UserServiceImpl) GetUser(id string) (models.User, error) {
	var user models.User
	result := us.db.First(&user, id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (us *UserServiceImpl) TestGetMethod() {

	rows, err := us.database.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("Error1: ", err.Error())
		return
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var (
			id       int
			email    string
			password string
			role     int
		)
		if err := rows.Scan(&id, &email, &password, &role); err != nil {
			fmt.Println("Error0: ", err.Error())
			return
		}
		// You can process the data or return it as JSON
		fmt.Println(id, email, password, role)
	}

	fmt.Println("Data fetched successfully!")

}
