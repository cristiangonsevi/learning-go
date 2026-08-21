package handler

import "github.com/go-playground/validator/v10"

func transalateErrors(err error) map[string]string {
	errs := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		switch e.Tag() {
		case "required":
			errs[field] = field + " is required"
		case "min":
			errs[field] = field + " must be at least " + e.Param() + " characters"
		case "max":
			errs[field] = field + " must be maximum " + e.Param() + " characters"
		case "email":
			errs[field] = "Invalid email"
		}
	}
	return errs
}
