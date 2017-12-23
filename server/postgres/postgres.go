package postgres

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

func Connect(host, user, password, dbName string) (*gorm.DB, error) {
	var connectionString string
	if password == "" {
		connectionString = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", host, user, dbName)
	} else {
		connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)
	}

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectTest(t *testing.T) *gorm.DB {
	var (
		host     = os.Getenv("POSTGRES_HOST")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbName   = os.Getenv("POSTGRES_TEST_DB")
	)

	db, err := Connect(host, user, password, dbName)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
