package delivery

import (
	"errors"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/module/auth"
	"murakali/internal/module/auth/delivery/body"
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

func (h *authHandlers) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.RefreshTokenCookie, "", -1, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
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

	token, err := h.authUC.Login(c, requestBody)
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

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.RefreshTokenCookie, token.RefreshToken, h.cfg.JWT.RefreshExpMin*60, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, body.LoginResponse{AccessToken: token.AccessToken}, http.StatusOK)
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
	registerToken, err := c.Cookie(constant.RegisterTokenCookie)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	claims, err := jwt.ExtractJWT(registerToken, h.cfg.JWT.JwtSecretKey)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	var requestBody body.RegisterUserRequest
	if c.ShouldBind(&requestBody) != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.authUC.RegisterUser(c, claims["email"].(string), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAuth, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.RegisterTokenCookie, "", -1, "/", h.cfg.Server.Domain, true, true)
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

	registerToken, err := h.authUC.VerifyOTP(c, requestBody)
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

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.RegisterTokenCookie, registerToken, h.cfg.JWT.RefreshExpMin*60, "/", h.cfg.Server.Domain, true, true)
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

func (h *authHandlers) ResetPasswordUser(c *gin.Context) {
	ResetPasswordToken, err := c.Cookie(constant.ResetPasswordTokenCookie)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	claims, err := jwt.ExtractJWT(ResetPasswordToken, h.cfg.JWT.JwtSecretKey)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	var requestBody body.ResetPasswordUserRequest
	if c.ShouldBind(&requestBody) != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	_, err = h.authUC.ResetPasswordUser(c, claims["email"].(string), &requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAuth, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		invalidFields, err = requestBody.Validate()
		if err != nil {
			response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.ResetPasswordTokenCookie, "", -1, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, nil, http.StatusCreated)
}

func (h *authHandlers) ResetPasswordVerifyOTP(c *gin.Context) {
	var requestBody body.ResetPasswordVerifyOTPRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	ResetPasswordToken, err := h.authUC.ResetPasswordVerifyOTP(c, requestBody)
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

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.ResetPasswordTokenCookie, ResetPasswordToken, h.cfg.JWT.RefreshExpMin*60, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}
