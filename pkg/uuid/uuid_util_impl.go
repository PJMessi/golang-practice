package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

type UtilImpl struct {
}

func NewUtil() Util {
	return &UtilImpl{}
}

func (u *UtilImpl) GenUuidV4() (string, error) {
	uuidObj, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("GenUuidV4: %w", err)
	}
	return uuidObj.String(), nil
}
