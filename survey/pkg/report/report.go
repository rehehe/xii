package report

import (
	"time"
)

type Report struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	Device    string    `json:"device,omitempty"`
	Location  string    `json:"location,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
