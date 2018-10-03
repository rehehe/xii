package relational

import (
	"time"
)

type Report struct {
	ID        uint `gorm:"primary_key,auto_increment"`
	Text      string
	Device    string
	Location  string
	CreatedAt time.Time
}
