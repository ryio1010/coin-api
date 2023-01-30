package interactor

import (
	"coin-api/common"
	"coin-api/common/enum"
	models "coin-api/domain/model"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
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
		log.Log().Msg(fmt.Sprintf("バリデーションエラー CoinAddUseForm : %s", common.CreateJsonString(&form)))
		log.Error().Stack().Err(err)

		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()), err)
	}

	// coin残高処理対象ユーザーの取得
	uidUint := common.StringToUint(form.UserId)
	user, err := c.userRepo.SelectById(uidUint)
	if err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	// 残高不足確認
	amountInt, _ := strconv.Atoi(form.Amount)
	if form.Operation == string(enum.USE) {
		if *user.CoinBalance < amountInt {
			// 消費量が残高を上回る場合はエラー
			log.Log().Msg(fmt.Sprintf("コイン残高不足エラー user : %s", common.CreateJsonString(&user)))
			log.Error().Stack().Err(err)

			return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, "コイン残高不足エラー"), err)
		}
		// 区分がUSEの場合は符号を-に変換
		amountInt = -amountInt
	}

	// 残高の算出
	balance := *user.CoinBalance + amountInt
	user.CoinBalance = &balance

	// 残高の更新
	updated, err := c.userRepo.Update(user)
	if err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	// 履歴の追加
	operationTime := time.Now()
	target := models.CoinHistory{
		Operation:          form.Operation,
		OperationTimestamp: operationTime,
		UserId:             uidUint,
		Amount:             amountInt,
	}
	history, err := c.coinRepo.Insert(&target)
	if err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	return c.op.OutputCoin(model.CoinResponseFromDomainModel(history, *updated.CoinBalance))
}

func (c *CoinUseCase) SendCoin(form *model.CoinSendForm) error {
	// formのバリデーション
	err := form.ValidateCoinSendForm()
	if err != nil {
		log.Log().Msg(fmt.Sprintf("バリデーションエラー CoinSendForm : %s", common.CreateJsonString(&form)))
		log.Error().Stack().Err(err)

		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()), err)
	}

	// Senderの残高取得
	senderUidUint := common.StringToUint(form.Sender)
	sender, err := c.userRepo.SelectById(senderUidUint)
	if err != nil {
		log.Log().Msg(fmt.Sprintf("Senderユーザー取得に失敗 user : %s", common.CreateJsonString(&sender)))
		log.Error().Stack().Err(err)

		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	// Receiverの残高取得
	receiverUidUint := common.StringToUint(form.Receiver)
	receiver, err := c.userRepo.SelectById(receiverUidUint)
	if err != nil {
		log.Log().Msg(fmt.Sprintf("Receiverユーザー取得に失敗 user : %s", common.CreateJsonString(&receiver)))

		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	// Sender残高の確認
	amountInt, _ := strconv.Atoi(form.Amount)
	if *sender.CoinBalance < amountInt {
		// 消費量が残高を上回る場合はエラー
		log.Log().Msg(fmt.Sprintf("コイン残高不足エラー user : %s", common.CreateJsonString(&sender)))
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, "コイン残高不足エラー"), err)
	}

	// Sender残高の更新
	senderBalance := *sender.CoinBalance + (-amountInt)
	sender.CoinBalance = &senderBalance

	senderInfo, err := c.userRepo.Update(sender)
	if err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	// Receiver残高の更新
	receiverBalance := *receiver.CoinBalance + amountInt
	receiver.CoinBalance = &receiverBalance

	_, err = c.userRepo.Update(receiver)
	if err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	// Sender履歴
	operationTime := time.Now()
	senderInsertion := models.CoinHistory{
		Operation:          string(enum.SEND),
		OperationTimestamp: operationTime,
		UserId:             senderUidUint,
		Amount:             -amountInt,
	}
	// Receiver履歴
	receiverInsertion := models.CoinHistory{
		Operation:          string(enum.RECEIVE),
		OperationTimestamp: operationTime,
		UserId:             receiverUidUint,
		Amount:             amountInt,
	}

	// 履歴一括作成
	histories := make([]*models.CoinHistory, 0)
	histories = append(histories, &senderInsertion, &receiverInsertion)

	_, err = c.coinRepo.BatchInsert(histories)
	if err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	return c.op.OutputCoinSend(model.CoinSendResponseFromDomainModel(senderUidUint, receiverUidUint, amountInt, *senderInfo.CoinBalance))
}

func (c *CoinUseCase) SelectHistoriesByUserId(uid string) error {
	// uidのバリデーション
	err := validation.Validate(uid, validation.Required, is.Digit)
	if err != nil {
		log.Log().Msg(fmt.Sprintf("バリデーションエラー ユーザーID : %s", uid))
		log.Error().Stack().Err(err)

		return c.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()), err)
	}

	// idに紐づく履歴取得
	uidUint := common.StringToUint(uid)
	histories, err := c.coinRepo.SelectHistoriesByUserId(uidUint)
	if err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	// response用に詰め替え
	response := make([]*model.CoinHistoryResponse, 0)
	for _, v := range histories {
		entity := model.CoinHistoryResponseFromDomainModel(&v)
		response = append(response, entity)
	}

	return c.op.OutputCoinHistory(response)
}
