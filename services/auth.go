package services

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(username, password string) (UserToken, UserClaims, error)
	Register(username, password string) (bool, error)
}

type authService struct {
	db               *gorm.DB
	userService      UserService
	userTokenService UserTokenService
}

func NewAuthService(db *gorm.DB, userService UserService, userTokenService UserTokenService) AuthService {
	return &authService{db: db, userService: userService, userTokenService: userTokenService}
}

func (s *authService) Login(username, password string) (UserToken, UserClaims, error) {
	user, err := s.userService.GetByUsername(username)
	var userToken UserToken

	if err != nil {
		return userToken, UserClaims{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return userToken, UserClaims{}, err
	}

	userToken, claims, err := s.userTokenService.Create(&user, CreateOptions{})

	return userToken, claims, err
}

func (s *authService) Register(username, password string) (bool, error) {
	if s.userService.Exists(username) {
		return false, ErrUserAlreadyExists
	}

	if _, err := s.userService.Create(username, password); err != nil {
		return false, err
	}

	return true, nil
}
