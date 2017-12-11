package choice

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/ws"
)

func IncrementNumVotes(wsConn *ws.Conn) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			db     = c.MustGet("db").(*gorm.DB)
			id     = c.Params.ByName("id")
			tx     = db.Begin()
			choice Choice
		)

		if err := tx.Where(&Choice{ID: id}).First(&choice).Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"message": "Invalid choice ID."})
			return
		}

		if err := tx.Model(&choice).Update("num_votes", choice.NumVotes+1).Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"message": err})
			return
		}

		tx.Commit()

		wsConn.BroadcastUpdate(choice.PollID, choice.ID, choice.NumVotes)
		c.JSON(200, gin.H{"id": id, "text": choice.Text, "num_votes": choice.NumVotes})
	}
}
