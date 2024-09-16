package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Tender struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ServType    string    `json:"serviceType"`
	OrgId       string    `json:"-"`
	Version     int64     `json:"version"`
	Created     time.Time `json:"createdAt"`
}

func (t *Tender) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&t.Description, validation.Required, validation.Length(1, 500)),
		validation.Field(&t.ServType, validation.Required, validation.In("Construction", "Delivery", "Manufacture")),
		// validation.Field(&t.Status, validation.Required, validation.In("Created", "Published")),
	)
}

func (t *Tender) ValidateEdition() error {
	// if t.Name == "" && t.Description == "" && t.ServType == "" {
	// 	return ErrFieldsEmpty
	// }
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Name, validation.Length(1, 100)),
		validation.Field(&t.Description, validation.Length(1, 500)),
		validation.Field(&t.ServType, validation.In("Construction", "Delivery", "Manufacture")),
	)
}
