package relational

import (
	"time"
)

type Report struct {
	ID                uint `gorm:"primary_key"`
	Text              string
	Device            string
	Location          string
	CreatedAt         time.Time
	Provider          string `gorm:"unique_index:idx_provider_name_id"`
	ProviderID        uint   `gorm:"unique_index:idx_provider_name_id"`
	ProviderCreatedAt time.Time
}
