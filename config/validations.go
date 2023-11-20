package config

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// fn for custom message when validating data
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo es obligatorio"
	case "min":
		return "El número o cantidad de caracteres minimo(s) permitido(s) se ha sobrepasado"
	case "max":
		return "El número o cantidad de caracteres máximo(s) permitido(s) se ha sobrepasado"
	case "oneof":
		return "El valor del campo no existe dentro de las opciones"
	case "gt":
		return "El valor del campo debe ser más grande"
	}
	return fe.Error() // default error
}

func ValidateData(reqPayload interface{}) []CustomValidationError {
	validate := validator.New()
	err := validate.Struct(reqPayload)
	if err != nil {
		var customErrors []CustomValidationError

		for _, err := range err.(validator.ValidationErrors) {
			customError := CustomValidationError{
				Field:   err.Field(),
				Message: msgForTag(err),
			}
			customErrors = append(customErrors, customError)
		}

		return customErrors
	}
	return nil
}
