package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/richardpanda/quick-poll/server/choice"
	"github.com/richardpanda/quick-poll/server/poll"
	"github.com/richardpanda/quick-poll/server/postgres"
	"github.com/richardpanda/quick-poll/server/router"
	"github.com/richardpanda/quick-poll/server/vote"
	"github.com/richardpanda/quick-poll/server/ws"
)

func main() {
	var (
		host     = os.Getenv("POSTGRES_HOST")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbName   = os.Getenv("POSTGRES_DB")
		conn     = ws.NewConn()
	)

	db, err := postgres.Connect(host, user, password, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&poll.Poll{}, &choice.Choice{}, &vote.Vote{})

	r := router.New(db, conn)
	r.Run(":8080")
}
