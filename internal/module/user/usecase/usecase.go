package usecase

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
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

func (u *userUC) GetAddress(ctx context.Context, userID string, pgn *pagination.Pagination, queryRequest *body.GetAddressQueryRequest) (*pagination.Pagination, error) {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
		}

		return nil, err
	}

	totalRows, err := u.userRepo.GetTotalAddress(ctx, userModel.ID.String(), queryRequest.Name, queryRequest.IsDefaultBool, queryRequest.IsShopDefaultBool)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	var addresses []*model.Address
	addresses, err = u.userRepo.GetAddresses(ctx, userModel.ID.String(), queryRequest.Name, queryRequest.IsDefaultBool, queryRequest.IsShopDefaultBool, pgn)
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

	link := fmt.Sprintf("http://%s/verify/email?code=%s&email=%s", u.cfg.Server.Origin, hashedOTP, email)

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
	cardNumber, err := u.userRepo.CheckDefaultSealabsPay(ctx, userid)
	if err != nil {
		return err
	}

	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		if u.userRepo.SetDefaultSealabsPayTrans(ctx, tx, cardNumber) != nil {
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
	otp, err := util.GenerateOTP(6)
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
