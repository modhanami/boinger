package services

import (
	"errors"
	"github.com/modhanami/boinger/models"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserCreationFailed = errors.New("failed to create user")
	ErrPasswordHashFailed = errors.New("failed to hash password")
)

type UserService interface {
	Create(username, password string) (models.UserModel, error)
	Exists(username string) bool
	GetById(id uint) (models.UserModel, error)
	GetByUsername(username string) (models.UserModel, error)
	GetByUid(uid string) (models.UserModel, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) Create(username, password string) (models.UserModel, error) {
	var user models.UserModel
	if exists := s.Exists(username); exists {
		return user, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, ErrPasswordHashFailed
	}

	user.Uid = ksuid.New().String()
	user.Username = username
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	if err := s.db.Create(&user).Error; err != nil {
		return user, ErrUserCreationFailed
	}

	return user, nil
}

func (s *userService) Exists(username string) bool {
	var user models.UserModel
	return s.db.Where("username = ?", username).First(&user).Error == nil
}

func (s *userService) GetById(id uint) (models.UserModel, error) {
	var user models.UserModel
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		} else {
			return user, err
		}
	}
	return user, nil
}

func (s *userService) GetByUsername(username string) (models.UserModel, error) {
	var user models.UserModel
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return user, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetByUid(uid string) (models.UserModel, error) {
	var user models.UserModel
	if err := s.db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return user, ErrUserNotFound
	}
	return user, nil
}
