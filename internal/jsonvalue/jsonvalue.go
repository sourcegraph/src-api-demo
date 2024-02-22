package jsonvalue

import (
	"encoding/json"
)

type JSONValue struct{ Value any }

func MarshalJSONValue(v *JSONValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func UnmarshalJSONValue(b []byte, v *JSONValue) error {
	return json.Unmarshal(b, &v.Value)
}

func (v JSONValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

func (v *JSONValue) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.Value)
}
