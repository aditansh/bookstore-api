package utils

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {

	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidatePriceRange(priceRange []float64) bool {

	// Check if priceRange has exactly two elements and the first is less than or equal to the second
	if len(priceRange) != 2 || priceRange[0] > priceRange[1] {
		return false
	}

	return true
}
