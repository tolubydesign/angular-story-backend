package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}

// NOTE: The Laws of Reflection - https://go.dev/blog/laws-of-reflection
func ValidateLimitedStringVariable(str string) error {
	value := reflect.TypeOf(str)
	if value.Kind() != reflect.String {
		return errors.New("Variable is not a valid string.")
	}

	if strings.ReplaceAll(str, " ", "") == "" {
		return errors.New("A valid string must be provided.")
	}

	strLength := len(str)
	if strLength > 255 {
		return errors.New("String is too long.")
	}

	additionalChecks := AdditionalStringValidation(str)
	if additionalChecks != nil {
		return additionalChecks
	}

	return nil
}

func AdditionalStringValidation(str string) error {
	value := reflect.TypeOf(str)
	if value.Kind() != reflect.String {
		return errors.New("Additional validation found. Invalid value provided.")
	}

	return nil
}
