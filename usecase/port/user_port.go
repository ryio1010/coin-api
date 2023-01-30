package ports

import (
	"coin-api/usecase/model"
)

type UserInputPort interface {
	RegisterUser(user *model.UserAddForm) error
	GetBalanceByUserId(uid string) error
}

type UserOutputPort interface {
	OutputUser(user *model.UserResponse) error
	OutputUserBalance(balance *model.UserBalanceResponse) error
	OutputError(res *model.ErrorResponse, err error) error
}
