package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints/response"
	"github.com/modhanami/boinger/endpoints/utils"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services"
	"net/http"
	"strconv"
)

func MakeListEndpoint(s services.BoingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		boings, err := s.List()
		if err != nil {
			c.JSON(500, response.ErrorResponseFromError(err))
			return
		}
		c.JSON(200, boings)
	}
}

type CreateRequest struct {
	Text string
}

func MakeCreateEndpoint(s services.BoingService, u services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(500, response.ErrorResponseFromError(err))
			return
		}

		userClaims := utils.GetUserClaimsFromContext(c)
		if userClaims == nil {
			c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Unauthorized"))
			return
		}

		user, err := u.GetById(userClaims.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponseFromError(err))
			return
		}

		err = s.Create(request.Text, user.ID)
		if err != nil {
			c.JSON(500, response.ErrorResponseFromError(err))
			return
		}
		c.Status(201)
	}
}

func MakeGetByIdEndpoint(s services.BoingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 0)
		if err != nil {
			c.JSON(500, response.ErrorResponseFromError(err))
			return
		}

		boing, err := s.GetById(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponseFromError(err))
			return
		} else if boing == (models.Boing{}) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse("Boing not found"))
			return
		}

		c.JSON(http.StatusOK, boing)
	}
}
