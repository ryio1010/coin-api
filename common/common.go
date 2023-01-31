package common

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func CreateJsonString(T any) string {
	s, err := json.Marshal(&T)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(s)
}

func StringToUint(v string) uint {
	// 対象文字列が数値であること前提(Validation後に使用)
	vInt, _ := strconv.ParseUint(v, 10, 64)
	return uint(vInt)
}
