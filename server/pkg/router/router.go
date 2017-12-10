package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/middleware"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
)

func New(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.SetDB(db))
	r.POST("/v1/choices/:id", choice.IncrementNumVotes)
	r.GET("/v1/polls/:id", poll.ReadOne)
	r.POST("/v1/polls", poll.Create)
	return r
}

func NewTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return New(db)
}
