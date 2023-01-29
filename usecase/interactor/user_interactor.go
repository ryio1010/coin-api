package interactor

import (
	models "coin-api/domain/model"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserUseCase struct {
	op       ports.UserOutputPort
	userRepo repository.IUserRepository
}

func NewUserUseCase(uop ports.UserOutputPort, ur repository.IUserRepository) ports.UserInputPort {
	return &UserUseCase{
		op:       uop,
		userRepo: ur,
	}
}

func (u *UserUseCase) RegisterUser(user *model.UserAddForm) error {
	// formのバリデーション
	err := user.ValidateUserAddForm()
	if err != nil {
		log.Fatal().Err(err)
		return u.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	// パスワードハッシュ化
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// Insert対象データ作成
	insertionTarget := models.User{
		Username: user.UserName,
		Password: string(passwordHash),
	}

	// ユーザー登録処理実行
	insertedUser, err := u.userRepo.Insert(insertionTarget)

	if err != nil {
		log.Fatal().Err(err)
		return u.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	return u.op.OutputUser(model.UserFromDomainModel(insertedUser))
}

func (u *UserUseCase) GetBalanceByUserId(uid uint) error {
	// uidのバリデーション
	err := validation.Validate(uid, validation.Required)
	if err != nil {
		log.Fatal().Err(err)
		return u.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	// ユーザー取得処理実行
	user, err := u.userRepo.SelectById(uid)

	if err != nil {
		log.Fatal().Err(err)
		return u.op.OutputError(model.CreateResponse(http.StatusBadRequest, err.Error()))
	}

	return u.op.OutputUserBalance(model.UserBalanceFromDomainModel(user))
}
