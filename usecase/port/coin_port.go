package ports

import (
	"coin-api/usecase/model"
)

type CoinInputPort interface {
	SelectHistoryByUserId(uid uint) error
	AddUseCoin(form *model.CoinAddUseForm) error
	SendCoin(form *model.CoinSendForm) error
}

type CoinOutputPort interface {
	OutputCoin(coin *model.CoinResponse) error
	OutputCoinSend(coin *model.CoinSendResponse) error
	OutputCoinHistory(histories []*model.CoinHistoryResponse) error
	OutputError(err *model.ErrorResponse) error
}
