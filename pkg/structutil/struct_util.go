package structutil

import (
	"encoding/json"
	"fmt"
)

func ConvertToBytes(structData interface{}) ([]byte, error) {
	responseInBytes, err := json.Marshal(structData)
	if err != nil {
		return nil, fmt.Errorf("structutil.ConvertToBytes(): %w", err)
	}

	return responseInBytes, nil
}

func ConvertFromBytes(byteData []byte, targetStruct interface{}) error {
	return json.Unmarshal(byteData, targetStruct)
}
