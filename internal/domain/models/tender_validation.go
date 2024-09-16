package models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func atLeastThreeFieldsNotEmpty(fields ...string) validation.RuleFunc {
	return func(value interface{}) error {
		count := 0
		for _, field := range fields {
			if field != "" {
				count++
			}
		}
		if count >= 1 {
			return nil
		}
		return fmt.Errorf("at least 1 fields must be non-empty")
	}
}
