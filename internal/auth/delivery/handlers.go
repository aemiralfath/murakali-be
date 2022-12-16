package delivery

import (
	"errors"
	"murakali/config"
	"murakali/internal/auth"
	"murakali/internal/auth/delivery/body"
	"murakali/internal/constant"
	"murakali/pkg/httperror"
	"murakali/pkg/jwt"
	"murakali/pkg/logger"
	"murakali/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log}
}

func (h *authHandlers) Login(c *gin.Context) {
	var requestBody body.LoginRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	accessToken, refreshToken, err := h.authUC.Login(c, requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAuth, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetCookie(constant.AccessTokenCookie, accessToken, h.cfg.JWT.AccessExpMin*60, "/", h.cfg.Server.Domain, false, true)
	c.SetCookie(constant.RefreshTokenCookie, refreshToken, h.cfg.JWT.RefreshExpMin*60, "/", h.cfg.Server.Domain, false, true)
	response.SuccessResponse(c.Writer, body.LoginResponse{AccessToken: accessToken}, http.StatusOK)
}

func (h *authHandlers) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(constant.RefreshTokenCookie)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	claims, err := jwt.ExtractJWT(refreshToken, h.cfg.JWT.JwtSecretKey)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	accessToken, err := h.authUC.RefreshToken(c, claims["id"].(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAuth, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetCookie(constant.AccessTokenCookie, accessToken, h.cfg.JWT.AccessExpMin*60, "/", h.cfg.Server.Domain, false, true)
	response.SuccessResponse(c.Writer, body.LoginResponse{AccessToken: accessToken}, http.StatusOK)
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
			h.logger.Errorf("HandlerAuth, Error: %s", err)
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
			h.logger.Errorf("HandlerAuth, Error: %s", err)
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
			h.logger.Errorf("HandlerAuth, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *authHandlers) ResetPasswordEmail(c *gin.Context) {
	var requestBody body.ResetPasswordEmailRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	_, err = h.authUC.ResetPasswordEmail(c, requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAuth, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusCreated)
}
