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
	r.POST(
		"/v1/polls/:pollID/choices/:choiceID",
		middleware.ValidatePollID(),
		middleware.ValidateChoiceID(),
		middleware.CheckDuplicateVote(),
		choice.IncrementNumVotes(conn),
	)
	r.GET("/v1/polls/:id", poll.ReadOne)
	r.POST("/v1/polls", poll.Create)
	r.GET("/v1/ws", ws.OpenConnection(conn))
	return r
}

func NewTestRouter(db *gorm.DB, conn *ws.Conn) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return New(db, conn)
}
