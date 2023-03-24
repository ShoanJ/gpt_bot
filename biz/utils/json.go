package utils

import "encoding/json"

func JsonMarshal(obj any) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
