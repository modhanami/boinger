package services

import (
	"errors"
	"github.com/modhanami/boinger/logger"
	"github.com/modhanami/boinger/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserCreationFailed = errors.New("failed to create user")
)

type UserService interface {
	Create(user *models.User) (*models.User, error)
	ExistsByUsername(username string) (bool, error)
	GetById(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
}

type userService struct {
	db  *gorm.DB
	log logger.Logger
}

func NewUserService(db *gorm.DB, log logger.Logger) UserService {
	return &userService{db: db, log: log}
}

func (s *userService) Create(user *models.User) (*models.User, error) {
	if err := s.db.Create(&user).Error; err != nil {
		s.log.Error("failed to create user", "username", user.Username, "error", err)
		return nil, ErrUserCreationFailed
	}

	s.log.Info("user created", "username", user.Username)
	return user, nil
}

func (s *userService) ExistsByUsername(username string) (bool, error) { // TODO: unexport this
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (s *userService) GetById(id uint) (*models.User, error) {
	s.log.Info("getting user by username", "username", id)
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Info("user not found", "username", id)
			return &user, ErrUserNotFound
		} else {
			return &user, err
		}
	}
	return &user, nil
}

func (s *userService) GetByUsername(username string) (*models.User, error) {
	s.log.Info("getting user by username", "username", username)
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Info("user not found", "username", username)
			return &user, ErrUserNotFound
		} else {
			return &user, err
		}
	}
	return &user, nil
}
