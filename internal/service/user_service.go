package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"projetoapi/internal/domain"
	"projetoapi/internal/dto"
	"projetoapi/internal/repository"
)

// UserService defines the interface for user business logic
type UserService interface {
	Login(req dto.LoginRequest) (*domain.User, error)
	Register(req dto.RegisterRequest) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Login(req dto.LoginRequest) (*domain.User, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !checkPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userService) Register(req dto.RegisterRequest) (*domain.User, error) {
	exists, err := s.repo.UsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username: req.Username,
		Password: hash,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByUsername(username string) (*domain.User, error) {
	return s.repo.FindByUsername(username)
}

// HashPassword creates a bcrypt hash of the password (exported for seeding)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// HashPasswordForSeeding creates a bcrypt hash for database seeding
func HashPasswordForSeeding(password string) string {
	hash, _ := HashPassword(password)
	return hash
}

// checkPasswordHash compares a password with its hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
