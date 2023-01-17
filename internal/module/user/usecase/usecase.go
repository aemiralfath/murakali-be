package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/util"
	smtp "murakali/pkg/email"
	"murakali/pkg/httperror"
	"murakali/pkg/jwt"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userUC struct {
	cfg      *config.Config
	txRepo   *postgre.TxRepo
	userRepo user.Repository
}

func NewUserUseCase(cfg *config.Config, txRepo *postgre.TxRepo, userRepo user.Repository) user.UseCase {
	return &userUC{cfg: cfg, txRepo: txRepo, userRepo: userRepo}
}

func (u *userUC) CreateAddress(ctx context.Context, userID string, requestBody body.CreateAddressRequest) error {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}

		return err
	}

	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		if requestBody.IsDefault {
			defaultAddress, getErr := u.userRepo.GetDefaultUserAddress(ctx, userModel.ID.String())
			if getErr != nil && getErr != sql.ErrNoRows {
				return getErr
			}

			if defaultAddress != nil && defaultAddress.IsDefault {
				if errUpdate := u.userRepo.UpdateDefaultAddress(ctx, tx, false, defaultAddress); errUpdate != nil {
					return errUpdate
				}
			}
		}

		if requestBody.IsShopDefault {
			defaultShopAddress, getErr := u.userRepo.GetDefaultShopAddress(ctx, userModel.ID.String())
			if getErr != nil && getErr != sql.ErrNoRows {
				return getErr
			}

			if defaultShopAddress != nil && defaultShopAddress.IsShopDefault {
				if errUpdate := u.userRepo.UpdateDefaultShopAddress(ctx, tx, false, defaultShopAddress); errUpdate != nil {
					return errUpdate
				}
			}
		}

		err = u.userRepo.CreateAddress(ctx, tx, userModel.ID.String(), requestBody)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (u *userUC) UpdateAddressByID(ctx context.Context, userID, addressID string, requestBody body.UpdateAddressRequest) error {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}

		return err
	}

	address, err := u.userRepo.GetAddressByID(ctx, userModel.ID.String(), addressID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.AddressNotExistMessage)
		}

		return err
	}

	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		if requestBody.IsDefault {
			defaultAddress, getErr := u.userRepo.GetDefaultUserAddress(ctx, userModel.ID.String())
			if getErr != nil && getErr != sql.ErrNoRows {
				return getErr
			}

			if defaultAddress != nil && defaultAddress.IsDefault {
				if errUpdate := u.userRepo.UpdateDefaultAddress(ctx, tx, false, defaultAddress); errUpdate != nil {
					return errUpdate
				}
			}
		}

		if requestBody.IsShopDefault {
			defaultShopAddress, getErr := u.userRepo.GetDefaultShopAddress(ctx, userModel.ID.String())
			if getErr != nil && getErr != sql.ErrNoRows {
				return getErr
			}

			if defaultShopAddress != nil && defaultShopAddress.IsShopDefault {
				if errUpdate := u.userRepo.UpdateDefaultShopAddress(ctx, tx, false, defaultShopAddress); errUpdate != nil {
					return errUpdate
				}
			}
		}

		address.Name = requestBody.Name
		address.ProvinceID = requestBody.ProvinceID
		address.CityID = requestBody.CityID
		address.Province = requestBody.Province
		address.City = requestBody.City
		address.District = requestBody.District
		address.SubDistrict = requestBody.SubDistrict
		address.AddressDetail = requestBody.AddressDetail
		address.ZipCode = requestBody.ZipCode
		address.IsDefault = requestBody.IsDefault
		address.IsShopDefault = requestBody.IsShopDefault
		address.UpdatedAt.Valid = true
		address.UpdatedAt.Time = time.Now()
		err = u.userRepo.UpdateAddress(ctx, tx, address)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *userUC) GetAddress(ctx context.Context, userID string, pgn *pagination.Pagination,
	queryRequest *body.GetAddressQueryRequest) (*pagination.Pagination, error) {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}
		return nil, err
	}

	addresses := make([]*model.Address, 0)
	if !queryRequest.IsDefaultBool && !queryRequest.IsShopDefaultBool {
		totalRows, err := u.userRepo.GetTotalAddress(ctx, userModel.ID.String(), queryRequest.Name)
		if err != nil {
			return nil, err
		}

		totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
		pgn.TotalRows = totalRows
		pgn.TotalPages = totalPages

		addresses, err = u.userRepo.GetAllAddresses(ctx, userModel.ID.String(), queryRequest.Name, pgn)
		if err != nil {
			return nil, err
		}
	}

	if queryRequest.IsDefaultBool && !queryRequest.IsShopDefaultBool {
		address, err := u.userRepo.GetDefaultUserAddress(ctx, userModel.ID.String())
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, httperror.New(http.StatusBadRequest, response.DefaultAddressNotFound)
			}
			return nil, err
		}

		addresses = append(addresses, address)

		totalRows := int64(1)
		totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
		pgn.TotalRows = totalRows
		pgn.TotalPages = totalPages
	}

	if !queryRequest.IsDefaultBool && queryRequest.IsShopDefaultBool {
		address, err := u.userRepo.GetDefaultShopAddress(ctx, userModel.ID.String())
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, httperror.New(http.StatusBadRequest, response.ShopAddressNotFound)
			}
			return nil, err
		}

		addresses = append(addresses, address)

		totalRows := int64(1)
		totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
		pgn.TotalRows = totalRows
		pgn.TotalPages = totalPages
	}

	pgn.Rows = addresses
	return pgn, nil
}

