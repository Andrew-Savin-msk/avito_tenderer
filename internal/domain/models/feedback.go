package models

import "time"

type Feedback struct {
	Id      string    `json:"id"`
	Desc    string    `json:"description"`
	Created time.Time `json:"createdAt"`
}
