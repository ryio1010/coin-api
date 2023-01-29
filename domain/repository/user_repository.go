package repository

import (
	"coin-api/domain/model"
)

type IUserRepository interface {
	SelectById(id uint) (*model.User, error)
	Insert(user model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
}
