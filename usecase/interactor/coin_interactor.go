package interactor

import (
	models "coin-api/domain/model"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"errors"
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
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	// coin残高処理対象ユーザーの取得
	user, err := c.userRepo.SelectById(form.UserId)

	if err != nil {
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	// 処理区分取得
	operation := form.Operation

	if operation == string(model.USE) {
		if *user.CoinBalance < form.Amount {
			log.Error().Msg("コイン残高不足エラー")
			return c.op.OutputError(model.CreateResponse(http.StatusInternalServerError, "コイン残高不足エラー"))
		}
		// Useの場合は符号を-に変換
		form.Amount = -form.Amount
	}

	// 残高の算出
	balance := *user.CoinBalance + form.Amount

	// 残高の更新と履歴の追加は１つのトランザクションで実施
	// 残高の更新
	user.CoinBalance = &balance
	fmt.Println(user.CoinBalance)
	updated, err := c.userRepo.Update(user)

	// 履歴の追加
	operationTime := time.Now()
	insertionTarget := models.CoinHistory{
		Operation:          form.Operation,
		OperationTimestamp: operationTime,
		UserId:             form.UserId,
		Amount:             form.Amount,
	}
	history, err := c.coinRepo.Insert(&insertionTarget)

	if err != nil {
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	return c.op.OutputCoin(model.CoinResponseFromDomainModel(history, *updated.CoinBalance))
}

func (c *CoinUseCase) SendCoin(form *model.CoinSendForm) error {
	// formのバリデーション
	err := form.ValidateCoinSendForm()
	if err != nil {
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	// Sender/Receiverの残高確認
	sender, err := c.userRepo.SelectById(form.Sender)
	receiver, err := c.userRepo.SelectById(form.Receiver)

	if err != nil {
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	// Sender残高の算出
	senderBalance := *sender.CoinBalance + (-form.Amount)
	if senderBalance <= 0 {
		log.Fatal().Err(errors.New("コイン残高不足エラー"))
		return c.op.OutputError(model.CreateResponse(http.StatusInternalServerError, "コイン残高不足エラー"))
	}

	// Receiver残高の算出
	receiverBalance := *receiver.CoinBalance + form.Amount

	// 残高の更新
	sender.CoinBalance = &senderBalance
	receiver.CoinBalance = &receiverBalance
	senderInfo, err := c.userRepo.Update(sender)
	_, err = c.userRepo.Update(receiver)

	// Sender履歴の追加
	operationTime := time.Now()
	senderInsertion := models.CoinHistory{
		Operation:          string(model.SEND),
		OperationTimestamp: operationTime,
		UserId:             form.Sender,
		Amount:             form.Amount,
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
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	return c.op.OutputCoin(&model.CoinResponse{UserId: form.Sender, Operation: "Add", Amount: form.Amount, Balance: *senderInfo.CoinBalance})
}

func (c *CoinUseCase) SelectHistoryByUserId(uid uint) error {
	// uidのバリデーション
	err := validation.Validate(uid, validation.Required)
	if err != nil {
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	histories, err := c.coinRepo.SelectHistoriesByUserId(uid)
	if err != nil {
		log.Fatal().Err(err)
		return c.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}
	response := make([]*model.CoinHistoryResponse, 0)
	for _, v := range histories {
		entity := model.CoinHistoryResponseFromDomainModel(&v)
		response = append(response, entity)
	}

	return c.op.OutputCoinHistory(response)
}
