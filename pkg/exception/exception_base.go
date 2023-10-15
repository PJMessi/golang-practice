package exception

import (
	"fmt"

	"github.com/pjmessi/golang-practice/pkg/strutil"
)

type Base struct {
	Type    string             `json:"type"`
	Message string             `json:"message"`
	Details *map[string]string `json:"details"`
}

func (e *Base) Error() string {
	message := fmt.Sprintf("%s: %s", e.Type, e.Message)
	if e.Details != nil {
		message += fmt.Sprintf(". %s", *e.Details)
	}
	return message
}

func newException(ex *Base, defaultType, defaultMessage string) *Base {
	if ex == nil {
		b := Base{}
		ex = &b
	}

	if strutil.IsEmpty(ex.Message) {
		ex.Message = defaultMessage
	}

	if strutil.IsEmpty(ex.Type) {
		ex.Type = defaultType
	}

	return ex
}
