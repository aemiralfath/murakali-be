package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"murakali/internal/middleware"
	authDelivery "murakali/internal/module/auth/delivery"
	authRepository "murakali/internal/module/auth/repository"
	authUseCase "murakali/internal/module/auth/usecase"
	userDelivery "murakali/internal/module/user/delivery"
	userRepository "murakali/internal/module/user/repository"
	userUseCase "murakali/internal/module/user/usecase"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"
)

func (s *Server) MapHandlers() error {
	txRepo := postgre.NewTxRepository(s.db)

	authRepo := authRepository.NewAuthRepository(s.db, s.redisClient)
	authUC := authUseCase.NewAuthUseCase(s.cfg, txRepo, authRepo)
	authHandlers := authDelivery.NewAuthHandlers(s.cfg, authUC, s.log)

	userRepo := userRepository.NewUserRepository(s.db, s.redisClient)
	userUC := userUseCase.NewUserUseCase(s.cfg, txRepo, userRepo)
	userHandlers := userDelivery.NewUserHandlers(s.cfg, userUC, s.log)

	s.gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3001"
		},
		MaxAge: 12 * time.Hour,
	}))

	s.gin.NoRoute(func(c *gin.Context) {
		response.ErrorResponse(c.Writer, response.NotFoundMessage, http.StatusNotFound)
	})

	v1 := s.gin.Group("/api/v1")
	authGroup := v1.Group("/auth")
	userGroup := v1.Group("/user")

	mw := middleware.NewMiddlewareManager(s.cfg, []string{"*"}, s.log)
	authDelivery.MapAuthRoutes(authGroup, authHandlers)
	userDelivery.MapUserRoutes(userGroup, userHandlers, mw)

	return nil
}
