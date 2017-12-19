package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
)

func ValidatePollID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			db          = c.MustGet("db").(*gorm.DB)
			pollID      = c.Params.ByName("pollID")
			currentPoll = poll.Poll{ID: pollID}
		)

		if err := db.First(&currentPoll).Error; err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid poll ID."})
			return
		}

		c.Set("currentPoll", currentPoll)
		c.Next()
	}
}
