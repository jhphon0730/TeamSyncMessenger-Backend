package service

import (
	"database/sql"
)

type UserService interface {
	TestUser() string
}

type userService struct {
	DB *sql.DB
}

func NewUserService(DB *sql.DB) *userService {
	return &userService{
		DB: DB,
	}
}

func (us *userService) TestUser() string {
	return "대충 사용자 정보"
}
