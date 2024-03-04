package util

import (
	"encoding/json"

	"github.com/go-openapi/spec"
)

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func MapToSpecSchema(obj interface{}) (*spec.Schema, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	result := &spec.Schema{}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
