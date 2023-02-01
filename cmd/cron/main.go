package main

import (
	"fmt"
	"log"
	"murakali/config"
	"murakali/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)

func main() {
	log.Println("Starting cron server")
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

	cronJob := cron.New()
	defer cronJob.Stop()

	_, err = cronJob.AddFunc("@every 1m", func() {
		updateOnDelivery(cfg, appLogger)
	})
	if err != nil {
		appLogger.Warn("FatalConfig: %v", err)
	}

	_, err = cronJob.AddFunc("@every 1m", func() {
		updateExpiredAt(cfg, appLogger)
	})
	if err != nil {
		appLogger.Warn("FatalConfig: %v", err)
	}

	_, err = cronJob.AddFunc("@every 1m", func() {
		updateRejectedRefund(cfg, appLogger)
	})
	if err != nil {
		appLogger.Warn("FatalConfig: %v", err)
	}

	// product -> fav count, rating_avg
	// shop -> total_product, total_rating, rating avg

	go cronJob.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	appLogger.Info("Cron stop")
}

func updateOnDelivery(cfg *config.Config, appLogger logger.Logger) {
	appLogger.Info("cron update delivery start")
	url := fmt.Sprintf("https://%s/api/v1/seller/delivery", cfg.Server.Domain)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		appLogger.Warnf("request error: ", err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		appLogger.Warn("response error: ", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		appLogger.Warn("status code error: ", res.StatusCode)
		return
	}

	appLogger.Infof("update delivery success")
}

func updateExpiredAt(cfg *config.Config, appLogger logger.Logger) {
	appLogger.Info("cron update expired at start")
	url := fmt.Sprintf("https://%s/api/v1/seller/expired", cfg.Server.Domain)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		appLogger.Warnf("request error: ", err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		appLogger.Warn("response error: ", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		appLogger.Warn("status code error: ", res.StatusCode)
		return
	}

	appLogger.Infof("update expired success")
}

func updateRejectedRefund(cfg *config.Config, appLogger logger.Logger) {
	appLogger.Info("cron update rejected at start")
	url := fmt.Sprintf("https://%s/api/v1/user/rejected-refund", cfg.Server.Domain)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		appLogger.Warnf("request error: ", err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		appLogger.Warn("response error: ", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		appLogger.Warn("status code error: ", res.StatusCode)
		return
	}

	appLogger.Infof("update rejected success")
}
