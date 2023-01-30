package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"column:username"`
	Password    string `gorm:"column:password"`
	CoinBalance *int   `gorm:"column:coinbalance"`
}

func CreateJsonStringFromUserModel(user *User) string {
	s, err := json.Marshal(&user)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(s)
}
