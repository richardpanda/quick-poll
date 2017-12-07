package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

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
