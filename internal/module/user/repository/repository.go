package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/postgre"

	"github.com/go-redis/redis/v8"
	"murakali/pkg/pagination"
	"time"
)

type userRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewUserRepository(psql *sql.DB, client *redis.Client) user.Repository {
	return &userRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}

func (r *userRepo) CreateAddress(ctx context.Context, tx postgre.Transaction, userID string, requestBody body.CreateAddressRequest) error {
	_, err := tx.ExecContext(
		ctx,
		CreateAddressQuery,
		userID,
		requestBody.Name,
		requestBody.ProvinceID,
		requestBody.CityID,
		requestBody.Province,
		requestBody.City,
		requestBody.District,
		requestBody.SubDistrict,
		requestBody.AddressDetail,
		requestBody.ZipCode,
		requestBody.IsDefault,
		requestBody.IsShopDefault)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetAddressByIDQuery, addressID, userID).Scan(
		&address.ID,
		&address.UserID,
		&address.Name,
		&address.ProvinceID,
		&address.CityID,
		&address.Province,
		&address.City,
		&address.District,
		&address.SubDistrict,
		&address.AddressDetail,
		&address.ZipCode,
		&address.IsDefault,
		&address.IsShopDefault,
		&address.CreatedAt,
		&address.UpdatedAt); err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *userRepo) GetDefaultUserAddress(ctx context.Context, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetDefaultAddressQuery, userID, true).Scan(&address.ID, &address.UserID, &address.IsDefault); err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *userRepo) GetDefaultShopAddress(ctx context.Context, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetDefaultShopAddressQuery, userID, true).
		Scan(&address.ID, &address.UserID, &address.IsShopDefault); err != nil {
		return nil, err
	}

	return &address, nil
}


func (r *userRepo) GetTotalAddress(ctx context.Context, userID, name string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalAddressQuery, userID, fmt.Sprintf("%%%s%%", name)).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepo) UpdateAddress(ctx context.Context, tx postgre.Transaction, address *model.Address) error {
	_, err := tx.ExecContext(ctx, UpdateAddressByIDQuery,
		address.Name,
		address.ProvinceID,
		address.CityID,
		address.Province,
		address.City,
		address.District,
		address.SubDistrict,
		address.AddressDetail,
		address.ZipCode,
		address.IsDefault,
		address.IsShopDefault,
		address.UpdatedAt,
		address.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateDefaultAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error {
	_, err := tx.ExecContext(ctx, UpdateDefaultAddressQuery, status, address.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateDefaultShopAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error {
	_, err := tx.ExecContext(ctx, UpdateDefaultShopAddressQuery, status, address.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) DeleteAddress(ctx context.Context, addressID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteAddressByIDQuery, addressID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetAddresses(ctx context.Context, userID, name string, pgn *pagination.Pagination) ([]*model.Address, error) {
	addresses := make([]*model.Address, 0)
	res, err := r.PSQL.QueryContext(
		ctx, GetAddressesQuery,
		userID,
		fmt.Sprintf("%%%s%%", name),
		pgn.GetSort(),
		pgn.GetLimit(),
		pgn.GetOffset())

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var address model.Address
		if errScan := res.Scan(
			&address.ID,
			&address.UserID,
			&address.Name,
			&address.ProvinceID,
			&address.CityID,
			&address.Province,
			&address.City,
			&address.District,
			&address.SubDistrict,
			&address.AddressDetail,
			&address.ZipCode,
			&address.IsDefault,
			&address.IsShopDefault,
			&address.CreatedAt,
			&address.UpdatedAt); errScan != nil {
			return nil, err
		}

		addresses = append(addresses, &address)
	}

	if res.Err() != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByIDQuery, id).
		Scan(&userModel.ID, &userModel.RoleID, &userModel.Email, &userModel.Username, &userModel.PhoneNo,
			&userModel.FullName, &userModel.Gender, &userModel.BirthDate, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *userRepo) CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error) {
	var emailHistory model.EmailHistory
	if err := r.PSQL.QueryRowContext(ctx, CheckEmailHistoryQuery, email).
		Scan(&emailHistory.ID, &emailHistory.Email); err != nil {
		return nil, err
	}

	return &emailHistory, nil
}
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByUsernameQuery, username).
		Scan(&userModel.ID, &userModel.Email, &userModel.Username, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *userRepo) GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByPhoneNoQuery, phoneNo).
		Scan(&userModel.ID, &userModel.Email, &userModel.PhoneNo, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}
func (r *userRepo) UpdateUserField(ctx context.Context, userModel *model.User) error {
	_, err := r.PSQL.ExecContext(
		ctx, UpdateUserFieldQuery, userModel.Username, userModel.FullName,
		userModel.PhoneNo, userModel.BirthDate, userModel.Gender, userModel.UpdatedAt, userModel.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateUserEmail(ctx context.Context, tx postgre.Transaction, userModel *model.User) error {
	_, err := tx.ExecContext(
		ctx, UpdateUserEmailQuery, userModel.Email, userModel.UpdatedAt, userModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error {
	_, err := tx.ExecContext(ctx, CreateEmailHistoryQuery, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) InsertNewOTPKey(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	duration, err := time.ParseDuration(constant.OtpDuration)
	if err != nil {
		return err
	}

	if err := r.RedisClient.Set(ctx, key, otp, duration); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *userRepo) GetOTPValue(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	res := r.RedisClient.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *userRepo) DeleteOTPValue(ctx context.Context, email string) (int64, error) {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	res := r.RedisClient.Del(ctx, key)
	if res.Err() != nil {
		return -1, res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return -1, err
	}

	return value, nil
}

func (r *userRepo) GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error) {

	responses := make([]*model.SealabsPay, 0)
	res, err := r.PSQL.QueryContext(ctx, GetSealabsPayByIdQuery, userid)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var response model.SealabsPay
		if err := res.Scan(
			&response.CardNumber, &response.UserID, &response.Name, &response.IsDefault, &response.ActiveDate, &response.CreatedAt, &response.UpdatedAt, &response.DeletedAt); err != nil {
			return nil, err
		}
		responses = append(responses, &response)
	}
	return responses, nil
}

func (r *userRepo) CheckDefaultSealabsPay(ctx context.Context, userid string) (*string, error) {
	var temp *string
	if err := r.PSQL.QueryRowContext(ctx, CheckDefaultSealabsPayQuery, userid).
		Scan(&temp); err != nil {
		return nil, err
	}

	return temp, nil
}

func (r *userRepo) SetDefaultSealabsPayTrans(ctx context.Context, tx postgre.Transaction, card_number *string) error {
	if _, err := tx.ExecContext(ctx, SetDefaultSealabsPayQuery, card_number); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) AddSealabsPay(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest) error {

	if _, err := tx.ExecContext(ctx, CreateSealabsPayQuery, request.CardNumber, request.UserID, request.Name, request.IsDefault, request.ActiveDateTime); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) SetDefaultSealabsPay(ctx context.Context, card_number string, userid string) error {
	if _, err := r.PSQL.ExecContext(ctx, SetDefaultSealabsPayQuery, card_number, userid); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) PatchSealabsPay(ctx context.Context, card_number string) error {

	if _, err := r.PSQL.ExecContext(ctx, PatchSealabsPayQuery, card_number); err != nil {
		return err
	}
	return nil
}

func (r *userRepo) DeleteSealabsPay(ctx context.Context, card_number string) error {

	if _, err := r.PSQL.ExecContext(ctx, DeleteSealabsPayQuery, card_number); err != nil {
		return err
	}
	return nil
}
