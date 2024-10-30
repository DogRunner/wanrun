package dto

import "github.com/golang-jwt/jwt/v4"

// JWTのClaims
type AccountClaims struct {
	ID    string `json:"id"`
	JwtID string `json:"jwt_id"`
	jwt.RegisteredClaims
}
