package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/middleware"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
)

func New(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.SetDB(db))

	poll.AddRoutes(r)

	return r
}

func newTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return New(db)
}
