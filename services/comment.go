package services

import (
	"errors"
	"fmt"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services/common"
	"github.com/modhanami/boinger/services/usercontext"
	"gorm.io/gorm"
)

var (
	ErrCommentCreationFailed = fmt.Errorf("failed to create comment")
	ErrCommentNotFound       = fmt.Errorf("comment not found")
)

type CommentService interface {
	Create(boingId, userId uint, text string) error
	Delete(user usercontext.UserContext, id uint) error
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

func (s *commentService) Delete(user usercontext.UserContext, id uint) error {
	userId := user.UserID()
	var comment models.Comment
	if err := s.db.Where("id = ?", id, userId).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return common.ErrUnexpectedDBError
	}

	if comment.UserID != userId {
		return common.ErrUserNotAuthorized
	}

	return s.db.Delete(&comment).Error
}
