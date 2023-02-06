package table

import (
	"database/sql"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/pkg/postgre"
	"strconv"
	"time"
)

const InsertEmailHistoryQuery = `INSERT INTO "email_history" (email) VALUES ($1)`
const InsertUserQuery = `INSERT INTO "user" 
    	(id, role_id, username, email, phone_no, fullname, password, gender, photo_url, birth_date, is_sso, is_verify)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
const InsertUserAddressQuery = `INSERT INTO "address" 
    (user_id, name, province_id, city_id, province, city, district, sub_district, address_detail, zip_code, is_default, is_shop_default)
VALUES ($1, 'Home', 33, 327, 'Sumatera Selatan', 'Palembang', 'Ilir Timur II', '2 Ilir', 'no 91', '30118', $2, $3)`
const CreateUserWallet = `INSERT INTO "wallet" (user_id, balance, pin, attempt_count, active_date)
	VALUES ($1, 0, '$2a$10$haIhdlIHObH0yGMyCx5Zl.s5b7sV/x3GWact0Yd2xREXof3UAzUl6', 0, CURRENT_TIMESTAMP);`
const CreateUserSLP = `INSERT INTO "sealabs_pay" (card_number, user_id, name, is_default, active_date, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

type UserFaker struct {
	Size       int
	RoleID     int
	Gender     string
	UserID     []string
	Email      []string
	CardNumber []string
}

func NewUserFaker(size, roleID int, gender string, userID, email, cardNumber []string) ISeeder {
	return &UserFaker{Size: size, RoleID: roleID, Gender: gender, UserID: userID, Email: email, CardNumber: cardNumber}
}

func (u *UserFaker) GenerateData(tx postgre.Transaction) error {
	for i, val := range u.Email {
		id, err := uuid.Parse(u.UserID[i])
		if err != nil {
			return err
		}

		if err := u.GenerateDataUser(tx, id, val, u.CardNumber[i], u.CardNumber[i]); err != nil {
			return err
		}
	}

	for i := 0; i < u.Size; i++ {
		iString := strconv.Itoa(i)
		if err := u.GenerateDataUser(tx, uuid.New(), faker.Email(), iString, ""); err != nil {
			return err
		}
	}

	return nil
}

func (u *UserFaker) GenerateDataUser(tx postgre.Transaction, id uuid.UUID, email, phoneNo, cardNumber string) error {
	data := u.GenerateUser(id, email, phoneNo)
	_, err := tx.Exec(InsertUserQuery,
		data.ID, data.RoleID, data.Username, data.Email, data.PhoneNo, data.FullName, data.Password, data.Gender,
		data.PhotoURL, data.BirthDate, data.IsSSO, data.IsVerify)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(CreateUserWallet, id); err != nil {
		return err
	}

	if cardNumber != "" {
		if _, err := tx.Exec(CreateUserSLP, cardNumber, data.ID, faker.Name(), true, time.Now().AddDate(1, 0, 0), time.Now()); err != nil {
			return err
		}
	}

	if errEmail := u.GenerateEmailHistory(tx, data.Email); errEmail != nil {
		return errEmail
	}

	if errAddress := u.GenerateUserAddress(tx, data.ID.String()); errAddress != nil {
		return errAddress
	}

	return nil
}

func (u *UserFaker) GenerateEmailHistory(tx postgre.Transaction, email string) error {
	_, err := tx.Exec(InsertEmailHistoryQuery, email)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserFaker) GenerateUserAddress(tx postgre.Transaction, userID string) error {
	isShopDefault := false
	if u.RoleID == constant.RoleSeller {
		isShopDefault = true
	}
	_, err := tx.Exec(InsertUserAddressQuery, userID, true, isShopDefault)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserFaker) GenerateUser(id uuid.UUID, email, phoneNo string) *model.User {
	fullName := faker.Name()
	username := faker.Username()
	password := "$2a$10$cNhdZVN.pgsfK1xUQ00p7eK5Fh7iClrtJB9SY5un.H55Mi/dtQzCa"
	photoURL := "https://res.cloudinary.com/dhpao1zxi/image/upload/v1675678498/seeder_mgcsz6.png"

	return &model.User{
		ID:        id,
		RoleID:    u.RoleID,
		Username:  &username,
		PhoneNo:   &phoneNo,
		Email:     email,
		FullName:  &fullName,
		Password:  &password,
		Gender:    &u.Gender,
		PhotoURL:  &photoURL,
		BirthDate: sql.NullTime{Valid: true, Time: time.Now()},
		IsSSO:     false,
		IsVerify:  true,
	}
}
