package delivery

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"murakali/config"
	"murakali/internal/module/user"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

type userHandlers struct {
	cfg    *config.Config
	userUC user.UseCase
	logger logger.Logger
}

func NewUserHandlers(cfg *config.Config, userUC user.UseCase, log logger.Logger) user.Handlers {
	return &userHandlers{cfg: cfg, userUC: userUC, logger: log}
}

func (h *userHandlers) GetAddress(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	pgn := &pagination.Pagination{}
	name := h.ValidateQueryAddress(c, pgn)

	addresses, err := h.userUC.GetAddress(c, userID.(string), name, pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, addresses, http.StatusOK)
}

func (h *userHandlers) ValidateQueryAddress(c *gin.Context, pgn *pagination.Pagination) string {
	name := strings.TrimSpace(c.Query("name"))
	sort := strings.TrimSpace(c.Query("sort"))
	sortBy := strings.TrimSpace(c.Query("sortBy"))
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))

	var sortFilter string
	var sortByFilter string
	var limitFilter int
	var pageFilter int

	switch sort {
	case "asc":
		sortFilter = sort
	default:
		sortFilter = "desc"
	}

	switch sortBy {
	case "province":
		sortByFilter = sortBy
	default:
		sortByFilter = "created_at"
	}

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 10
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}

	pgn.Limit = limitFilter
	pgn.Page = pageFilter
	pgn.Sort = fmt.Sprintf("%s %s", sortByFilter, sortFilter)

	return name
}
