package service

import (
	"TeamSyncMessenger-Backend/DTO"
	"TeamSyncMessenger-Backend/model"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
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
	ComparePasswords(hashedPassword, password string) error
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

// 사용자가 제공한 비밀번호를 bcrypt 해시를 사용하여 암호화
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// 사용자가 제공한 비밀번호와 저장된 bcrypt 해시를 비교하는 함수
func (us *userService) ComparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func (us *userService) CreateUser(registerUserDTO DTO.RegisterUserDTO) (DTO.RegisterUserDTO, error) {
	hashedPassword, err := hashPassword(registerUserDTO.Password)
	if err != nil {
		return DTO.RegisterUserDTO{}, err
	}
	registerUserDTO.Password = hashedPassword

	_, err = us.DB.Exec("INSERT INTO users (username, email, password) VALUES(?, ?, ?)", &registerUserDTO.Username, &registerUserDTO.Email, &registerUserDTO.Password)
	if err != nil {
		return DTO.RegisterUserDTO{}, err
	}

	return registerUserDTO, nil
}
