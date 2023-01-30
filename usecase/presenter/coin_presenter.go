package presenter

import (
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CoinPresenter struct {
	ctx *gin.Context
}

func NewCoinOutputPort(context *gin.Context) ports.CoinOutputPort {
	return &CoinPresenter{
		ctx: context,
	}
}

func (c *CoinPresenter) OutputCoinHistory(histories []*model.CoinHistoryResponse) error {
	c.ctx.JSON(http.StatusOK, histories)
	return nil
}

func (c *CoinPresenter) OutputCoin(coin *model.CoinResponse) error {
	c.ctx.JSON(http.StatusOK, coin)
	return nil
}

func (c *CoinPresenter) OutputCoinSend(coin *model.CoinSendResponse) error {
	c.ctx.JSON(http.StatusOK, coin)
	return nil
}

func (c *CoinPresenter) OutputError(res *model.ErrorResponse, err error) error {
	c.ctx.JSON(res.ErrorCode, res)
	return err
}
