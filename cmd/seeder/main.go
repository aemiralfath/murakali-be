package main

import (
	"log"
	"murakali/config"
	"murakali/pkg/logger"
	"murakali/pkg/postgre"
	"murakali/sql/fakers"
)

func main() {
	log.Println("Starting seeder server")
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

	pgDB, err := postgre.NewPG(cfg, appLogger)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	appLogger.Infof("Postgres connected")

	if err := fakers.DBSeed(pgDB); err != nil {
		appLogger.Fatalf("seeder init: %s", err)
	}
	appLogger.Infof("Seeder done")
}
