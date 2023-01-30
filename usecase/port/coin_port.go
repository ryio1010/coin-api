package ports

import (
	"coin-api/usecase/model"
)

type CoinInputPort interface {
	SelectHistoriesByUserId(uid string) error
	AddUseCoin(form *model.CoinAddUseForm) error
	SendCoin(form *model.CoinSendForm) error
}

type CoinOutputPort interface {
	OutputCoin(coin *model.CoinResponse) error
	OutputCoinSend(coin *model.CoinSendResponse) error
	OutputCoinHistory(histories []*model.CoinHistoryResponse) error
	OutputError(res *model.ErrorResponse, err error) error
}
