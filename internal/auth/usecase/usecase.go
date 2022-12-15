package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/config"
	"murakali/internal/auth"
	"murakali/internal/auth/delivery/body"
	"murakali/internal/model"
	"murakali/internal/util"
	smtp "murakali/pkg/email"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
)

type authUC struct {
	cfg      *config.Config
	authRepo auth.Repository
}

func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo}
}

func (u *authUC) Register(ctx context.Context, body body.RegisterRequest) (*model.User, error) {
	emailHistory, err := u.authRepo.CheckEmailHistory(ctx, body.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if emailHistory != nil {
		return nil, httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage)
	}

	user, err := u.authRepo.GetUserByEmail(ctx, body.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if user != nil {
		if user.IsVerify {
			return nil, httperror.New(http.StatusBadRequest, response.UserAlreadyExistMessage)
		} else {
			if err := u.SendOTPEmail(ctx, user.Email); err != nil {
				return nil, err
			}

			return user, nil
		}
	}

	user, err = u.authRepo.CreateUser(ctx, body.Email)
	if err != nil {
		return nil, err
	}

	if err := u.SendOTPEmail(ctx, user.Email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUC) SendOTPEmail(ctx context.Context, email string) error {
	otp, err := util.GenerateOTP(6)
	if err != nil {
		return err
	}

	if err := u.authRepo.InsertNewOTPKey(ctx, email, otp); err != nil {
		return err
	}

	subject := "Email Verification!"
	msg := fmt.Sprintf("<html><body>Hi!<br>Please input this OTP for verify your email: %s<br></body></html>", otp)
	go smtp.SendEmail(u.cfg, email, subject, msg)

	return nil
}
