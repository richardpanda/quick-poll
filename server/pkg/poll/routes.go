package poll

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine) {
	r.GET("/v1/polls/:id", ReadPoll)
	r.POST("/v1/polls", CreatePoll)
}
