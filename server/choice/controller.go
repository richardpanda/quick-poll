package choice

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/vote"
	"github.com/richardpanda/quick-poll/server/ws"
	uuid "github.com/satori/go.uuid"
)

func IncrementNumVotes(wsConn *ws.Conn) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			db            = c.MustGet("db").(*gorm.DB)
			checkIP       = c.MustGet("checkIP").(bool)
			pollID        = c.Params.ByName("id")
			currentChoice = c.MustGet("currentChoice").(Choice)
			choiceID      = c.Params.ByName("choiceID")
			tx            = db.Begin()
		)

		if err := tx.Model(&currentChoice).Update("num_votes", currentChoice.NumVotes+1).Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"message": err})
			return
		}

		if checkIP {
			v := vote.Vote{
				ID:        uuid.NewV4().String(),
				PollID:    pollID,
				IPAddress: c.ClientIP(),
			}
			if err := tx.Create(&v).Error; err != nil {
				tx.Rollback()
				c.JSON(400, gin.H{"message": err})
				return
			}
		}

		tx.Commit()
		wsConn.BroadcastUpdate(pollID, currentChoice.ID, currentChoice.NumVotes)
		c.JSON(200, gin.H{"id": choiceID, "text": currentChoice.Text, "num_votes": currentChoice.NumVotes})
	}
}
