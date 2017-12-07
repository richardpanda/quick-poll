package poll

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	uuid "github.com/satori/go.uuid"
)

func CreatePoll(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var requestBody POSTPollsRequestBody
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Fatal(err)
	}

	choices := make([]choice.Choice, len(requestBody.Choices))
	for i, text := range requestBody.Choices {
		choices[i] = choice.Choice{ID: uuid.NewV4().String(), Text: text}
	}

	poll := Poll{ID: uuid.NewV4().String(), Question: requestBody.Question, Choices: choices}
	newDB := db.Create(&poll)
	if newDB.Error != nil {
		c.JSON(400, gin.H{"message": newDB.Error})
		return
	}

	c.JSON(200, gin.H{"id": poll.ID, "question": requestBody.Question, "choices": choices})
}