func (u *userUC) GetOrder(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	totalRows, err := u.userRepo.GetTotalOrder(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	orders, err := u.userRepo.GetOrders(ctx, userID, pgn)
	if err != nil {
		return nil, err
	}
	pgn.Rows = orders
	return pgn, nil
}

func (u *userUC) GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error) {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}

		return nil, err
	}

	address, err := u.userRepo.GetAddressByID(ctx, userModel.ID.String(), addressID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.AddressNotExistMessage)
		}

		return nil, err
	}

	return address, nil
}

func (u *userUC) DeleteAddressByID(ctx context.Context, userID, addressID string) error {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}

		return err
	}

	address, err := u.userRepo.GetAddressByID(ctx, userModel.ID.String(), addressID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.AddressNotExistMessage)
		}

		return err
	}

	if address.IsDefault || address.IsShopDefault {
		return httperror.New(http.StatusBadRequest, response.AddressIsDefaultMessage)
	}

	if err := u.userRepo.DeleteAddress(ctx, address.ID.String()); err != nil {
		return err
	}

	return nil
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

	link := fmt.Sprintf("%s/verify/email?code=%s&email=%s", u.cfg.Server.Origin, hashedOTP, email)

	subject := "Change email!"
	msg := smtp.VerificationEmailLinkOTPBody(link)
	go smtp.SendEmail(u.cfg, email, subject, msg)

	return nil
}

func (u *userUC) GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error) {
	slp, err := u.userRepo.GetSealabsPay(ctx, userid)
	if err != nil {
		return nil, err
	}

	return slp, nil
}

