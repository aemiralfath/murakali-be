package delivery

import (
	"errors"
	"murakali/config"
	"murakali/internal/module/cart"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cartHandlers struct {
	cfg    *config.Config
	cartUC cart.UseCase
	logger logger.Logger
}

func NewCartHandlers(cfg *config.Config, cartUC cart.UseCase, log logger.Logger) cart.Handlers {
	return &cartHandlers{cfg: cfg, cartUC: cartUC, logger: log}
}

func (h *cartHandlers) GetCartHoverHome(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestParam body.CartHomeRequest
	if err := c.ShouldBind(&requestParam); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestParam.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	carts, err := h.cartUC.GetCartHoverHome(c, userID.(string), requestParam.Limit)
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

	response.SuccessResponse(c.Writer, carts, http.StatusCreated)
}
