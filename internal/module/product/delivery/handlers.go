package delivery

import (
	"errors"
	"fmt"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *productHandlers) CreateFavoriteProduct(c *gin.Context) {
	var requestBody body.GetProductRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	err = h.productUC.CreateFavoriteProduct(c, requestBody.ProductID, userID.(string))
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

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *productHandlers) DeleteFavoriteProduct(c *gin.Context) {
	var requestBody body.GetProductRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	err = h.productUC.DeleteFavoriteProduct(c, requestBody.ProductID, userID.(string))
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

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
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

func (h *productHandlers) GetAllProductImage(c *gin.Context) {
	productID := c.Param("product_id")
	productImages, err := h.productUC.GetAllProductImage(c, productID)
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

	response.SuccessResponse(c.Writer, productImages, http.StatusOK)
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

	sortBy = strings.ToLower(sortBy)
	sort = strings.ToLower(sort)

	minPrice := strings.TrimSpace(c.Query("min_price"))
	maxPrice := strings.TrimSpace(c.Query("max_price"))
	minRating := strings.TrimSpace(c.Query("min_rating"))
	maxRating := strings.TrimSpace(c.Query("max_rating"))

	category := strings.TrimSpace(c.Query("category"))

	category = strings.ToLower(category)

	shop := strings.TrimSpace(c.Query("shop_id"))

	listedStatus := strings.TrimSpace(c.Query("listed_status"))

	province := strings.TrimSpace(c.Query("province_ids"))

	var sortFilter string
	var limitFilter, pageFilter int
	var minPriceFilter, maxPriceFilter, minRatingFilter, maxRatingFilter float64

	listedStatusFilter, _ := strconv.Atoi(listedStatus)
	switch listedStatusFilter {
	case 0:
		listedStatusFilter = 0
	case 1:
		listedStatusFilter = 1
	case 2:
		listedStatusFilter = 2
	default:
		listedStatusFilter = 0
	}

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 12
	} else if limitFilter > 100 {
		limitFilter = 100
	}

	switch sort {
	case constant.ASC:
		sortFilter = sort
	default:
		sortFilter = constant.DESC
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}

	var pgn *pagination.Pagination
	switch sortBy {
	case "created_at":
		pgn = &pagination.Pagination{
			Limit: limitFilter,
			Page:  pageFilter,
			Sort:  "p." + sortBy + " " + sortFilter,
		}
	case "recommended":
		pgn = &pagination.Pagination{
			Limit: limitFilter,
			Page:  pageFilter,
			Sort:  "view_count " + sortFilter + ", " + "unit_sold " + sortFilter,
		}
	case "min_price":
		pgn = &pagination.Pagination{
			Limit: limitFilter,
			Page:  pageFilter,
			Sort:  "min_price" + " " + sortFilter,
		}
	case "unit_sold":
		pgn = &pagination.Pagination{
			Limit: limitFilter,
			Page:  pageFilter,
			Sort:  "unit_sold" + " " + sortFilter,
		}
	case "view_count":
		pgn = &pagination.Pagination{
			Limit: limitFilter,
			Page:  pageFilter,
			Sort:  "view_count" + " " + sortFilter,
		}
	default:
		pgn = &pagination.Pagination{
			Limit: limitFilter,
			Page:  pageFilter,
			Sort:  "p.created_at" + " " + sortFilter,
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
	} else if minRatingFilter > 5 {
		minRatingFilter = 0
	}

	maxRatingFilter, err = strconv.ParseFloat(maxRating, 64)
	if err != nil || maxRatingFilter > 5 || maxRatingFilter <= 0 {
		maxRatingFilter = 5
	}

	searchFilter := fmt.Sprintf("%%%s%%", search)
	categoryFilter := fmt.Sprintf("%%%s%%", category)

	var provinceFilter []string
	if province != "" {
		provinceFilter = strings.Split(province, ",")
	}

	if maxRatingFilter < minRatingFilter {
		minRatingFilter = 0
		maxRatingFilter = 5
	}
	query := &body.GetProductQueryRequest{
		Search:       searchFilter,
		Shop:         shop,
		Category:     categoryFilter,
		MinPrice:     minPriceFilter,
		MaxPrice:     maxPriceFilter,
		MinRating:    minRatingFilter,
		MaxRating:    maxRatingFilter,
		Province:     provinceFilter,
		ListedStatus: listedStatusFilter,
	}
	return pgn, query
}

func (h *productHandlers) CreateProduct(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.CreateProductRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.ValidateCreateProduct()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.productUC.CreateProduct(c, requestBody, userID.(string)); err != nil {
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

func (h *productHandlers) UpdateListedStatus(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.productUC.UpdateListedStatus(c, productID.String()); err != nil {
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

func (h *productHandlers) UpdateListedStatusBulk(c *gin.Context) {
	var requestBody body.UpdateProductListedStatusBulkRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.ValidateUpdateProductListedStatusBulk()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.productUC.UpdateProductListedStatusBulk(c, requestBody); err != nil {
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

func (h *productHandlers) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.UpdateProductRequest
	if err = c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.ValidateUpdateProduct()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.productUC.UpdateProduct(c, requestBody, userID.(string), productID.String()); err != nil {
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

func (h *productHandlers) UploadProductPicture(c *gin.Context) {
	type Sizer interface {
		Size() int64
	}

	var imgURL string
	var img body.ImageRequest

	err := c.ShouldBind(&img)
	if err != nil {
		response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	data, _, err := c.Request.FormFile("Img")
	if err != nil {
		response.ErrorResponse(c.Writer, body.ImageIsEmpty, http.StatusInternalServerError)
		return
	}

	if data.(Sizer).Size() > constant.ImgMaxSize {
		response.ErrorResponse(c.Writer, response.PictureSizeTooBig, http.StatusInternalServerError)
		return
	}

	if data == nil {
		response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	imgURL = util.UploadImageToCloudinary(c, h.cfg, data)

	response.SuccessResponse(c.Writer, imgURL, http.StatusOK)
}

// get product review by product id
func (h *productHandlers) GetProductReviews(c *gin.Context) {
	productID := c.Param("product_id")
	pgn, query := h.ValidateQueryReview(c)

	reviews, err := h.productUC.GetProductReviews(c, pgn, productID, query)
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

	response.SuccessResponse(c.Writer, reviews, http.StatusOK)
}

func (h *productHandlers) CreateProductReview(c *gin.Context) {
	var requestBody body.ReviewProductRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	err = h.productUC.CreateProductReview(c, requestBody, userID.(string))
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

	response.SuccessResponse(c.Writer, nil, http.StatusCreated)
}

func (h *productHandlers) DeleteProductReview(c *gin.Context) {
	reviewId := c.Param("review_id")

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	err := h.productUC.DeleteProductReview(c, reviewId, userID.(string))
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

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *productHandlers) ValidateQueryReview(c *gin.Context) (*pagination.Pagination, *body.GetReviewQueryRequest) {
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))
	sort := strings.TrimSpace(c.Query("sort"))

	var limitFilter int
	var pageFilter int
	var sortFilter string

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 5
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}

	switch sort {
	case "asc":
		sortFilter = sort
	default:
		sortFilter = "desc"
	}

	pgn := &pagination.Pagination{
		Limit: limitFilter,
		Page:  pageFilter,
		Sort:  "r.created_at " + sortFilter,
	}

	rating := strings.TrimSpace(c.Query("rating"))
	showComment := strings.TrimSpace(c.Query("show_comment"))
	showImage := strings.TrimSpace(c.Query("show_image"))
	userId := strings.TrimSpace(c.Query("user_id"))

	var ratingFilter int
	ratingFilterInput := "0"
	var showCommentFilter bool
	var showImageFilter bool
	userIdFilter := userId

	ratingFilter, err = strconv.Atoi(rating)

	if err == nil && ratingFilter >= 1 && ratingFilter <= 5 {
		ratingFilterInput = strconv.Itoa(ratingFilter)
	}

	if showComment == constant.FALSE {
		showCommentFilter = false
	} else {
		showCommentFilter = true
	}

	if showImage == constant.FALSE {
		showImageFilter = false
	} else {
		showImageFilter = true
	}

	if _, err := uuid.Parse(userId); err != nil {
		userIdFilter = ""
	}

	query := &body.GetReviewQueryRequest{
		Rating:      ratingFilterInput,
		ShowComment: showCommentFilter,
		ShowImage:   showImageFilter,
		UserID:      userIdFilter,
	}

	return pgn, query
}

func (h *productHandlers) GetTotalReviewRatingByProductID(c *gin.Context) {
	productID := c.Param("product_id")

	reviewRating, err := h.productUC.GetTotalReviewRatingByProductID(c, productID)
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

	response.SuccessResponse(c.Writer, reviewRating, http.StatusOK)
}
