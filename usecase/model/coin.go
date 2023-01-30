package model

import (
	"coin-api/common/enum"
	"coin-api/domain/model"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"time"
)

type CoinAddUseForm struct {
	UserId    string `json:"userid"`
	Operation string `json:"operation"`
	Amount    string `json:"amount"`
}

func (c CoinAddUseForm) ValidateCoinAddUseForm() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UserId, validation.Required, is.Digit),
		validation.Field(&c.Operation, validation.Required, validation.In(string(enum.ADD), string(enum.USE))),
		validation.Field(&c.Amount, validation.Required, is.Digit),
	)
}

type CoinSendForm struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (c CoinSendForm) ValidateCoinSendForm() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Sender, validation.Required, is.Digit),
		validation.Field(&c.Receiver, validation.Required, is.Digit),
		validation.Field(&c.Amount, validation.Required, is.Digit),
	)
}

type CoinResponse struct {
	UserId    uint   `json:"userid"`
	Operation string `json:"operation"`
	Amount    int    `json:"amount"`
	Balance   int    `json:"balance"`
}

type CoinHistoryResponse struct {
	Operation          string    `json:"operation"`
	OperationTimestamp time.Time `json:"operation_timestamp"`
	Amount             int       `json:"amount"`
}

type CoinSendResponse struct {
	Sender        uint `json:"sender"`
	Receiver      uint `json:"receiver"`
	Amount        int  `json:"amount"`
	SenderBalance int  `json:"sender_balance"`
}

func CoinResponseFromDomainModel(c *model.CoinHistory, balance int) *CoinResponse {
	h := &CoinResponse{
		UserId:    c.UserId,
		Operation: c.Operation,
		Amount:    c.Amount,
		Balance:   balance,
	}

	return h
}

func CoinSendResponseFromDomainModel(sender uint, receiver uint, amount int, balance int) *CoinSendResponse {
	h := &CoinSendResponse{
		Sender:        sender,
		Receiver:      receiver,
		Amount:        amount,
		SenderBalance: balance,
	}

	return h
}

func CoinHistoryResponseFromDomainModel(c *model.CoinHistory) *CoinHistoryResponse {
	h := &CoinHistoryResponse{
		Operation:          c.Operation,
		OperationTimestamp: c.OperationTimestamp,
		Amount:             c.Amount,
	}

	return h
}
