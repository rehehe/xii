package aggregator

import (
	"time"
)

type Report struct {
	ID                uint      `json:"id"`
	Text              string    `json:"text,omitempty"`
	Device            string    `json:"device,omitempty"`
	Location          string    `json:"location,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	Provider          string    `json:"provider"`
	ProviderID        uint      `json:"provider_id"`
	ProviderCreatedAt time.Time `json:"provider_created_at"`
}
