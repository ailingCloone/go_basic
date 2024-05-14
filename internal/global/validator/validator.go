package validator

import (
	"github.com/go-playground/validator/v10"
)

// typeValidator checks if the value is one of the allowed types (EMAIL, CONTACT, IC)
func TypeValidator(fl validator.FieldLevel) bool {
	allowedTypes := map[string]bool{
		"EMAIL":   true,
		"CONTACT": true,
		"IC":      true,
	}
	typ := fl.Field().String()
	_, ok := allowedTypes[typ]
	return ok
}

func FromValidator(fl validator.FieldLevel) bool {
	allowedTypes := map[string]bool{
		"staff":    true,
		"customer": true,
	}
	typ := fl.Field().String()
	_, ok := allowedTypes[typ]
	return ok
}

func PageValidator(fl validator.FieldLevel) bool {
	allowedTypes := map[string]bool{
		"login":          true,
		"register":       true,
		"profile_update": true,
	}
	typ := fl.Field().String()
	_, ok := allowedTypes[typ]
	return ok
}

func RegisterListStatusValidator(fl validator.FieldLevel) bool {

	allowedTypes := map[string]bool{
		"1": true, // Pending
		"2": true, // Approve
		"3": true, // Reject
		"4": true, // All
	}

	typ := fl.Field().String()
	_, ok := allowedTypes[typ]
	return ok
}
