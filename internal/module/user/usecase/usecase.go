package usecase

import (
	"context"
	"database/sql"
	"math"
	"murakali/config"
	"murakali/internal/module/user"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
)

type userUC struct {
	cfg      *config.Config
	txRepo   *postgre.TxRepo
	userRepo user.Repository
}

func NewUserUseCase(cfg *config.Config, txRepo *postgre.TxRepo, userRepo user.Repository) user.UseCase {
	return &userUC{cfg: cfg, txRepo: txRepo, userRepo: userRepo}
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
