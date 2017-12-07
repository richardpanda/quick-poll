package poll

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine) {
	r.POST("/v1/polls", CreatePoll)
}