func (u *userUC) AddSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error {
	slpCount, err := u.userRepo.CheckUserSealabsPay(ctx, userid)
	if err != nil {
		return err
	}
	cardCount, err := u.userRepo.CheckDeletedSealabsPay(ctx, request.CardNumber)
	if err != nil {
		return err
	}

	if slpCount == 0 {
		if cardCount == 0 {
			err = u.userRepo.AddSealabsPay(ctx, request, userid)
			if err != nil {
				return err
			}
		} else {
			err = u.userRepo.UpdateUserSealabsPay(ctx, request, userid)
			if err != nil {
				return err
			}
		}
	} else {
		cardNumber, err := u.userRepo.CheckDefaultSealabsPay(ctx, userid)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if *cardNumber == request.CardNumber {
			return httperror.New(http.StatusBadRequest, response.SealabsCardAlreadyExist)
		}

		if cardCount == 0 {
			err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
				if u.userRepo.SetDefaultSealabsPayTrans(ctx, tx, cardNumber) != nil {
					return err
				}

				err = u.userRepo.AddSealabsPayTrans(ctx, tx, request, userid)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return httperror.New(http.StatusBadRequest, response.SealabsCardAlreadyExist)
			}
		} else {
			err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
				if u.userRepo.SetDefaultSealabsPayTrans(ctx, tx, cardNumber) != nil {
					return err
				}

				err = u.userRepo.UpdateUserSealabsPayTrans(ctx, tx, request, userid)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *userUC) PatchSealabsPay(ctx context.Context, cardNumber, userid string) error {
	err := u.userRepo.PatchSealabsPay(ctx, cardNumber)
	if err != nil {
		return err
	}

	if u.userRepo.SetDefaultSealabsPay(ctx, cardNumber, userid) != nil {
		return err
	}
	return nil
}

func (u *userUC) DeleteSealabsPay(ctx context.Context, cardNumber string) error {
	err := u.userRepo.DeleteSealabsPay(ctx, cardNumber)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUC) ActivateWallet(ctx context.Context, userID, pin string) error {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	wallet, err := u.userRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if wallet != nil {
		return httperror.New(http.StatusBadRequest, response.WalletAlreadyActivated)
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	walletData := &model.Wallet{}
	walletData.UserID = userModel.ID
	walletData.Balance = 0
	walletData.PIN = string(hashedPin)
	walletData.AttemptCount = 0
	walletData.ActiveDate.Valid = true
	walletData.ActiveDate.Time = time.Now()

	if err := u.userRepo.CreateWallet(ctx, walletData); err != nil {
		return err
	}

	return nil
}

func (u *userUC) RegisterMerchant(ctx context.Context, userID, shopName string) error {
	count, err := u.userRepo.CheckShopByID(ctx, userID)
	if err != nil {
		return err
	}
	if count != 0 {
		return httperror.New(http.StatusBadRequest, response.UserAlreadyHaveShop)
	}

	count, err = u.userRepo.CheckShopUnique(ctx, shopName)
	if err != nil {
		return err
	}
	if count != 0 {
		return httperror.New(http.StatusBadRequest, response.ShopAlreadyExists)
	}

	if _, errWallet := u.userRepo.GetWalletByUserID(ctx, userID); errWallet != nil {
		if errWallet == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.WalletIsNotActivated)
		}
		return errWallet
	}

	err = u.userRepo.AddShop(ctx, userID, shopName)
	if err != nil {
		return err
	}

	err = u.userRepo.UpdateRole(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUC) GetUserProfile(ctx context.Context, userID string) (*body.ProfileResponse, error) {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return nil, err
	}

	profileInfo := &body.ProfileResponse{
		Role:        userModel.RoleID,
		UserName:    userModel.Username,
		Email:       userModel.Email,
		PhoneNumber: userModel.PhoneNo,
		FullName:    userModel.FullName,
		Gender:      userModel.Gender,
		BirthDate:   userModel.BirthDate.Time,
		PhotoURL:    userModel.PhotoURL,
		IsVerify:    userModel.IsVerify,
	}

	return profileInfo, nil
}

