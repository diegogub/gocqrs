package lib

import (
	"encoding/json"
)

func Map(i interface{}) map[string]interface{} {
	d := make(map[string]interface{})
	b, _ := json.Marshal(i)
	json.Unmarshal(b, &d)
	return d
}
