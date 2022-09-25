package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints/response"
	"github.com/modhanami/boinger/endpoints/utils"
	"github.com/modhanami/boinger/logger"
	"github.com/modhanami/boinger/services"
	"net/http"
	"strconv"
)

type commentHandler struct {
	s   services.CommentService
	log logger.Logger
}

func NewCommentHandler(s services.CommentService, log logger.Logger) *commentHandler {
	return &commentHandler{s: s, log: log}
}

type CreateCommentRequest struct {
	BoingId uint
	Text    string
}

func (h *commentHandler) Create(c *gin.Context) {
	log := h.log.With("context", "commentHandler.Create")
	userClaims := utils.GetUserClaimsFromContext(c)
	if userClaims == nil {
		h.log.Error("user claims not found in context")
		c.Status(http.StatusUnauthorized)
		return
	}
	log = log.With("userId", userClaims.ID)

	rawBoingId := c.Param("id")
	boingId, err := strconv.Atoi(rawBoingId)
	if err != nil {
		log.Error("failed to parse boing id", "error", err)
		c.Status(http.StatusBadRequest)
		return
	}

	var request CreateCommentRequest
	err = c.BindJSON(&request)
	if err != nil {
		log.Debug("failed to bind request")
		c.JSON(500, response.NewErrorResponse("failed to bind request"))
		return
	}
	log = log.With("boingId", boingId, "text", request.Text)

	if boingId == 0 || request.Text == "" {
		log.Debug("invalid request", "boingId", boingId, "text", request.Text)
		//TODO: return validation errors like dinkdonk
		c.JSON(400, response.NewErrorResponse("invalid request"))
		return
	}

	err = h.s.Create(uint(boingId), userClaims.ID, request.Text)
	if err != nil {
		log.Error("failed to create comment")
		c.JSON(500, response.NewErrorResponse("failed to create comment"))
		return
	}

	c.Status(201)
}
