package model

import (
	"coin-api/domain/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UserResponse struct {
	UserId   uint   `json:"userid"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

type UserBalanceResponse struct {
	UserId  uint `json:"userid"`
	Balance int  `json:"balance"`
}

type UserAddForm struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func UserFromDomainModel(m *model.User) *UserResponse {
	u := &UserResponse{
		UserId:   m.ID,
		Name:     m.Username,
		Password: m.Password,
		Balance:  *m.CoinBalance,
	}

	return u
}

func UserBalanceFromDomainModel(m *model.User) *UserBalanceResponse {
	u := &UserBalanceResponse{
		UserId:  m.ID,
		Balance: *m.CoinBalance,
	}

	return u
}

func (u UserAddForm) ValidateUserAddForm() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserName, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
}
