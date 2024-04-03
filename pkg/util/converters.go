package util

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ConvertUnstructured(unstructured *unstructured.Unstructured, obj interface{}) error {
	b, err := json.Marshal(unstructured)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}

	return nil
}

func ConvertToUnstructured(obj interface{}) (*unstructured.Unstructured, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	unstructuredObj := &unstructured.Unstructured{}
	err = unstructuredObj.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}

	return unstructuredObj, nil
}
