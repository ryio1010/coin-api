package interactor

import (
	"coin-api/common"
	models "coin-api/domain/model"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserUseCase struct {
	op ports.UserOutputPort
	ur repository.IUserRepository
}

func NewUserUseCase(uop ports.UserOutputPort, ur repository.IUserRepository) ports.UserInputPort {
	return &UserUseCase{
		op: uop,
		ur: ur,
	}
}

func (u *UserUseCase) RegisterUser(form *model.UserAddForm) error {
	// formのバリデーション
	if err := form.ValidateUserAddForm(); err != nil {
		log.Log().Msg(fmt.Sprintf("バリデーションエラー UserAddForm : %s", common.CreateJsonString(&form)))
		log.Error().Stack().Err(err)

		return u.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()), err)
	}

	// パスワードハッシュ化
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	// Insert対象データ作成
	target := models.User{
		Username: form.UserName,
		Password: string(pwHash),
	}

	// ユーザー登録処理実行
	user, err := u.ur.Insert(&target)
	if err != nil {
		log.Error().Stack().Err(err)
		return u.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	return u.op.OutputUser(model.UserFromDomainModel(user))
}

func (u *UserUseCase) GetBalanceByUserId(uid string) error {
	// uidのバリデーション
	if err := validation.Validate(uid, validation.Required, is.Digit); err != nil {
		log.Error().Msg(fmt.Sprintf("バリデーションエラー ユーザーID : %s", uid))
		log.Error().Stack().Err(err)

		return u.op.OutputError(model.CreateErrorResponse(http.StatusBadRequest, err.Error()), err)
	}

	// ユーザー取得処理実行
	uidUint := common.StringToUint(uid)
	user, err := u.ur.SelectById(uidUint)
	if err != nil {
		log.Error().Stack().Err(err)
		return u.op.OutputError(model.CreateErrorResponse(http.StatusInternalServerError, err.Error()), err)
	}

	return u.op.OutputUserBalance(model.UserBalanceFromDomainModel(user))
}
