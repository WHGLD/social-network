package common

import "github.com/golang-jwt/jwt/v4"

var JwtKey = []byte("secret-key-for-s-network")

type Claims struct {
	UserID           string `json:"user_id"`
	RegisteredClaims jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
