package usecase

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/auth"
	"murakali/internal/module/auth/delivery/body"
	"murakali/internal/util"
	smtp "murakali/pkg/email"
	"murakali/pkg/httperror"
	"murakali/pkg/jwt"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUC struct {
	cfg      *config.Config
	txRepo   *postgre.TxRepo
	authRepo auth.Repository
}

func NewAuthUseCase(cfg *config.Config, txRepo *postgre.TxRepo, authRepo auth.Repository) auth.UseCase {
	return &authUC{cfg: cfg, txRepo: txRepo, authRepo: authRepo}
}

func (u *authUC) Login(ctx context.Context, requestBody body.LoginRequest) (*model.Token, error) {
	user, err := u.authRepo.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}

		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(requestBody.Password)) != nil {
		return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
	}

	accessToken, err := jwt.GenerateJWTAccessToken(user.ID.String(), user.RoleID, u.cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateJWTRefreshToken(user.ID.String(), u.cfg)
	if err != nil {
		return nil, err
	}

	return &model.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (u *authUC) RefreshToken(ctx context.Context, id string) (*model.AccessToken, error) {
	user, err := u.authRepo.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}

		return nil, err
	}

	accessToken, err := jwt.GenerateJWTAccessToken(user.ID.String(), user.RoleID, u.cfg)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (u *authUC) RegisterEmail(ctx context.Context, requestBody body.RegisterEmailRequest) (*model.User, error) {
	emailHistory, err := u.authRepo.CheckEmailHistory(ctx, requestBody.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if emailHistory != nil {
		return nil, httperror.New(http.StatusBadRequest, response.EmailAlreadyExistMessage)
	}

	user, err := u.authRepo.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if user != nil {
		if user.IsVerify {
			return nil, httperror.New(http.StatusBadRequest, response.UserAlreadyExistMessage)
		}

		if errEmail := u.SendOTPEmail(ctx, user.Email); errEmail != nil {
			return nil, errEmail
		}

		return user, nil
	}

	user, err = u.authRepo.CreateUser(ctx, requestBody.Email)
	if err != nil {
		return nil, err
	}

	if err := u.SendOTPEmail(ctx, user.Email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUC) RegisterUser(ctx context.Context, email string, requestBody body.RegisterUserRequest) error {
	user, err := u.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return err
	}

	usernameUser, err := u.authRepo.GetUserByUsername(ctx, requestBody.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if usernameUser != nil {
		return httperror.New(http.StatusBadRequest, response.UserNameAlreadyExistMessage)
	}

	phoneNoUser, err := u.authRepo.GetUserByPhoneNo(ctx, requestBody.PhoneNo)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if phoneNoUser != nil {
		return httperror.New(http.StatusBadRequest, response.PhoneNoAlreadyExistMessage)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	password := string(hashedPassword)
	if !user.IsVerify {
		user.PhoneNo = &requestBody.PhoneNo
		user.FullName = &requestBody.FullName
		user.Username = &requestBody.Username
		user.Password = &password
		user.IsVerify = true
		user.UpdatedAt.Time = time.Now()
		user.UpdatedAt.Valid = true
		err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
			if errUpdate := u.authRepo.UpdateUser(ctx, tx, user); errUpdate != nil {
				return errUpdate
			}

			if errCreate := u.authRepo.CreateEmailHistory(ctx, tx, email); errCreate != nil {
				return errCreate
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *authUC) ResetPasswordEmail(ctx context.Context, requestBody body.ResetPasswordEmailRequest) (*model.User, error) {
	emailHistory, err := u.authRepo.CheckEmailHistory(ctx, requestBody.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if emailHistory == nil {
		return nil, httperror.New(http.StatusBadRequest, response.EmailNotExistMessage)
	}

	user, err := u.authRepo.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if user == nil {
		return nil, httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
	}

	if !user.IsVerify {
		return nil, httperror.New(http.StatusBadRequest, response.UserNotVerifyMessage)
	}

	if err := u.SendLinkOTPEmail(ctx, user.Email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUC) ResetPasswordUser(ctx context.Context, email string, requestBody *body.ResetPasswordUserRequest) (*model.User, error) {
	emailHistory, err := u.authRepo.CheckEmailHistory(ctx, email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if emailHistory == nil {
		return nil, httperror.New(http.StatusBadRequest, response.EmailNotExistMessage)
	}

	user, err := u.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if user == nil {
		return nil, httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
	}

	if !user.IsVerify {
		return nil, httperror.New(http.StatusBadRequest, response.UserNotVerifyMessage)
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(requestBody.Password))
	if err == nil {
		return nil, httperror.New(http.StatusBadRequest, response.PasswordSameOldPasswordMessage)
	}

	if strings.Contains(strings.ToLower(requestBody.Password), strings.ToLower(*user.Username)) {
		return nil, httperror.New(http.StatusBadRequest, response.PasswordContainUsernameMessage)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err = u.authRepo.UpdatePassword(ctx, user, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUC) VerifyOTP(ctx context.Context, requestBody body.VerifyOTPRequest) (string, error) {
	value, err := u.authRepo.GetOTPValue(ctx, requestBody.Email)
	if err != nil {
		return "", httperror.New(http.StatusBadRequest, response.OTPAlreadyExpiredMessage)
	}

	if value != requestBody.OTP {
		return "", httperror.New(http.StatusBadRequest, response.OTPIsNotValidMessage)
	}

	registerToken, err := jwt.GenerateJWTRegisterToken(requestBody.Email, u.cfg)
	if err != nil {
		return "", err
	}

	_, err = u.authRepo.DeleteOTPValue(ctx, requestBody.Email)
	if err != nil {
		return "", err
	}

	return registerToken, nil
}

func (u *authUC) ResetPasswordVerifyOTP(ctx context.Context, requestBody body.ResetPasswordVerifyOTPRequest) (string, error) {
	value, err := u.authRepo.GetOTPValue(ctx, requestBody.Email)
	if err != nil {
		return "", httperror.New(http.StatusBadRequest, response.OTPAlreadyExpiredMessage)
	}

	h := sha256.New()
	h.Write([]byte(value))
	hashedOTP := fmt.Sprintf("%x", h.Sum(nil))

	if hashedOTP != requestBody.Code {
		return "", httperror.New(http.StatusBadRequest, response.OTPIsNotValidMessage)
	}

	resetPasswordToken, err := jwt.GenerateJWTResetPasswordToken(requestBody.Email, hashedOTP, u.cfg)
	if err != nil {
		return "", err
	}

	_, err = u.authRepo.DeleteOTPValue(ctx, requestBody.Email)
	if err != nil {
		return "", err
	}

	return resetPasswordToken, nil
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

func (u *authUC) SendLinkOTPEmail(ctx context.Context, email string) error {
	otp, err := util.GenerateOTP(6)
	if err != nil {
		return err
	}

	if err := u.authRepo.InsertNewOTPKey(ctx, email, otp); err != nil {
		return err
	}

	h := sha256.New()
	h.Write([]byte(otp))
	hashedOTP := fmt.Sprintf("%x", h.Sum(nil))

	link := fmt.Sprintf("http://%s/verify?code=%s&email=%s", u.cfg.Server.Origin, hashedOTP, email)

	subject := "Reset Password!"
	msg := smtp.VerificationEmailLinkOTPBody(link)
	go smtp.SendEmail(u.cfg, email, subject, msg)

	return nil
}
