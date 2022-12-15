package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"murakali/pkg/response"
	"net/http"
	"time"
)

func (s *Server) MapHandlers() error {

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
	authGroup.GET("/", func(c *gin.Context) {
		response.SuccessResponse(c.Writer, "success", 200)
	})

	return nil
}
