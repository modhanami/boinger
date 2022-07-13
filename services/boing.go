package services

import (
	"errors"
	"github.com/modhanami/boinger/models"
	"github.com/segmentio/ksuid"
	"time"

	"gorm.io/gorm"
)

var (
	ErrBoingNotFound       = errors.New("boing not found")
	ErrBoingCreationFailed = errors.New("failed to create boing")
)

type BoingService interface {
	List() ([]models.BoingModel, error)
	Get(id uint) (models.BoingModel, error)
	Create(text string, userId uint) error
}

type boingService struct {
	db *gorm.DB
}

func NewBoingService(db *gorm.DB) BoingService {
	return &boingService{db: db}
}

func (s *boingService) List() ([]models.BoingModel, error) {
	var boings []models.BoingModel
	if err := s.db.Find(&boings).Error; err != nil {
		return nil, err
	}
	return boings, nil
}

func (s *boingService) Get(id uint) (models.BoingModel, error) {
	var boing models.BoingModel
	if err := s.db.First(&boing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return boing, ErrBoingNotFound
		} else {
			return boing, err
		}
	}
	return boing, nil
}

func (s *boingService) Create(text string, userId uint) error {
	uid := ksuid.New().String()

	var boing models.BoingModel
	boing.Uid = uid
	boing.Text = text
	boing.UserId = userId
	boing.CreatedAt = time.Now()

	if err := s.db.Create(&boing).Error; err != nil {
		return ErrBoingCreationFailed
	}
	return nil
}
