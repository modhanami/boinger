package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/services"
)

func MakeTimelineEndpoint(s services.TimelineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		timeline, err := s.List()
		if err != nil {
			c.JSON(500, ErrorResponseFromError(err))
			return
		}
		c.JSON(200, timeline)
	}
}
