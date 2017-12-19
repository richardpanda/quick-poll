package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	"github.com/richardpanda/quick-poll/server/pkg/vote"
)

func CreateChoicesTable(db *gorm.DB) {
	db.CreateTable(&choice.Choice{})
}

func CreatePollsTable(db *gorm.DB) {
	db.CreateTable(&poll.Poll{})
}

func CreateVotesTable(db *gorm.DB) {
	db.CreateTable(&vote.Vote{})
}

func DBConnection(t *testing.T) (*gorm.DB, func()) {
	var (
		username         = os.Getenv("DB_USERNAME")
		name             = os.Getenv("TEST_DB_NAME")
		connectionString = fmt.Sprintf("user=%s dbname=%s sslmode=disable", username, name)
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		t.Fatal(err)
	}

	return db, func() { db.Close() }
}

func DropChoicesTable(db *gorm.DB) {
	db.DropTableIfExists("choices")
}

func DropPollsTable(db *gorm.DB) {
	db.DropTableIfExists("polls")
}

func DropVotesTable(db *gorm.DB) {
	db.DropTableIfExists("votes")
}
