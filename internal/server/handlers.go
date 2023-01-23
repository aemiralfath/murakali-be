package server

import (
	"murakali/internal/middleware"
	adminDelivery "murakali/internal/module/admin/delivery"
	adminRepository "murakali/internal/module/admin/repository"
	adminUseCase "murakali/internal/module/admin/usecase"
	authDelivery "murakali/internal/module/auth/delivery"
	authRepository "murakali/internal/module/auth/repository"
	authUseCase "murakali/internal/module/auth/usecase"
	cartDelivery "murakali/internal/module/cart/delivery"
	cartRepository "murakali/internal/module/cart/repository"
	cartUseCase "murakali/internal/module/cart/usecase"
	locationDelivery "murakali/internal/module/location/delivery"
	locationRepository "murakali/internal/module/location/repository"
	locationUseCase "murakali/internal/module/location/usecase"
	productDelivery "murakali/internal/module/product/delivery"
	productRepository "murakali/internal/module/product/repository"
	productUseCase "murakali/internal/module/product/usecase"
	sellerDelivery "murakali/internal/module/seller/delivery"
	sellerRepository "murakali/internal/module/seller/repository"
	sellerUseCase "murakali/internal/module/seller/usecase"
	userDelivery "murakali/internal/module/user/delivery"
	userRepository "murakali/internal/module/user/repository"
	userUseCase "murakali/internal/module/user/usecase"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	txRepo := postgre.NewTxRepository(s.db)

	adminRepo := adminRepository.NewAdminRepository(s.db, s.redisClient)
	adminUC := adminUseCase.NewAdminUseCase(s.cfg, txRepo, adminRepo)
	adminHandlers := adminDelivery.NewAdminHandlers(s.cfg, adminUC, s.log)

	authRepo := authRepository.NewAuthRepository(s.db, s.redisClient)
	authUC := authUseCase.NewAuthUseCase(s.cfg, txRepo, authRepo)
	authHandlers := authDelivery.NewAuthHandlers(s.cfg, authUC, s.log)

	userRepo := userRepository.NewUserRepository(s.db, s.redisClient)
	userUC := userUseCase.NewUserUseCase(s.cfg, txRepo, userRepo)
	userHandlers := userDelivery.NewUserHandlers(s.cfg, userUC, s.log)

	productRepo := productRepository.NewProductRepository(s.db, s.redisClient)
	productUC := productUseCase.NewProductUseCase(s.cfg, txRepo, productRepo)
	productHandlers := productDelivery.NewProductHandlers(s.cfg, productUC, s.log)

	cartRepo := cartRepository.NewCartRepository(s.db, s.redisClient)
	cartUC := cartUseCase.NewCartUseCase(s.cfg, txRepo, cartRepo)
	cartHandlers := cartDelivery.NewCartHandlers(s.cfg, cartUC, s.log)

	locationRepo := locationRepository.NewLocationRepository(s.db, s.redisClient)
	locationUC := locationUseCase.NewLocationUseCase(s.cfg, txRepo, locationRepo)
	locationHandlers := locationDelivery.NewLocationHandlers(s.cfg, locationUC, s.log)

	sellerRepo := sellerRepository.NewSellerRepository(s.db, s.redisClient)
	sellerUC := sellerUseCase.NewSellerUseCase(s.cfg, txRepo, sellerRepo)
	sellerHandlers := sellerDelivery.NewSellerHandlers(s.cfg, sellerUC, s.log)

	s.gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{s.cfg.Server.Origin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == s.cfg.Server.Origin
		},
		MaxAge: 12 * time.Hour,
	}))

	s.gin.Static("/docs", "dist/")
	s.gin.NoRoute(func(c *gin.Context) {
		response.ErrorResponse(c.Writer, response.NotFoundMessage, http.StatusNotFound)
	})

	v1 := s.gin.Group("/api/v1")
	authGroup := v1.Group("/auth")
	userGroup := v1.Group("/user")
	productGroup := v1.Group("/product")
	cartGroup := v1.Group("/cart")
	locationGroup := v1.Group("/location")
	sellerGroup := v1.Group("/seller")
	adminGroup := v1.Group("/admin")

	mw := middleware.NewMiddlewareManager(s.cfg, []string{"*"}, s.log)
	authDelivery.MapAuthRoutes(authGroup, authHandlers)
	userDelivery.MapUserRoutes(userGroup, userHandlers, mw)
	productDelivery.MapProductRoutes(productGroup, productHandlers, mw)
	cartDelivery.MapCartRoutes(cartGroup, cartHandlers, mw)
	locationDelivery.MapAuthRoutes(locationGroup, locationHandlers)
	sellerDelivery.MapSellerRoutes(sellerGroup, sellerHandlers, mw)
	adminDelivery.MapAdminRoutes(adminGroup, adminHandlers, mw)

	return nil
}
