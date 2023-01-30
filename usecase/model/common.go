package model

import (
	"encoding/json"
	"fmt"
)

func CreateJsonString(T any) string {
	s, err := json.Marshal(&T)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(s)
}
