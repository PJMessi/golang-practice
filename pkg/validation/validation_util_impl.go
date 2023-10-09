package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
)

type UtilImpl struct {
	Util
	validatorIns *validator.Validate
}

func NewUtil() (Util, error) {
	validatorIns := validator.New()

	if validatorIns == nil {
		return nil, fmt.Errorf("validator.NewUtil: validatorIns is not initialized")
	}

	return &UtilImpl{
		validatorIns: validatorIns,
	}, nil
}

// ValidateStruct validates the provided struct based on struct tags besides the struct field names
func (v *UtilImpl) ValidateStruct(s interface{}) error {
	err := v.validatorIns.Struct(s)

	var validationErrs validator.ValidationErrors
	var invalidValidationErr *validator.InvalidValidationError

	if errors.As(err, &validationErrs) {
		return v.prepareValErrDetails(validationErrs)

	} else if errors.As(err, &invalidValidationErr) {
		return &ValidationError{}
	}

	return err
}

func (v *UtilImpl) prepareValErrDetails(valErrs validator.ValidationErrors) ValidationError {
	details := map[string]string{}

	for _, valErr := range valErrs {
		field := valErr.StructField()
		tag := valErr.Tag()
		details[field] = fmt.Sprintf("validation failed for tag: '%s'", tag)
	}

	return ValidationError{
		Details: details,
	}
}
