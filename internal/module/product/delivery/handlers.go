package delivery

import (
	"errors"
	"fmt"
	"murakali/config"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type productHandlers struct {
	cfg       *config.Config
	productUC product.UseCase
	logger    logger.Logger
}

func NewProductHandlers(cfg *config.Config, productUC product.UseCase, log logger.Logger) product.Handlers {
	return &productHandlers{cfg: cfg, productUC: productUC, logger: log}
}

func (h *productHandlers) GetCategories(c *gin.Context) {
	categoriesResponse, err := h.productUC.GetCategories(c)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetBanners(c *gin.Context) {
	banners, err := h.productUC.GetBanners(c)
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

	response.SuccessResponse(c.Writer, banners, http.StatusOK)
}

func (h *productHandlers) GetCategoriesByNameLevelOne(c *gin.Context) {
	var requestPath body.CategoryRequest
	requestPath.NameLevelOne = c.Param("name_lvl_one")

	categoriesResponse, err := h.productUC.GetCategoriesByName(c, requestPath.NameLevelOne)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetCategoriesByNameLevelTwo(c *gin.Context) {
	var requestPath body.CategoryRequest
	requestPath.NameLevelTwo = c.Param("name_lvl_two")

	categoriesResponse, err := h.productUC.GetCategoriesByName(c, requestPath.NameLevelTwo)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetCategoriesByNameLevelThree(c *gin.Context) {
	var requestPath body.CategoryRequest
	requestPath.NameLevelThree = c.Param("name_lvl_three")

	categoriesResponse, err := h.productUC.GetCategoriesByName(c, requestPath.NameLevelThree)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetRecommendedProducts(c *gin.Context) {
	pgn := h.ValidateQueryRecommendProduct(c)
	RecommendedProducts, err := h.productUC.GetRecommendedProducts(c, pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, RecommendedProducts, http.StatusOK)
}

func (h *productHandlers) GetProducts(c *gin.Context) {
	pgn, query := h.ValidateQueryProduct(c)
	SearchProducts, err := h.productUC.GetProducts(c, pgn, query)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, SearchProducts, http.StatusOK)
}

func (h *productHandlers) GetFavoriteProducts(c *gin.Context) {
	pgn, query := h.ValidateQueryProduct(c)

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	SearchProducts, err := h.productUC.GetFavoriteProducts(c, pgn, query, userID.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, SearchProducts, http.StatusOK)
}

func (h *productHandlers) GetProductDetail(c *gin.Context) {
	productID := c.Param("product_id")
	productDetail, err := h.productUC.GetProductDetail(c, productID)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, productDetail, http.StatusOK)
}

func (h *productHandlers) ValidateQueryRecommendProduct(c *gin.Context) *pagination.Pagination {
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))

	var limitFilter int
	var pageFilter int

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 18
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}
	pgn := &pagination.Pagination{
		Limit: limitFilter,
		Page:  pageFilter,
		Sort:  "unit_sold DESC",
	}

	return pgn
}

func (h *productHandlers) ValidateQueryProduct(c *gin.Context) (*pagination.Pagination, *body.GetProductQueryRequest) {
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))

	search := strings.TrimSpace(c.Query("search"))

	sort := strings.TrimSpace(c.Query("sort"))
	sortBy := strings.TrimSpace(c.Query("sort_by"))

	minPrice := strings.TrimSpace(c.Query("min_price"))
	maxPrice := strings.TrimSpace(c.Query("max_price"))
	minRating := strings.TrimSpace(c.Query("min_rating"))
	maxRating := strings.TrimSpace(c.Query("max_rating"))

	category := strings.TrimSpace(c.Query("category"))
	shop := strings.TrimSpace(c.Query("shop_name"))

	province := strings.TrimSpace(c.Query("province"))

	var limitFilter, pageFilter int
	var minPriceFilter, maxPriceFilter, minRatingFilter, maxRatingFilter float64

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 12
	}

	if sortBy == "" {
		sortBy = "p.unit_sold"
	}
	if sort == "" {
		sort = "desc"
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}

	pgn := &pagination.Pagination{
		Limit: limitFilter,
		Page:  pageFilter,
		Sort:  sortBy + " " + sort,
	}

	if sortBy == "recommended" {
		pgn = &pagination.Pagination{
			Limit: limitFilter,
			Page:  pageFilter,
			Sort:  "view_count " + sort + ", " + "unit_sold " + sort,
		}
	}

	minPriceFilter, err = strconv.ParseFloat(minPrice, 64)
	if err != nil || minPriceFilter <= 0 {
		minPriceFilter = 0
	}

	maxPriceFilter, err = strconv.ParseFloat(maxPrice, 64)
	if err != nil || maxPriceFilter == 0 {
		maxPriceFilter = 999999999999
	}

	minRatingFilter, err = strconv.ParseFloat(minRating, 64)
	if err != nil || minRatingFilter <= 0 {
		minRatingFilter = 0
	}

	maxRatingFilter, err = strconv.ParseFloat(maxRating, 64)
	if err != nil || maxRatingFilter > 5 || maxRatingFilter <= 0 {
		maxRatingFilter = 5
	}

	searchFilter := fmt.Sprintf("%%%s%%", search)
	categoryFilter := fmt.Sprintf("%%%s%%", category)
	shopFilter := fmt.Sprintf("%%%s%%", shop)

	// provinceFilter := strings.SplitAfter(province, ",")
	// provinceIntFilter := make([]int, 0)
	// for _, i := range provinceFilter {
	// 	 temp := strings.TrimRight(i, ",")
	//     j, err := strconv.Atoi(temp)
	//     if err != nil {
	//         fmt.Println(err)
	//     }
	//     provinceIntFilter = append(provinceIntFilter, j)
	// }

	query := &body.GetProductQueryRequest{
		Search: searchFilter,

		Shop:      shopFilter,
		Category:  categoryFilter,
		MinPrice:  minPriceFilter,
		MaxPrice:  maxPriceFilter,
		MinRating: minRatingFilter,
		MaxRating: maxRatingFilter,
		Province:  province,
	}

	return pgn, query
}
