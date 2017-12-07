package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	"github.com/richardpanda/quick-poll/server/pkg/router"
)

func main() {
	var (
		username         = os.Getenv("DB_USERNAME")
		name             = os.Getenv("DB_NAME")
		connectionString = fmt.Sprintf("user=%s dbname=%s sslmode=disable", username, name)
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(&poll.Poll{}, &choice.Choice{})

	r := router.New(db)
	r.Run(":8080")
}
