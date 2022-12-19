package main

import (
	"log"
	"murakali/config"
	"murakali/internal/server"
	"murakali/pkg/logger"
	"murakali/pkg/postgre"
	"murakali/pkg/redis"
)

func main() {
	log.Println("Starting api server")
	cfgFile, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	pgDB, err := postgre.NewPG(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	appLogger.Infof("Postgres connected")

	redisClient, err := redis.NewRedis(cfg)
	if err != nil {
		appLogger.Fatalf("redis init: %s", err)
	}

	appLogger.Infof("Redis connected")

	s := server.NewServer(cfg, pgDB, redisClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
