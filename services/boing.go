package services

import (
	"errors"
	"github.com/modhanami/boinger/logger"
	"github.com/modhanami/boinger/models"
	"gorm.io/gorm"
)

var (
	ErrBoingNotFound       = errors.New("boing not found")
	ErrBoingCreationFailed = errors.New("failed to create boing")
)

type BoingService interface {
	List() ([]*models.Boing, error)
	GetById(id uint) (*models.Boing, error)
	Create(text string, userId uint) error
}

type boingService struct {
	db  *gorm.DB
	log logger.Logger
}

func NewBoingService(db *gorm.DB, logger logger.Logger) BoingService {
	return &boingService{db: db, log: logger}
}

func (s *boingService) List() ([]*models.Boing, error) {
	var boings []*models.Boing
	if err := s.db.Find(&boings).Error; err != nil {
		return nil, ErrUnexpectedDBError
	}
	return boings, nil
}

func (s *boingService) GetById(id uint) (*models.Boing, error) {
	var boing models.Boing
	l := s.log.With("context", "boingService.GetById", "boingId", id)
	if err := s.db.First(&boing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			l.Info("boing not found")
			return nil, ErrBoingNotFound
		} else {
			l.Error("unexpected db error", "error", err)
			return nil, ErrUnexpectedDBError
		}
	}

	l.Info("boing found", "boingId", boing.ID)
	return &boing, nil
}

func (s *boingService) Create(text string, userId uint) error {
	boing := models.NewBoing(text, userId)
	l := s.log.With("context", "boingService.Create")

	if err := s.db.Create(&boing).Error; err != nil {
		l.Error("failed to create boing", "error", err)
		return ErrBoingCreationFailed
	}

	l.Info("boing created", "boingId", boing.ID)
	return nil
}
