package delivery

import (
	"errors"
	"github.com/gin-gonic/gin"
	"murakali/config"
	"murakali/internal/auth"
	"murakali/internal/auth/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/response"
	"net/http"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log}
}

func (h *authHandlers) RegisterEmail(c *gin.Context) {
	var requestBody body.RegisterEmailRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	_, err = h.authUC.RegisterEmail(c, requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerRegister, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusCreated)
}

func (h *authHandlers) RegisterUser(c *gin.Context) {
	var requestBody body.RegisterUserRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.authUC.RegisterUser(c, requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerRegister, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *authHandlers) VerifyOTP(c *gin.Context) {
	var requestBody body.VerifyOTPRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.authUC.VerifyOTP(c, requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerRegister, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}
