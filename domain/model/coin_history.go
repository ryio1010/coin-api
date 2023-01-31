package model

import (
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
