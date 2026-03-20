package jwt

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"projetoapi/internal/domain"
)

// Service defines the interface for JWT operations
type Service interface {
	GenerateToken(username string) (string, error)
	ValidateToken(c *gin.Context) bool
	GetUsernameFromToken(c *gin.Context) (string, error)
	GetTokenFromRequest(c *gin.Context) (string, error)
}

// service implements JWT Service
type service struct {
	secretKey []byte
}

// NewService creates a new JWT service
func NewService(secretKeyPath string) (Service, error) {
	b, err := ioutil.ReadFile(secretKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret key: %w", err)
	}
	return &service{secretKey: []byte(strings.TrimSpace(string(b)))}, nil
}

// NewServiceWithKey creates a new JWT service with a direct key
func NewServiceWithKey(key []byte) Service {
	return &service{secretKey: key}
}

func (s *service) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &domain.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *service) ValidateToken(c *gin.Context) bool {
	token, err := s.GetTokenFromRequest(c)
	if err != nil {
		return false
	}

	claims := &domain.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return false
	}

	if tkn != nil {
		return tkn.Valid
	}
	return false
}

func (s *service) GetUsernameFromToken(c *gin.Context) (string, error) {
	tokenString, err := s.GetTokenFromRequest(c)
	if err != nil {
		return "", err
	}

	claims := &domain.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", fmt.Errorf("invalid token")
}

func (s *service) GetTokenFromRequest(c *gin.Context) (string, error) {
	reqToken := c.Request.Header.Get("Authorization")

	if strings.Contains(reqToken, "Bearer") {
		if strings.TrimSpace(reqToken) == "" {
			return "", fmt.Errorf("empty authorization header")
		}
		splitToken := strings.Split(reqToken, "Bearer")
		return strings.TrimSpace(splitToken[1]), nil
	}

	return strings.TrimSpace(reqToken), nil
}
