package service

import (
	"TeamSyncMessenger-Backend/DTO"
	"TeamSyncMessenger-Backend/model"
	"database/sql"
	"log"
)

func CreateUsersTable(DB *sql.DB) {
	_, err := DB.Exec(`
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT,
				email TEXT,
				password TEXT,
				last_message TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

type UserService interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (model.User, error)
	CreateUser(registerUserDTO DTO.RegisterUserDTO) (DTO.RegisterUserDTO, error)
}

type userService struct {
	DB *sql.DB
}

func NewUserService(DB *sql.DB) *userService {
	CreateUsersTable(DB)
	return &userService{
		DB: DB,
	}
}

func (us *userService) GetUsers() ([]model.User, error) {
	rows, err := us.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.LastMessage, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) GetUserByUsername(username string) (model.User, error) {
	var validUser model.User

	err := us.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&validUser.ID)
	if err != nil {
		return validUser, err
	}

	return validUser, nil
}

func (us *userService) CreateUser(registerUserDTO DTO.RegisterUserDTO) (DTO.RegisterUserDTO, error) {
	_, err := us.DB.Exec("INSERT INTO users (username, email, password) VALUES(?, ?, ?)", &registerUserDTO.Username, &registerUserDTO.Email, &registerUserDTO.Password)
	if err != nil {
		return DTO.RegisterUserDTO{}, err
	}

	return registerUserDTO, nil
}
