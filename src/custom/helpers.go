package custom

import "encoding/json"

func RenameJSONKey(body *[]byte, oldKey string, newKey string) ([]byte, error) {
	var data map[string]interface{}

	err := json.Unmarshal(*body, &data)
	if err != nil {
		return nil, err
	}

	if value, exists := data[oldKey]; exists {
		data[newKey] = value
		delete(data, oldKey)
	}

	return json.Marshal(data)
}
