package services

import (
	"errors"
	"github.com/modhanami/boinger/models"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

var (
	ErrBoingNotFound       = errors.New("boing not found")
	ErrBoingCreationFailed = errors.New("failed to create boing")
)

type BoingService interface {
	List() ([]models.Boing, error)
	GetById(id uint) (models.Boing, error)
	Create(text string, userId uint) error
}

type boingService struct {
	db *gorm.DB
}

func NewBoingService(db *gorm.DB) BoingService {
	return &boingService{db: db}
}

func (s *boingService) List() ([]models.Boing, error) {
	var boings []models.Boing
	if err := s.db.Find(&boings).Error; err != nil {
		return nil, ErrUnexpectedDBError
	}
	return boings, nil
}

func (s *boingService) GetById(id uint) (models.Boing, error) {
	var boing models.Boing
	if err := s.db.First(&boing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return boing, ErrBoingNotFound
		} else {
			return boing, ErrUnexpectedDBError
		}
	}
	return boing, nil
}

func (s *boingService) Create(text string, userId uint) error {
	uid := ksuid.New().String()

	boing := models.NewBoing(uid, text, userId)

	if err := s.db.Create(&boing).Error; err != nil {
		return ErrBoingCreationFailed
	}
	return nil
}
