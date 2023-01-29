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

type CoinSendForm struct {
	Sender   uint `json:"sender"`
	Receiver uint `json:"receiver"`
	Amount   int  `json:"amount"`
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

func CoinResponseFromDomainModel(c *model.CoinHistory, balance int) *CoinResponse {
	h := &CoinResponse{
		UserId:    c.UserId,
		Operation: c.Operation,
		Amount:    c.Amount,
		Balance:   balance,
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

func (c CoinAddUseForm) ValidateCoinAddUseForm() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UserId, validation.Required),
		validation.Field(&c.Operation, validation.Required),
		validation.Field(&c.Amount, validation.Required),
	)
}

func (c CoinSendForm) ValidateCoinSendForm() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Sender, validation.Required),
		validation.Field(&c.Receiver, validation.Required),
		validation.Field(&c.Amount, validation.Required),
	)
}
