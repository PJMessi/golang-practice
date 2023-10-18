package uuidutil

import (
	"fmt"

	"github.com/google/uuid"
)

func GenUuidV4() (string, error) {
	uuidObj, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("uuid.GenUuidV4(): %w", err)
	}
	return uuidObj.String(), nil
}
