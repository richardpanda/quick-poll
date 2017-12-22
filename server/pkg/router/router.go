package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/middleware"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	"github.com/richardpanda/quick-poll/server/pkg/ws"
)

func New(db *gorm.DB, conn *ws.Conn) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.SetDB(db))

	v1 := r.Group("/v1")
	v1.POST(
		"/polls/:id/choices/:choiceID",
		middleware.ValidatePollID(),
		middleware.ValidateChoiceID(),
		middleware.CheckDuplicateVote(),
		choice.IncrementNumVotes(conn),
	)
	v1.GET(
		"/polls/:id/ws",
		middleware.ValidatePollID(),
		ws.OpenConnection(conn),
	)
	v1.GET("/polls/:id", poll.ReadOne)
	v1.POST("/polls", poll.Create)

	return r
}

func NewTestRouter(db *gorm.DB, conn *ws.Conn) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return New(db, conn)
}
