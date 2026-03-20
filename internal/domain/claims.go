package domain

import "github.com/dgrijalva/jwt-go"

// Claims represents JWT claims with embedded StandardClaims
type Claims struct {
	Username           string `json:"username"`
	jwt.StandardClaims `swaggerignore:"true"`
}
