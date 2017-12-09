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

	if c.Request.Body == nil {
		c.JSON(400, gin.H{"message": "Request body is missing."})
		return
	}

	var requestBody POSTPollsRequestBody
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(500, gin.H{"message": err})
		log.Fatal(err)
	}

	if requestBody.Question == "" {
		c.JSON(400, gin.H{"message": "Question is required."})
		return
	}

	if len(requestBody.Choices) <= 1 {
		c.JSON(400, gin.H{"message": "Please provide at least two choices."})
		return
	}

	choices := make([]choice.Choice, len(requestBody.Choices))
	for i, text := range requestBody.Choices {
		choices[i] = choice.New(text)
	}

	poll := Poll{ID: uuid.NewV4().String(), Question: requestBody.Question, Choices: choices}
	newDB := db.Create(&poll)
	if newDB.Error != nil {
		c.JSON(400, gin.H{"message": newDB.Error})
		return
	}

	c.JSON(200, gin.H{"id": poll.ID, "question": requestBody.Question, "choices": choices})
}

func orderChoicesCreatedAt(db *gorm.DB) *gorm.DB {
	return db.Order("choices.created_at ASC")
}

func ReadPoll(c *gin.Context) {
	var (
		db   = c.MustGet("db").(*gorm.DB)
		id   = c.Params.ByName("id")
		poll = Poll{ID: id}
	)

	err := db.Preload("Choices", orderChoicesCreatedAt).Find(&poll).Error
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid poll ID."})
		return
	}

	c.JSON(200, poll)
}
