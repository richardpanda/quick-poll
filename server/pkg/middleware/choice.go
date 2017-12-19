package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
)

func ValidateChoiceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			db            = c.MustGet("db").(*gorm.DB)
			pollID        = c.Params.ByName("id")
			choiceID      = c.Params.ByName("choiceID")
			currentChoice = choice.Choice{ID: choiceID, PollID: pollID}
		)

		if err := db.First(&currentChoice).Error; err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid choice ID."})
			return
		}

		c.Set("currentChoice", currentChoice)
		c.Next()
	}
}
