package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
)

type UtilImpl struct {
	Util
	validator *validator.Validate
}

func NewUtil() Util {
	return &UtilImpl{
		validator: validator.New(),
	}
}

func (v *UtilImpl) ValidateStruct(s interface{}) error {
	if v.validator == nil {
		panic("validator not initialized")
	}

	err := v.validator.Struct(s)
	if err != nil {
		return err
	}

	return nil
}

func (v *UtilImpl) FormatValidationError(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
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