func (u *userUC) UploadProfilePicture(ctx context.Context, imgURL, userID string) error {
	err := u.userRepo.UpdateProfileImage(ctx, imgURL, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUC) VerifyPasswordChange(ctx context.Context, userID string) error {
	userInfo, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	err = u.SendOTPEmail(ctx, userInfo.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUC) SendOTPEmail(ctx context.Context, email string) error {
	otp, err := util.GenerateRandomAlpaNumeric(6)
	if err != nil {
		return err
	}

	if err := u.userRepo.InsertNewOTPKey(ctx, email, otp); err != nil {
		return err
	}

	subject := "Change Password Verification!"
	msg := smtp.VerificationEmailBody(otp)
	go smtp.SendEmail(u.cfg, email, subject, msg)

	return nil
}

func (u *userUC) VerifyOTP(ctx context.Context, requestBody body.VerifyOTPRequest, userID string) (string, error) {
	userInfo, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	value, err := u.userRepo.GetOTPValue(ctx, userInfo.Email)
	if err != nil {
		return "", httperror.New(http.StatusBadRequest, response.OTPAlreadyExpiredMessage)
	}

	if value != requestBody.OTP {
		return "", httperror.New(http.StatusBadRequest, response.OTPIsNotValidMessage)
	}

	changePasswordToken, err := jwt.GenerateJWTChangePasswordToken(userInfo.ID.String(), u.cfg)
	if err != nil {
		return "", err
	}

	_, err = u.userRepo.DeleteOTPValue(ctx, userInfo.Email)
	if err != nil {
		return "", err
	}

	return changePasswordToken, nil
}

func (u *userUC) ChangePassword(ctx context.Context, userID, newPassword string) error {
	userInfo, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	userPassword, err := u.userRepo.GetPasswordByID(ctx, userID)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(newPassword))
	if err == nil {
		return httperror.New(http.StatusBadRequest, response.PasswordSameOldPasswordMessage)
	}

	if strings.Contains(strings.ToLower(newPassword), strings.ToLower(*userInfo.Username)) {
		return httperror.New(http.StatusBadRequest, response.PasswordContainUsernameMessage)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	password := string(hashedPassword)

	err = u.userRepo.UpdatePasswordByID(ctx, userID, password)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUC) TopUpWallet(ctx context.Context, userID string, requestBody body.TopUpWalletRequest) (string, error) {
	wallet, err := u.userRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", httperror.New(http.StatusBadRequest, response.WalletIsNotActivated)
		}
		return "", err
	}

	card, err := u.userRepo.GetSealabsPayUser(ctx, userID, requestBody.CardNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", httperror.New(http.StatusBadRequest, response.SealabsCardNotFound)
		}
		return "", err
	}

	transactionID, err := u.txRepo.WithTransactionReturnData(func(tx postgre.Transaction) (interface{}, error) {
		transaction := &model.Transaction{}
		transaction.WalletID = &wallet.ID
		transaction.CardNumber = &card.CardNumber
		transaction.TotalPrice = float64(requestBody.Amount)
		transaction.ExpiredAt.Valid = true
		transaction.ExpiredAt.Time = time.Now().Add(time.Hour * 24)

		transactionID, errTrans := u.userRepo.CreateTransaction(ctx, tx, transaction)
		if errTrans != nil {
			return nil, errTrans
		}

		return transactionID.String(), nil
	})

	if err != nil {
		return "", err
	}

	return transactionID.(string), nil
}

func (u *userUC) CreateSLPPayment(ctx context.Context, transactionID string) (string, error) {
	transaction, err := u.userRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", httperror.New(http.StatusBadRequest, response.TransactionIDNotExist)
		}

		return "", err
	}

	if time.Until(transaction.ExpiredAt.Time) < 0 {
		return "", httperror.New(http.StatusBadRequest, response.TransactionAlreadyExpired)
	}

	if transaction.PaidAt.Valid || transaction.CanceledAt.Valid {
		return "", httperror.New(http.StatusBadRequest, response.TransactionAlreadyFinished)
	}

	if transaction.CardNumber == nil {
		return "", httperror.New(http.StatusBadRequest, response.InvalidPaymentMethod)
	}

	signFormat := fmt.Sprintf("%s:%d:%s", *transaction.CardNumber, int(transaction.TotalPrice), u.cfg.External.SlpMerchantCode)
	h := hmac.New(sha256.New, []byte(u.cfg.External.SlpAPIKey))
	h.Write([]byte(signFormat))
	sign := hex.EncodeToString(h.Sum(nil))

	redirectURL, err := u.GetRedirectURL(transaction, sign)
	if err != nil {
		return "", err
	}

	return redirectURL, nil
}

