package services

import (
	"fmt"
	"github.com/modhanami/boinger/models"
	"gorm.io/gorm"
)

var (
	ErrCommentCreationFailed = fmt.Errorf("failed to create comment")
)

type CommentService interface {
	Create(boingId, userId uint, text string) error
}

type commentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) CommentService {
	return &commentService{db: db}
}

func (s *commentService) Create(boingId, userId uint, text string) error {
	comment := models.NewComment(boingId, userId, text)
	fmt.Println(boingId, userId, text)

	if err := s.db.Create(comment).Error; err != nil {
		return ErrCommentCreationFailed
	}

	return nil
}
