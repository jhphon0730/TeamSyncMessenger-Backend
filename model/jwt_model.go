package model

import "gopkg.in/dgrijalva/jwt-go.v3"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var JwtKey = []byte("some_jwt_token")