func (u *userUC) CreateWalletPayment(ctx context.Context, transactionID string) error {
	transaction, err := u.userRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.TransactionIDNotExist)
		}

		return err
	}

	if time.Until(transaction.ExpiredAt.Time) < 0 {
		return httperror.New(http.StatusBadRequest, response.TransactionAlreadyExpired)
	}

	if transaction.PaidAt.Valid || transaction.CanceledAt.Valid {
		return httperror.New(http.StatusBadRequest, response.TransactionAlreadyFinished)
	}

	if transaction.WalletID == nil {
		return httperror.New(http.StatusBadRequest, response.InvalidPaymentMethod)
	}

	wallet, err := u.userRepo.GetWalletUser(ctx, transaction.WalletID.String())
	if err != nil {
		return err
	}

	if wallet.Balance-transaction.TotalPrice < 0 {
		return httperror.New(http.StatusBadRequest, response.WalletBalanceNotEnough)
	}

	orders, err := u.userRepo.GetOrderByTransactionID(ctx, transactionID)
	if err != nil {
		return err
	}

	errTx := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		transaction.PaidAt.Valid = true
		transaction.PaidAt.Time = time.Now()
		if errTransaction := u.userRepo.UpdateTransaction(ctx, tx, transaction); errTransaction != nil {
			return errTransaction
		}

		for _, order := range orders {
			order.OrderStatusID = constant.OrderStatusWaitingForSeller
			if errOrder := u.userRepo.UpdateOrder(ctx, tx, order); errOrder != nil {
				return errOrder
			}
		}

		walletHistory := &model.WalletHistory{}
		walletHistory.TransactionID = transaction.ID
		walletHistory.WalletID = wallet.ID
		walletHistory.From = wallet.ID.String()
		walletHistory.To = transaction.ID.String()
		walletHistory.Description = "Payment transaction " + transaction.ID.String()
		walletHistory.Amount = transaction.TotalPrice
		walletHistory.CreatedAt = time.Now()
		if errWallet := u.userRepo.InsertWalletHistory(ctx, tx, walletHistory); errWallet != nil {
			return errWallet
		}

		wallet.Balance -= transaction.TotalPrice
		wallet.UpdatedAt.Valid = true
		wallet.UpdatedAt.Time = time.Now()

		if errBalance := u.userRepo.UpdateWalletBalance(ctx, tx, wallet); errBalance != nil {
			return errBalance
		}

		if errCredit := u.CreditToMarketplaceAccount(ctx, tx, transaction); errCredit != nil {
			return errCredit
		}

		return nil
	})
	if errTx != nil {
		return err
	}

	return nil
}

func (u *userUC) GetRedirectURL(transaction *model.Transaction, sign string) (string, error) {
	var responseSLP body.SLPPaymentResponse

	url := fmt.Sprintf("%s/v1/transaction/pay", u.cfg.External.SlpURL)
	callbackURL := fmt.Sprintf("https://%s/api/v1/user/transaction/slp-payment/%s", u.cfg.Server.Domain, transaction.ID.String())
	if transaction.WalletID != nil {
		callbackURL = fmt.Sprintf("https://%s/api/v1/user/transaction/wallet-payment/%s", u.cfg.Server.Domain, transaction.ID.String())
	}

	payload := fmt.Sprintf(
		"card_number=%s&amount=%d&merchant_code=%s&redirect_url=%s&callback_url=%s&signature=%s",
		*transaction.CardNumber,
		int(transaction.TotalPrice),
		u.cfg.External.SlpMerchantCode,
		"https://www.google.com",
		callbackURL,
		sign)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 303 {
		readErr := json.NewDecoder(res.Body).Decode(&responseSLP)
		if readErr != nil {
			return "", err
		}

		return "", httperror.New(res.StatusCode, responseSLP.Message)
	}

	return res.Header.Get("Location"), nil
}

