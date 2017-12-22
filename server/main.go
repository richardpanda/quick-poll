package main

import (
	"fmt"
	"log"
	"os"

	"github.com/richardpanda/quick-poll/server/pkg/vote"
	"github.com/richardpanda/quick-poll/server/pkg/ws"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/richardpanda/quick-poll/server/pkg/choice"
	"github.com/richardpanda/quick-poll/server/pkg/poll"
	"github.com/richardpanda/quick-poll/server/pkg/router"
)

func main() {
	var (
		host             = os.Getenv("POSTGRES_HOST")
		user             = os.Getenv("POSTGRES_USER")
		password         = os.Getenv("POSTGRES_PASSWORD")
		dbName           = os.Getenv("POSTGRES_DB")
		connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)
		conn             = ws.NewConn()
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(&poll.Poll{}, &choice.Choice{}, &vote.Vote{})

	r := router.New(db, conn)
	r.Run(":8080")
}
