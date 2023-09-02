package exceptions

import "fmt"

type InvalidRequestException struct {
	Type    string
	Message string
	Details *string
}

func (e *InvalidRequestException) Error() string {
	message := fmt.Sprintf("%s: %s", e.Type, e.Message)
	if e.Details != nil {
		message += fmt.Sprintf("%s. %s", message, *e.Details)
	}
	return message
}

type NotFoundException struct {
	Type    string
	Message string
	Details *string
}

func (e *NotFoundException) Error() string {
	message := fmt.Sprintf("%s: %s", e.Type, e.Message)
	if e.Details != nil {
		message += fmt.Sprintf("%s. %s", message, *e.Details)
	}
	return message
}

type DuplicateException struct {
	Type    string
	Message string
	Details *string
}

func (e *DuplicateException) Error() string {
	message := fmt.Sprintf("%s: %s", e.Type, e.Message)
	if e.Details != nil {
		message += fmt.Sprintf("%s. %s", message, *e.Details)
	}
	return message
}
