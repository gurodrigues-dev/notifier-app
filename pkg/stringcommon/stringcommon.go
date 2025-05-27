package stringcommon

import (
	"encoding/json"
	"fmt"
)

func Empty(s string) bool {
	return s == ""
}

func SerializeToJSON(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize to JSON: %w", err)
	}
	return data, nil
}
