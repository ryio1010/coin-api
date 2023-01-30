package model

import (
	"coin-api/domain/model"
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type CoinAddUseForm struct {
	UserId    uint   `json:"userid"`
	Operation string `json:"operation"`
	Amount    int    `json:"amount"`
}

func (c CoinAddUseForm) ValidateCoinAddUseForm() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UserId, validation.Required),
		validation.Field(&c.Operation, validation.Required, validation.In(string(ADD), string(USE))),
		validation.Field(&c.Amount, validation.Required),
	)
}

type CoinSendForm struct {
	Sender   uint `json:"sender"`
	Receiver uint `json:"receiver"`
	Amount   int  `json:"amount"`
}

func (c CoinSendForm) ValidateCoinSendForm() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Sender, validation.Required),
		validation.Field(&c.Receiver, validation.Required),
		validation.Field(&c.Amount, validation.Required),
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

func CoinSendResponseFromDomainModel(c *CoinSendForm, balance int) *CoinSendResponse {
	h := &CoinSendResponse{
		Sender:        c.Sender,
		Receiver:      c.Receiver,
		Amount:        c.Amount,
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
