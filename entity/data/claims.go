package data

import "github.com/golang-jwt/jwt"

var JwtKey = []byte("library-ap!-ilham")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
