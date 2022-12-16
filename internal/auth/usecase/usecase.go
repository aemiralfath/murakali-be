package usecase

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"murakali/config"
	"murakali/internal/auth"
	"murakali/internal/auth/delivery/body"
	"murakali/internal/model"
	"murakali/internal/util"
	smtp "murakali/pkg/email"
	"murakali/pkg/httperror"
	"murakali/pkg/jwt"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"
)

type authUC struct {
	cfg      *config.Config
	txRepo   *postgre.TxRepo
	authRepo auth.Repository
}

func NewAuthUseCase(cfg *config.Config, txRepo *postgre.TxRepo, authRepo auth.Repository) auth.UseCase {
	return &authUC{cfg: cfg, txRepo: txRepo, authRepo: authRepo}
}

func (u *authUC) Login(ctx context.Context, requestBody body.LoginRequest) (string, string, error) {
	user, err := u.authRepo.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		return "", "", httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
	}

	accessToken, err := jwt.GenerateJWTAccessToken(user.ID.String(), user.RoleID, u.cfg)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateJWTRefreshToken(user.ID.String(), u.cfg)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *authUC) RegisterEmail(ctx context.Context, body body.RegisterEmailRequest) (*model.User, error) {
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

func (u *authUC) RegisterUser(ctx context.Context, body body.RegisterUserRequest) error {
	value, err := u.authRepo.GetOTPValue(ctx, body.Email)
	if err != nil {
		return httperror.New(http.StatusBadRequest, response.OTPAlreadyExpiredMessage)
	}

	if value != body.OTP {
		return httperror.New(http.StatusBadRequest, response.OTPIsNotValidMessage)
	}

	usernameUser, err := u.authRepo.GetUserByUsername(ctx, body.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if usernameUser != nil {
		return httperror.New(http.StatusBadRequest, response.UserNameAlreadyExistMessage)
	}

	phoneNoUser, err := u.authRepo.GetUserByPhoneNo(ctx, body.PhoneNo)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if phoneNoUser != nil {
		return httperror.New(http.StatusBadRequest, response.PhoneNoAlreadyExistMessage)
	}

	user, err := u.authRepo.GetUserByEmail(ctx, body.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if !user.IsVerify {
		user.PhoneNo = body.PhoneNo
		user.FullName = body.FullName
		user.Username = body.Username
		user.Password = string(hashedPassword)
		user.IsVerify = true
		user.UpdatedAt = time.Now()
		err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
			if err := u.authRepo.UpdateUser(ctx, tx, user); err != nil {
				return err
			}

			if err := u.authRepo.CreateEmailHistory(ctx, tx, body.Email); err != nil {
				return err
			}
			return err
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *authUC) VerifyOTP(ctx context.Context, body body.VerifyOTPRequest) error {
	value, err := u.authRepo.GetOTPValue(ctx, body.Email)
	if err != nil {
		return httperror.New(http.StatusBadRequest, response.OTPAlreadyExpiredMessage)
	}

	if value != body.OTP {
		return httperror.New(http.StatusBadRequest, response.OTPIsNotValidMessage)
	}

	return nil
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
	msg := smtp.VerificationEmailBody(otp)
	go smtp.SendEmail(u.cfg, email, subject, msg)

	return nil
}
