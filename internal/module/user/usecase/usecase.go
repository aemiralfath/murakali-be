package usecase

import (
	"context"
	"database/sql"
	"math"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
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

func (u *userUC) UpdateAddress(ctx context.Context, userID string, requestBody body.UpdateAddressRequest) error {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}

		return err
	}

	address, err := u.userRepo.GetAddressByID(ctx, requestBody.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.AddressNotExistMessage)
		}

		return err
	}

	if address.UserID != userModel.ID {
		return httperror.New(http.StatusBadRequest, response.UserNotMatchMessage)
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
