package validator

import (
	"devbubble-api/pkg/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New(validate *validator.Validate) *Validator {
	return &Validator{validate}
}

func (v *Validator) Validate(w http.ResponseWriter, s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Handle decoding error (e.g., invalid JSON format)
		json.HttpError(w, http.StatusBadRequest, validationErrors.Error())
		return validationErrors

	}

	return nil
}
