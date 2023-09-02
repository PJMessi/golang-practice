package validation

import (
	"fmt"

	"github.com/go-playground/validator"
)

type Validator struct {
	validator *validator.Validate
}

func CreateValidator() *Validator {
	return &Validator{}
}

func (v *Validator) InitializeValidator() {
	v.validator = validator.New()
}

func (v *Validator) ValidateStruct(s interface{}) error {
	if v.validator == nil {
		panic("validator not initialized")
	}

	err := v.validator.Struct(s)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) FormatValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errorMsg := ""
		for _, vErr := range errs {
			field := vErr.StructField()
			tag := vErr.Tag()
			errorMsg += fmt.Sprintf("'%s' validation failed for tag '%s'. ", field, tag)
		}
		return errorMsg
	}

	return err.Error()
}
