package interactor

import (
	"coin-api/common"
	"coin-api/common/enum"
	models "coin-api/domain/model"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"context"
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
	tranRepo repository.ITxRepository
}

func NewCoinUseCase(uop ports.CoinOutputPort, cr repository.ICoinRepository, ur repository.IUserRepository, tr repository.ITxRepository) ports.CoinInputPort {
	return &CoinUseCase{
		op:       uop,
		coinRepo: cr,
		userRepo: ur,
		tranRepo: tr,
	}
}

func (c *CoinUseCase) AddUseCoin(ctx context.Context, form *model.CoinAddUseForm) error {
	// formのバリデーション
	if err := form.ValidateCoinAddUseForm(); err != nil {
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

	// 履歴オブジェクト生成
	operationTime := time.Now()
	target := &models.CoinHistory{
		Operation:          form.Operation,
		OperationTimestamp: operationTime,
		UserId:             uidUint,
		Amount:             amountInt,
	}

	// 同一transaction内で残高の更新と履歴の追加を実行
	if _, err := c.tranRepo.DoInTx(ctx, c.AddUseCoinAndUpdateBalance(user, target)); err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	return c.op.OutputCoin(model.CoinResponseFromDomainModel(target, balance))
}

func (c *CoinUseCase) AddUseCoinAndUpdateBalance(user *models.User, history *models.CoinHistory) func(ctx context.Context) (interface{}, error) {
	return func(ctx context.Context) (interface{}, error) {
		// 残高更新
		if _, err := c.userRepo.Update(ctx, user); err != nil {
			log.Error().Err(err).Send()
			return nil, err
		}

		// 履歴追加
		if _, err := c.coinRepo.Insert(ctx, history); err != nil {
			log.Error().Err(err).Send()
			return nil, err
		}
		return nil, nil
	}
}

func (c *CoinUseCase) SendCoin(ctx context.Context, form *model.CoinSendForm) error {
	// formのバリデーション
	if err := form.ValidateCoinSendForm(); err != nil {
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

	// Sender残高の設定
	senderBalance := *sender.CoinBalance + (-amountInt)
	sender.CoinBalance = &senderBalance

	// Receiver残高の設定
	receiverBalance := *receiver.CoinBalance + amountInt
	receiver.CoinBalance = &receiverBalance

	// Sender履歴作成
	operationTime := time.Now()
	senderInsertion := models.CoinHistory{
		Operation:          string(enum.SEND),
		OperationTimestamp: operationTime,
		UserId:             senderUidUint,
		Amount:             -amountInt,
	}
	// Receiver履歴作成
	receiverInsertion := models.CoinHistory{
		Operation:          string(enum.RECEIVE),
		OperationTimestamp: operationTime,
		UserId:             receiverUidUint,
		Amount:             amountInt,
	}

	// 履歴まとめ
	histories := make([]*models.CoinHistory, 0)
	histories = append(histories, &senderInsertion, &receiverInsertion)

	// 同一transaction内で残高の更新と履歴の追加を実行
	if _, err := c.tranRepo.DoInTx(ctx, c.SendCoinAndUpdateBalances(sender, receiver, histories)); err != nil {
		log.Error().Stack().Err(err)
		return c.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	return c.op.OutputCoinSend(model.CoinSendResponseFromDomainModel(senderUidUint, receiverUidUint, amountInt, senderBalance))
}

func (c *CoinUseCase) SendCoinAndUpdateBalances(sender *models.User, receiver *models.User, histories []*models.CoinHistory) func(ctx context.Context) (interface{}, error) {
	return func(ctx context.Context) (interface{}, error) {
		// Sender残高更新
		if _, err := c.userRepo.Update(ctx, sender); err != nil {
			log.Error().Err(err).Send()
			return nil, err
		}

		// Receiver残高更新
		if _, err := c.userRepo.Update(ctx, receiver); err != nil {
			log.Error().Err(err).Send()
			return nil, err
		}

		// 履歴一括追加
		if _, err := c.coinRepo.BatchInsert(ctx, histories); err != nil {
			log.Error().Err(err).Send()
			return nil, err
		}
		return nil, nil
	}
}

func (c *CoinUseCase) SelectHistoriesByUserId(uid string) error {
	// uidのバリデーション
	if err := validation.Validate(uid, validation.Required, is.Digit); err != nil {
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
