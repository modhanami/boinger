package services

import (
	"errors"
	"github.com/modhanami/boinger/log"
	"github.com/modhanami/boinger/models"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserCreationFailed = errors.New("failed to create user")
)

type UserService interface {
	Create(username, password string) (models.User, error)
	ExistsByUsername(username string) (bool, error)
	GetById(id uint) (models.User, error)
	GetByUsername(username string) (models.User, error)
	GetByUid(uid string) (models.User, error)
	GetByUids(uids []string) ([]models.User, error)
}

type userService struct {
	db  *gorm.DB
	log log.Interface
}

func NewUserService(db *gorm.DB, log log.Interface) UserService {
	return &userService{db: db, log: log}
}

func (s *userService) Create(username, hashedPassword string) (models.User, error) {
	var user models.User
	user.Uid = ksuid.New().String()
	user.Username = username
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	if err := s.db.Create(&user).Error; err != nil {
		s.log.Error("failed to create user", "username", username, "error", err)
		return user, ErrUserCreationFailed
	}

	s.log.Info("user created", "username", username)
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

func (s *userService) GetById(id uint) (models.User, error) {
	s.log.Info("getting user by id", "id", id)
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Info("user not found", "id", id)
			return user, ErrUserNotFound
		} else {
			return user, err
		}
	}
	return user, nil
}

func (s *userService) GetByUsername(username string) (models.User, error) {
	s.log.Info("getting user by username", "username", username)
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Info("user not found", "username", username)
			return user, ErrUserNotFound
		} else {
			return user, err
		}
	}
	return user, nil
}

func (s *userService) GetByUid(uid string) (models.User, error) {
	s.log.Info("getting user by uid", "uid", uid)
	var user models.User
	if err := s.db.Where("uid = ?", uid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Info("user not found", "uid", uid)
			return user, ErrUserNotFound
		} else {
			return user, err
		}
	}
	return user, nil
}

func (s *userService) GetByUids(uids []string) ([]models.User, error) {
	s.log.Info("getting users by uids", "uids", uids)
	var users []models.User
	if err := s.db.Where("uid IN (?)", uids).Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
