package services

import (
	"errors"
	"github.com/modhanami/boinger/models"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Authenticate(username, password string) (models.User, error)
	Register(username, password string) (models.User, error)
}

type authService struct {
	db               *gorm.DB
	userService      UserService
	userTokenService UserTokenService
	passwordHasher   PasswordHasher
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

func NewAuthService(db *gorm.DB, userService UserService, userTokenService UserTokenService, passwordHasher PasswordHasher) AuthService {
	return &authService{
		db:               db,
		userService:      userService,
		userTokenService: userTokenService,
		passwordHasher:   passwordHasher,
	}
}

func (s *authService) Authenticate(username, password string) (models.User, error) {
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		return models.User{}, err
	}

	err = s.passwordHasher.ComparePassword(user.Password, password)
	if err != nil {
		return models.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *authService) Register(username, password string) (models.User, error) {
	exists, err := s.userService.ExistsByUsername(username)
	if err != nil {
		return models.User{}, err
	}
	if exists {
		return models.User{}, ErrUserAlreadyExists
	}

	hashedPassword, err := s.passwordHasher.HashPassword(password)
	if err != nil {
		return models.User{}, ErrInvalidCredentials
	}

	user, err := s.userService.Create(username, hashedPassword)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
