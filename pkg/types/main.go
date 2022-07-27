package types

import (
	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Picture string
	Email   string
	Name    string
	jwt.StandardClaims
}