func (u *userUC) UpdateTransaction(ctx context.Context, transactionID string, requestBody body.SLPCallbackRequest) error {
	transaction, err := u.userRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.TransactionIDNotExist)
		}

		return err
	}

	if transaction.PaidAt.Valid || transaction.CanceledAt.Valid {
		return httperror.New(http.StatusBadRequest, response.TransactionAlreadyFinished)
	}

	orders, err := u.userRepo.GetOrderByTransactionID(ctx, transactionID)
	if err != nil {
		return err
	}

	if requestBody.Status == constant.SLPStatusPaid && requestBody.Message == constant.SlPMessagePaid {
		err := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
			transaction.PaidAt.Valid = true
			transaction.PaidAt.Time = time.Now()
			if err := u.userRepo.UpdateTransaction(ctx, tx, transaction); err != nil {
				return err
			}

			for _, order := range orders {
				order.OrderStatusID = constant.OrderStatusWaitingForSeller
				if err := u.userRepo.UpdateOrder(ctx, tx, order); err != nil {
					return err
				}
			}

			if err := u.CreditToMarketplaceAccount(ctx, tx, transaction); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (u *userUC) UpdateWalletTransaction(ctx context.Context, transactionID string, requestBody body.SLPCallbackRequest) error {
	transaction, err := u.userRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.TransactionIDNotExist)
		}

		return err
	}

	wallet, err := u.userRepo.GetWalletUser(ctx, transaction.WalletID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.WalletIsNotActivated)
		}

		return err
	}

	if transaction.PaidAt.Valid || transaction.CanceledAt.Valid {
		return httperror.New(http.StatusBadRequest, response.TransactionAlreadyFinished)
	}

	if requestBody.Status == constant.SLPStatusCanceled && requestBody.Message == constant.SLPMessageCanceled {
		err := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
			transaction.CanceledAt.Valid = true
			transaction.CanceledAt.Time = time.Now()
			if err := u.userRepo.UpdateTransaction(ctx, tx, transaction); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if requestBody.Status == constant.SLPStatusPaid && requestBody.Message == constant.SlPMessagePaid {
		err := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
			transaction.PaidAt.Valid = true
			transaction.PaidAt.Time = time.Now()
			if err := u.userRepo.UpdateTransaction(ctx, tx, transaction); err != nil {
				return err
			}

			walletHistory := &model.WalletHistory{}
			walletHistory.TransactionID = transaction.ID
			walletHistory.WalletID = wallet.ID
			walletHistory.From = *transaction.CardNumber
			walletHistory.To = transaction.WalletID.String()
			walletHistory.Description = "Top up from " + *transaction.CardNumber
			walletHistory.Amount = transaction.TotalPrice
			walletHistory.CreatedAt = time.Now()
			if err := u.userRepo.InsertWalletHistory(ctx, tx, walletHistory); err != nil {
				return err
			}

			wallet.Balance += transaction.TotalPrice
			wallet.UpdatedAt.Valid = true
			wallet.UpdatedAt.Time = time.Now()

			if err := u.userRepo.UpdateWalletBalance(ctx, tx, wallet); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (u *userUC) CreditToMarketplaceAccount(ctx context.Context, tx postgre.Transaction, transaction *model.Transaction) error {
	walletMarketplace, err := u.userRepo.GetWalletByUserID(ctx, constant.AdminMarketplaceID)
	if err != nil {
		return err
	}

	walletHistory := &model.WalletHistory{}
	if transaction.WalletID == nil {
		walletHistory.From = *transaction.CardNumber
	}

	if transaction.CardNumber == nil {
		walletHistory.From = transaction.WalletID.String()
	}

	walletHistory.TransactionID = transaction.ID
	walletHistory.WalletID = walletMarketplace.ID
	walletHistory.To = walletMarketplace.ID.String()
	walletHistory.Description = "Payment transaction " + transaction.ID.String()
	walletHistory.Amount = transaction.TotalPrice
	walletHistory.CreatedAt = time.Now()

	if err := u.userRepo.InsertWalletHistory(ctx, tx, walletHistory); err != nil {
		return err
	}

	walletMarketplace.Balance += transaction.TotalPrice
	walletMarketplace.UpdatedAt.Valid = true
	walletMarketplace.UpdatedAt.Time = time.Now()

	if err := u.userRepo.UpdateWalletBalance(ctx, tx, walletMarketplace); err != nil {
		return err
	}

	return nil
}

func (u *userUC) GetWallet(ctx context.Context, userID string) (*model.Wallet, error) {
	wallet, err := u.userRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.WalletIsNotActivated)
		}

		return nil, err
	}

	return wallet, nil
}

func (u *userUC) GetWalletHistory(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	wallet, err := u.userRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	totalRows, err := u.userRepo.GetTotalWalletHistoryByWalletID(ctx, wallet.ID.String())
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	walletHistory, err := u.userRepo.GetWalletHistoryByWalletID(ctx, pgn, wallet.ID.String())
	if err != nil {
		return nil, err
	}

	pgn.Rows = walletHistory

	return pgn, nil
}

