package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"column:username"`
	Password    string `gorm:"column:password"`
	CoinBalance *int   `gorm:"column:coinbalance"`
}
