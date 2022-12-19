package usecase

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/util"
	smtp "murakali/pkg/email"
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

func (u *userUC) EditEmail(ctx context.Context, userID string, requestBody body.EditEmailRequest) (*model.User, error) {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}
	}

	emailHistory, err := u.userRepo.CheckEmailHistory(ctx, requestBody.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if emailHistory != nil {
		if emailHistory.Email != userModel.Email {
			return nil, httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage)
		}

		if emailHistory.Email == userModel.Email {
			return nil, httperror.New(http.StatusBadRequest, response.EmailSamePreviousEmailMessage)
		}
	}

	if err := u.SendLinkOTPEmail(ctx, requestBody.Email); err != nil {
		return nil, err
	}
	return userModel, nil
}

func (u *userUC) EditEmailUser(ctx context.Context, userID string, requestBody body.EditEmailUserRequest) (*model.User, error) {
	value, err := u.userRepo.GetOTPValue(ctx, requestBody.Email)
	if err != nil {
		return nil, httperror.New(http.StatusBadRequest, response.OTPAlreadyExpiredMessage)
	}

	h := sha256.New()
	h.Write([]byte(value))
	hashedOTP := fmt.Sprintf("%x", h.Sum(nil))

	if hashedOTP != requestBody.Code {
		return nil, httperror.New(http.StatusBadRequest, response.OTPIsNotValidMessage)
	}

	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}
	}

	userModel.Email = requestBody.Email
	userModel.UpdatedAt.Time = time.Now()
	userModel.UpdatedAt.Valid = true
	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		if u.userRepo.UpdateUserEmail(ctx, tx, userModel) != nil {
			return err
		}

		if u.userRepo.CreateEmailHistory(ctx, tx, userModel.Email) != nil {
			return err
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	_, err = u.userRepo.DeleteOTPValue(ctx, requestBody.Email)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

func (u *userUC) SendLinkOTPEmail(ctx context.Context, email string) error {
	otp, err := util.GenerateOTP(6)
	if err != nil {
		return err
	}

	if err := u.userRepo.InsertNewOTPKey(ctx, email, otp); err != nil {
		return err
	}

	h := sha256.New()
	h.Write([]byte(otp))
	hashedOTP := fmt.Sprintf("%x", h.Sum(nil))

	link := fmt.Sprintf("http://localhost:8080/api/v1/user/email?code=%s&email=%s", hashedOTP, email)

	subject := "Change email!"
	msg := smtp.VerificationEmailLinkOTPBody(link)
	go smtp.SendEmail(u.cfg, email, subject, msg)

	return nil
}