func (u *userUC) WalletStepUp(ctx context.Context, userID string, requestBody body.WalletStepUpRequest) (string, error) {
	wallet, err := u.userRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", httperror.New(http.StatusBadRequest, response.WalletIsNotActivated)
		}
		return "", err
	}

	if wallet.Balance-float64(requestBody.Amount) < 0 {
		return "", httperror.New(http.StatusBadRequest, response.WalletBalanceNotEnough)
	}

	if wallet.UnlockedAt.Valid && time.Until(wallet.UnlockedAt.Time) >= 0 {
		return "", httperror.New(http.StatusBadRequest, response.WalletIsBlocked)
	}

	blocked := false
	invalidPin := false
	if bcrypt.CompareHashAndPassword([]byte(wallet.PIN), []byte(requestBody.Pin)) != nil {
		invalidPin = true
		wallet.AttemptCount++
		wallet.AttemptAt.Valid = true
		wallet.AttemptAt.Time = time.Now()

		if wallet.AttemptCount >= 3 {
			blocked = true
			wallet.AttemptCount = 0
			wallet.UnlockedAt.Valid = true
			wallet.UnlockedAt.Time = time.Now().Add(time.Minute * 15)
		}
	}

	if !invalidPin {
		wallet.AttemptCount = 0
		wallet.AttemptAt.Valid = true
		wallet.AttemptAt.Time = time.Now()
	}

	if errWallet := u.userRepo.UpdateWallet(ctx, wallet); errWallet != nil {
		return "", errWallet
	}

	if blocked {
		return "", httperror.New(http.StatusBadRequest, response.WalletIsBlocked)
	}

	if invalidPin {
		return "", httperror.New(http.StatusBadRequest, response.WalletPinIsInvalid)
	}

	walletToken, err := jwt.GenerateJWTWalletToken(userID, u.cfg)
	if err != nil {
		return "", err
	}

	return walletToken, nil
}

