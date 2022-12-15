package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	authDelivery "murakali/internal/auth/delivery"
	authRepository "murakali/internal/auth/repository"
	authUseCase "murakali/internal/auth/usecase"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"
)

func (s *Server) MapHandlers() error {
	txRepo := postgre.NewTxRepository(s.db)
	aRepo, err := authRepository.NewAuthRepository(s.db, s.redisClient)
	if err != nil {
		return err
	}

	authUC := authUseCase.NewAuthUseCase(s.cfg, txRepo, aRepo)
	authHandlers := authDelivery.NewAuthHandlers(s.cfg, authUC, s.logger)

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

	authDelivery.MapAuthRoutes(authGroup, authHandlers)

	return nil
}
