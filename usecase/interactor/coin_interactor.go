package interactor

import (
	models "coin-api/domain/model"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type CoinUseCase struct {
	op       ports.CoinOutputPort
	coinRepo repository.ICoinRepository
	userRepo repository.IUserRepository
}

func NewCoinUseCase(uop ports.CoinOutputPort, cr repository.ICoinRepository, ur repository.IUserRepository) ports.CoinInputPort {
	return &CoinUseCase{
		op:       uop,
		coinRepo: cr,
		userRepo: ur,
	}
}

func (c *CoinUseCase) AddUseCoin(form *model.CoinAddUseForm) error {
	// formのバリデーション
	err := form.ValidateCoinAddUseForm()
	if err != nil {
		log.Log().Msg(fmt.Sprintf("バリデーションエラー CoinAddUseForm : %s", model.CreateJsonString(&form)))
		log.Error().Stack().Err(err).Msg("")
		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// coin残高処理対象ユーザーの取得
	user, err := c.userRepo.SelectById(form.UserId)
	if err != nil {
		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// 処理区分取得
	operation := form.Operation
	if operation == string(model.USE) {
		if *user.CoinBalance < form.Amount {
			// 消費量が残高を上回る場合はエラー
			log.Log().Msg(fmt.Sprintf("コイン残高不足エラー user : %s", model.CreateJsonString(&user)))
			return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, "コイン残高不足エラー"))
		}
		// 区分がUSEの場合は符号を-に変換
		form.Amount = -form.Amount
	}

	// 残高の算出
	balance := *user.CoinBalance + form.Amount
	// 残高の更新
	user.CoinBalance = &balance
	updated, err := c.userRepo.Update(user)

	// 履歴の追加
	operationTime := time.Now()
	target := models.CoinHistory{
		Operation:          form.Operation,
		OperationTimestamp: operationTime,
		UserId:             form.UserId,
		Amount:             form.Amount,
	}
	history, err := c.coinRepo.Insert(&target)
	if err != nil {
		log.Error().Err(err).Send()
		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.op.OutputCoin(model.CoinResponseFromDomainModel(history, *updated.CoinBalance))
}

func (c *CoinUseCase) SendCoin(form *model.CoinSendForm) error {
	// formのバリデーション
	err := form.ValidateCoinSendForm()
	if err != nil {
		log.Log().Msg(fmt.Sprintf("バリデーションエラー CoinSendForm : %s", model.CreateJsonString(&form)))
		log.Error().Stack().Err(err).Msg("")
		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Sender/Receiverの残高確認
	sender, err := c.userRepo.SelectById(form.Sender)
	if err != nil {
		log.Log().Msg(fmt.Sprintf("Senderユーザー取得に失敗 user : %s", model.CreateJsonString(&sender)))
		log.Error().Err(err).Send()
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	receiver, err := c.userRepo.SelectById(form.Receiver)
	if err != nil {
		log.Log().Msg(fmt.Sprintf("Receiverユーザー取得に失敗 user : %s", model.CreateJsonString(&receiver)))
		log.Error().Err(err).Send()
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Sender残高の確認
	if *sender.CoinBalance < form.Amount {
		// 消費量が残高を上回る場合はエラー
		log.Log().Msg(fmt.Sprintf("コイン残高不足エラー user : %s", model.CreateJsonString(&sender)))
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, "コイン残高不足エラー"))
	}

	// 残高の算出
	senderBalance := *sender.CoinBalance + (-form.Amount)
	receiverBalance := *receiver.CoinBalance + form.Amount

	// 残高の更新
	sender.CoinBalance = &senderBalance
	receiver.CoinBalance = &receiverBalance

	senderInfo, err := c.userRepo.Update(sender)
	_, err = c.userRepo.Update(receiver)
	if err != nil {
		log.Error().Err(err).Send()
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Sender履歴の追加
	operationTime := time.Now()
	senderInsertion := models.CoinHistory{
		Operation:          string(model.SEND),
		OperationTimestamp: operationTime,
		UserId:             form.Sender,
		Amount:             -(form.Amount),
	}

	// Receiver履歴の追加
	receiverInsertion := models.CoinHistory{
		Operation:          string(model.RECEIVE),
		OperationTimestamp: operationTime,
		UserId:             form.Receiver,
		Amount:             form.Amount,
	}
	_, err = c.coinRepo.Insert(&senderInsertion)
	_, err = c.coinRepo.Insert(&receiverInsertion)
	if err != nil {
		log.Error().Err(err).Send()
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.op.OutputCoinSend(model.CoinSendResponseFromDomainModel(form, *senderInfo.CoinBalance))
}

func (c *CoinUseCase) SelectHistoryByUserId(uid uint) error {
	// uidのバリデーション
	err := validation.Validate(uid, validation.Required)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// idに紐づく履歴取得
	histories, err := c.coinRepo.SelectHistoriesByUserId(uid)
	if err != nil {
		log.Error().Err(err).Send()
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// response用に詰め替え
	response := make([]*model.CoinHistoryResponse, 0)
	for _, v := range histories {
		entity := model.CoinHistoryResponseFromDomainModel(&v)
		response = append(response, entity)
	}

	return c.op.OutputCoinHistory(response)
}
