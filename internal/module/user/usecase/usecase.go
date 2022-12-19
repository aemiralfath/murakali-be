package usecase

import (
	"context"
	"database/sql"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"
)

type userUC struct {
	cfg      *config.Config
	txRepo   *postgre.TxRepo
	userRepo user.Repository
}

func NewUserUseCase(cfg *config.Config, txRepo *postgre.TxRepo, userRepo user.Repository) user.UseCase {
	return &userUC{cfg: cfg, txRepo: txRepo, userRepo: userRepo}
}

func (u *userUC) EditUser(ctx context.Context, userID string, requestBody body.EditUserRequest) (*model.User, error) {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}
	}

	usernameUser, err := u.userRepo.GetUserByUsername(ctx, requestBody.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if usernameUser != nil {
		if *userModel.Username != *usernameUser.Username {
			return nil, httperror.New(http.StatusBadRequest, response.UserNameAlreadyExistMessage)
		}
	}

	phoneNoUser, err := u.userRepo.GetUserByPhoneNo(ctx, requestBody.PhoneNo)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if phoneNoUser != nil {
		if *userModel.PhoneNo != *phoneNoUser.PhoneNo {
			return nil, httperror.New(http.StatusBadRequest, response.PhoneNoAlreadyExistMessage)
		}
	}

	if !userModel.IsVerify {
		return nil, httperror.New(http.StatusBadRequest, response.UserNotVerifyMessage)
	}

	birthDate, _ := time.Parse("02-01-2006", requestBody.BirthDate)

	userModel.Username = &requestBody.Username
	userModel.FullName = &requestBody.FullName
	userModel.PhoneNo = &requestBody.PhoneNo
	userModel.BirthDate.Time = birthDate
	userModel.BirthDate.Valid = true
	userModel.Gender = &requestBody.Gender
	userModel.UpdatedAt.Time = time.Now()
	userModel.UpdatedAt.Valid = true

	err = u.userRepo.UpdateUserField(ctx, userModel)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}
