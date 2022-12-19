package delivery

import (
	"murakali/config"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/logger"
	"murakali/pkg/response"
	"net/http"

	"murakali/pkg/httperror"

	"errors"

	"github.com/gin-gonic/gin"
)

type userHandlers struct {
	cfg    *config.Config
	userUC user.UseCase
	logger logger.Logger
}

func NewUserHandlers(cfg *config.Config, userUC user.UseCase, log logger.Logger) user.Handlers {
	return &userHandlers{cfg: cfg, userUC: userUC, logger: log}
}

func (h *userHandlers) GetSealabsPay(c *gin.Context) {
	userid, exist := c.Get("userID")

	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	result, err := h.userUC.GetSealabsPay(c, userid.(string))
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
	response.SuccessResponse(c.Writer, result, http.StatusOK)

}

func (h *userHandlers) AddSealabsPay(c *gin.Context) {
	userid, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.AddSealabsPayRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.AddSealabsPay(c, requestBody, userid.(string)); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) PatchSealabsPay(c *gin.Context) {
	userid, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}
	var requestBody body.PatchSealabsPayRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}
	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.PatchSealabsPay(c, requestBody.CardNumber, userid.(string)); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}
