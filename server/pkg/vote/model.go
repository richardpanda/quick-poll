package vote

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Vote struct {
	ID        string     `json:"id" sql:"type:uuid;primary_key"`
	PollID    string     `json:"-" gorm:"unique_index:idx_poll_id_ip_address;type:uuid REFERENCES polls(id)"`
	IPAddress string     `json:"-" gorm:"unique_index:idx_poll_id_ip_address;type:inet"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func CreateTable(db *gorm.DB) {
	db.CreateTable(&Vote{})
}

func DropTable(db *gorm.DB) {
	db.DropTableIfExists("votes")
}
