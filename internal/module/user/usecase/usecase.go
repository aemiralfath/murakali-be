package usecase

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"math"
	"fmt"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/util"
	smtp "murakali/pkg/email"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
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

func (u *userUC) CreateAddress(ctx context.Context, userID string, requestBody body.CreateAddressRequest) error {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}

		return err
	}

	if requestBody.IsShopDefault && userModel.RoleID != constant.RoleSeller {
		return httperror.New(http.StatusBadRequest, response.UserNotASellerMessage)
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

	if requestBody.IsShopDefault && userModel.RoleID != constant.RoleSeller {
		return httperror.New(http.StatusBadRequest, response.UserNotASellerMessage)
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

	return nil
}

func (u *userUC) GetAddress(ctx context.Context, userID, name string, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}

		return nil, err
	}

	totalRows, err := u.userRepo.GetTotalAddress(ctx, userModel.ID.String(), name)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	addresses, err := u.userRepo.GetAddresses(ctx, userModel.ID.String(), name, pgn)
	if err != nil {
		return nil, err
	}

	pgn.Rows = addresses
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

	link := fmt.Sprintf("http://localhost:8080/api/v1/user/email?code=%s&email=%s", hashedOTP, email)

	subject := "Change email!"
	msg := smtp.VerificationEmailLinkOTPBody(link)
	go smtp.SendEmail(u.cfg, email, subject, msg)

	return nil
}

func (u *userUC) GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error) {

	response, err := u.userRepo.GetSealabsPay(ctx, userid)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *userUC) AddSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error {

	card_number, err := u.userRepo.CheckDefaultSealabsPay(ctx, userid)
	if err != nil {
		return err
	}

	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		if u.userRepo.SetDefaultSealabsPay(ctx, tx, card_number) != nil {
			return err
		}

		err = u.userRepo.AddSealabsPay(ctx, tx, request)
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}
