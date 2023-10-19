package validation

type Handler interface {
	ValidateStruct(s interface{}) error
}

type ValidationError struct {
	Details map[string]string
}

func (v ValidationError) Error() string {
	return "invalid data"
}
