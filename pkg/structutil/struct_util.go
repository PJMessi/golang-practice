package structutil

import "encoding/json"

func ConvertToBytes(structData interface{}) ([]byte, error) {
	responseInBytes, err := json.Marshal(structData)
	if err != nil {
		return nil, err
	}

	return responseInBytes, nil
}

func ConvertFromBytes(byteData []byte, targetStruct interface{}) error {
	if err := json.Unmarshal(byteData, targetStruct); err != nil {
		return err
	}
	return nil
}
