package server

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"murakali/config"
	"murakali/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	gin         *gin.Engine
	cfg         *config.Config
	db          *sql.DB
	redisClient *redis.Client
	log         logger.Logger
}

func NewServer(cfg *config.Config, db *sql.DB, redisClient *redis.Client, log logger.Logger) *Server {
	return &Server{gin: gin.Default(), cfg: cfg, db: db, redisClient: redisClient, log: log}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		Handler:        s.gin,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Error starting Server: ", err)
		}
	}()

	if err := s.MapHandlers(); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	s.log.Info("Shutdown Server ...")

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	if err := server.Shutdown(ctx); err != nil {
		s.log.Fatal("Server Shutdown:", err)
		return err
	}

	<-ctx.Done()
	s.log.Infof("timeout.")
	s.log.Info("Server Exited Properly")
	return nil
}