func (u *userUC) CreateTransaction(ctx context.Context, userID string, requestBody body.CreateTransactionRequest) (string, error) {
	// TODO: Add voucher promotion stock & validation
	transactionData := &model.Transaction{}
	orderResponses := make([]*body.OrderResponse, 0)

	userModel, errUser := u.userRepo.GetUserByID(ctx, userID)
	if errUser != nil {
		if errUser == sql.ErrNoRows {
			return "", httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}
		return "", errUser
	}

	if requestBody.WalletID != "" {
		walletUser, errWallet := u.userRepo.GetWalletUser(ctx, requestBody.WalletID)
		if errWallet != nil {
			if errWallet == sql.ErrNoRows {
				return "", httperror.New(http.StatusBadRequest, response.WalletIsNotActivated)
			}
			return "", errWallet
		}

		if walletUser.UserID != userModel.ID {
			return "", httperror.New(http.StatusBadRequest, response.WalletIsNotActivated)
		}

		transactionData.WalletID = &walletUser.ID
	}

	if requestBody.CardNumber != "" {
		SealabsPayUser, errSealabpay := u.userRepo.GetSealabsPayUser(ctx, userModel.ID.String(), requestBody.CardNumber)
		if errSealabpay != nil {
			return "", errSealabpay
		}
		transactionData.CardNumber = &SealabsPayUser.CardNumber
	}

	if requestBody.VoucherMarketplaceID != "" {
		voucherMarketplace, errVoucherMP := u.userRepo.GetVoucherMarketplaceByID(ctx, requestBody.VoucherMarketplaceID)
		if errVoucherMP != nil {
			if errVoucherMP != sql.ErrNoRows {
				return "", errVoucherMP
			}
		}
		transactionData.VoucherMarketplaceID = &voucherMarketplace.ID
	}

	data, err := u.txRepo.WithTransactionReturnData(func(tx postgre.Transaction) (interface{}, error) {
		for _, cart := range requestBody.CartItems {
			orderData := &model.OrderModel{}
			cartShop, err := u.userRepo.GetShopByID(ctx, cart.ShopID)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, httperror.New(http.StatusBadRequest, response.UnknownShop)
				}
				return nil, err
			}

			var voucherShop *model.Voucher
			var voucherShopID *uuid.UUID
			if cart.VoucherShopID != "" {
				voucherShop, err = u.userRepo.GetVoucherShopByID(ctx, cart.VoucherShopID, cartShop.ID.String())
				if err != nil {
					if err != sql.ErrNoRows {
						return nil, err
					}
				}
				voucherShopID = &voucherShop.ID
			}

			courierShop, err := u.userRepo.GetCourierShopByID(ctx, cart.CourierID, cartShop.ID.String())
			if err != nil {
				return nil, httperror.New(http.StatusBadRequest, response.SelectShippingCourier)
			}
			orderResponse := &body.OrderResponse{
				Items: make([]*body.OrderItemResponse, 0),
			}

			isAvail := true
			for _, bodyProductDetail := range cart.ProductDetails {
				productDetailData, err := u.userRepo.GetProductDetailByID(ctx, tx, bodyProductDetail.ID)
				if err != nil {
					return nil, err
				}

				cartData, err := u.userRepo.GetCartItemUser(ctx, userModel.ID.String(), productDetailData.ID.String())
				if err != nil {
					return nil, httperror.New(http.StatusBadRequest, response.CartItemNotExist)
				}

				if int(productDetailData.Stock)-bodyProductDetail.Quantity < 0 {
					isAvail = false
					errCart := u.userRepo.DeleteCartItemByID(ctx, tx, cartData)
					if errCart != nil {
						return nil, errCart
					}
				}

				orderItem := &model.OrderItem{
					ProductDetailID: productDetailData.ID,
					Quantity:        bodyProductDetail.Quantity,
					ItemPrice:       bodyProductDetail.SubPrice,
					TotalPrice:      bodyProductDetail.SubPrice,
				}
				item := &body.OrderItemResponse{
					Item:              orderItem,
					ProductDetailData: productDetailData,
					CartItemData:      cartData,
				}
				orderResponse.Items = append(orderResponse.Items, item)
				orderData.TotalPrice += orderItem.TotalPrice
			}

			if !isAvail {
				return nil, httperror.New(http.StatusBadRequest, response.ProductQuantityNotAvailable)
			}

			orderData.ShopID = cartShop.ID
			orderData.UserID = userModel.ID
			orderData.VoucherShopID = voucherShopID
			orderData.CourierID = courierShop.ID
			orderData.DeliveryFee = cart.CourierFee
			orderData.OrderStatusID = 1

			orderResponse.OrderData = orderData
			transactionData.TotalPrice += orderData.TotalPrice + orderData.DeliveryFee

			orderResponses = append(orderResponses, orderResponse)
		}

		transactionData.ExpiredAt.Valid = true
		transactionData.ExpiredAt.Time = time.Now().Add(time.Hour * 24)
		transactionResponse := &body.TransactionResponse{
			TransactionData: transactionData,
			OrderResponses:  orderResponses,
		}

		transactionID, errTrans := u.userRepo.CreateTransaction(ctx, tx, transactionResponse.TransactionData)
		if errTrans != nil {
			return nil, errTrans
		}

		for _, o := range transactionResponse.OrderResponses {
			o.OrderData.TransactionID = *transactionID
			orderID, errOrder := u.userRepo.CreateOrder(ctx, tx, o.OrderData)
			if errOrder != nil {
				return nil, errOrder
			}
			for _, i := range o.Items {
				i.Item.OrderID = *orderID
				_, errItem := u.userRepo.CreateOrderItem(ctx, tx, i.Item)
				if errItem != nil {
					return nil, errItem
				}
				i.ProductDetailData.Stock -= i.CartItemData.Quantity
				errProduct := u.userRepo.UpdateProductDetailStock(ctx, tx, i.ProductDetailData)
				if errProduct != nil {
					return nil, errProduct
				}
				errCart := u.userRepo.DeleteCartItemByID(ctx, tx, i.CartItemData)
				if errCart != nil {
					return nil, errCart
				}
			}
		}
		return transactionID.String(), nil
	})
	if err != nil {
		return "", err
	}
	return data.(string), nil
}
