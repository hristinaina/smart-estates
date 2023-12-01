package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"smarthome-back/models"
)

type UserRepository interface {
	GetAll() []models.User
	SaveUser(user models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id int) (*models.User, error)
	ResetSuperAdminPassword(password string, id int) error
	EditSuperAdmin(name, surname, email string) error
	ResetPassword(email, password string) error
}

type UserRepositoryImpl struct {
	db *sql.DB
}

var ErrUserNotFound = errors.New("User not found")

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (res *UserRepositoryImpl) GetAll() []models.User {
	query := "SELECT * FROM user"
	rows, err := res.db.Query(query)
	if CheckIfError(err) {
		return nil
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var (
			user models.User
		)

		if err := rows.Scan(&user.Id, &user.Email, &user.Password,
			&user.Name, &user.Surname, &user.Role, &user.IsLogin); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.User{}
		}
		users = append(users, user)
	}

	return users
}

func (res *UserRepositoryImpl) SaveUser(user models.User) error {
	user.Id = res.generateId()

	query := "INSERT INTO User (Id, Email, Password, Name, Surname, Role, IsLogin)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?);"
	_, err := res.db.Exec(query, user.Id, user.Email, user.Password, user.Name, user.Surname, user.Role, user.IsLogin)
	if CheckIfError(err) {
		return fmt.Errorf("Failed to save user: %v", err)

	}
	return nil
}

func (res *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM user WHERE email = ?"
	row := res.db.QueryRow(query, email)

	var user models.User

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Role, &user.IsLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (res *UserRepositoryImpl) GetUserById(id int) (*models.User, error) {
	query := "SELECT * FROM user WHERE id = ?"
	row := res.db.QueryRow(query, id)

	var user models.User

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Role, &user.IsLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (res *UserRepositoryImpl) ResetSuperAdminPassword(password string, id int) error {
	updateStatement := "UPDATE user SET Password=?, IsLogin=? WHERE Id=?"

	_, err := res.db.Exec(updateStatement, password, true, id)
	return err
}

func (res *UserRepositoryImpl) EditSuperAdmin(name, surname, email string) error {
	updateStatement := "UPDATE user SET Name=?, Surname=? WHERE email=?"

	_, err := res.db.Exec(updateStatement, name, surname, email)
	return err
}

func (res *UserRepositoryImpl) ResetPassword(email, password string) error {
	updateStatement := "UPDATE user SET Password=? WHERE Email=?"

	_, err := res.db.Exec(updateStatement, password, email)
	return err
}

func (res *UserRepositoryImpl) generateId() int {
	id := 0
	users := res.GetAll()

	for _, user := range users {
		if user.Id > id {
			id = user.Id
		}
	}
	return id + 1
}

func CheckIfError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return true
	}
	return false
}
