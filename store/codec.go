package store

import "encoding/json"

func encode(value any) ([]byte, error) {
	return json.Marshal(value)
}

func decode(data []byte, target any) error {
	return json.Unmarshal(data, target)
}

func copyBytes(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}
