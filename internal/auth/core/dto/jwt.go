package dto

import "github.com/golang-jwt/jwt/v4"

// JWTのClaims
type AccountClaims struct {
	ID  string `json:"id"`
	JTI string `json:"jti"`
	jwt.RegisteredClaims
}
