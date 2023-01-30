package presenter

import (
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserPresenter struct {
	ctx *gin.Context
}

func NewUserOutputPort(context *gin.Context) ports.UserOutputPort {
	return &UserPresenter{
		ctx: context,
	}
}

func (u *UserPresenter) OutputUser(user *model.UserResponse) error {
	u.ctx.JSON(http.StatusOK, user)
	return nil
}

func (u *UserPresenter) OutputUserBalance(balance *model.UserBalanceResponse) error {
	u.ctx.JSON(http.StatusOK, balance)
	return nil
}

func (u *UserPresenter) OutputError(err *model.ErrorResponse) error {
	u.ctx.JSON(err.ErrorCode, err)
	return nil
}
