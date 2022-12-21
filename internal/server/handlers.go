package server

import (
	"fmt"
	"murakali/internal/middleware"
	authDelivery "murakali/internal/module/auth/delivery"
	authRepository "murakali/internal/module/auth/repository"
	authUseCase "murakali/internal/module/auth/usecase"
	cartDelivery "murakali/internal/module/cart/delivery"
	cartRepository "murakali/internal/module/cart/repository"
	cartUseCase "murakali/internal/module/cart/usecase"
	productDelivery "murakali/internal/module/product/delivery"
	productRepository "murakali/internal/module/product/repository"
	productUseCase "murakali/internal/module/product/usecase"
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

	s.gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://%s", s.cfg.Server.Origin)},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == fmt.Sprintf("http://%s", s.cfg.Server.Origin)
		},
		MaxAge: 12 * time.Hour,
	}))

	s.gin.NoRoute(func(c *gin.Context) {
		response.ErrorResponse(c.Writer, response.NotFoundMessage, http.StatusNotFound)
	})

	v1 := s.gin.Group("/api/v1")
	authGroup := v1.Group("/auth")
	userGroup := v1.Group("/user")
	productGroup := v1.Group("/product")
	cartGroup := v1.Group("/cart")

	mw := middleware.NewMiddlewareManager(s.cfg, []string{"*"}, s.log)
	authDelivery.MapAuthRoutes(authGroup, authHandlers)
	userDelivery.MapUserRoutes(userGroup, userHandlers, mw)
	productDelivery.MapProductRoutes(productGroup, productHandlers)
	cartDelivery.MapCartRoutes(cartGroup, cartHandlers, mw)

	return nil
}
