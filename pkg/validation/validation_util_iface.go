package validation

type Util interface {
	ValidateStruct(s interface{}) error
	FormatValidationError(err error) string
}
