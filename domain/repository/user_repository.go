package repository

import (
	"coin-api/domain/model"
	"context"
)

type IUserRepository interface {
	SelectById(id uint) (*model.User, error)
	Insert(user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
}
