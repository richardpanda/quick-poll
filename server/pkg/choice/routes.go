package choice

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine) {
	r.POST("/v1/choices/:id", IncrementNumVotes)
}
