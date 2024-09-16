package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Bid struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"-"`
	Status      string    `json:"status"`
	TenderId    string    `json:"-"`
	AuthorType  string    `json:"authorType"`
	AuthorId    string    `json:"authorId"`
	Version     int64     `json:"version"`
	Created     time.Time `json:"createdAt"`
}

func (b *Bid) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&b.Description, validation.Required, validation.Length(1, 1000)),
		validation.Field(&b.TenderId, validation.Required, validation.Length(36, 100)),
		validation.Field(&b.AuthorType, validation.Required, validation.In("Organization", "User")),
		validation.Field(&b.AuthorId, validation.Required, validation.Length(36, 100)),
	)
}
