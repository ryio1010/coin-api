package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type CoinHistory struct {
	gorm.Model
	Operation          string    `gorm:"column:operation"`
	OperationTimestamp time.Time `gorm:"column:operation_timestamp"`
	UserId             uint      `gorm:"column:userid"`
	Amount             int       `gorm:"column:amount"`
}

func CreateJsonStringFromHistoryModel(history *CoinHistory) string {
	s, err := json.Marshal(&history)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(s)
}
