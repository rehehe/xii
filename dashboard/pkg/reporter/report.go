package reporter

import (
	"time"
)

type Report struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text,omitempty"`
	Device    string    `json:"device,omitempty"`
	Location  string    `json:"location,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	// uncomment if you need this omitted fields
	//Provider          string    `json:"provider,omitempty"`
	//ProviderID        uint      `json:"provider_id,omitempty"`
	//ProviderCreatedAt time.Time `json:"provider_created_at,omitempty"`
}
