package types

import (
	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Picture string
	Email   string
	Name    string
	Id      string
	jwt.StandardClaims
}

type Microservice struct {
	Url  string
	Name string
}

type Payload struct {
	data []byte
}
