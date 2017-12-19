package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	"github.com/richardpanda/quick-poll/server/pkg/vote"
)

func CheckDuplicateVote() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			db          = c.MustGet("db").(*gorm.DB)
			currentPoll = c.MustGet("currentPoll").(poll.Poll)
			pollID      = c.Params.ByName("id")
		)

		c.Set("checkIP", false)
		if !currentPoll.CheckIP {
			c.Next()
			return
		}

		c.Set("checkIP", true)
		err := db.Where("poll_id = ? AND ip_address = ?", pollID, c.ClientIP()).First(&vote.Vote{}).Error
		if err != nil {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(400, gin.H{"message": "You have already voted on this poll."})
	}
}
