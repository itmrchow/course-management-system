package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/itmrchow/course-management-system/internal/config"
)

func main() {

	// 系統信號處理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ctx
	ctx, cancel := context.WithCancel(context.Background())

	// config
	config.InitConfig()

	// logger
	logger := config.InitLogger()

	// db
	db := config.NewPostgresDB(context.Background(), logger)
	if err := config.PingDB(ctx, logger, db); err != nil {
		logger.Err(err).Msg("failed to ping db")
	}

	select {
	case sig := <-sigChan:
		logger.Info().Msgf("收到系統信號: %v, 開始關閉服務", sig)
		cancel()
	case <-ctx.Done():
		logger.Info().Msg("服務開始關閉")
		// close db
		sqlDB, err := db.DB()
		if err != nil {
			logger.Err(err).Msg("failed to get sql db")
		}
		if err := sqlDB.Close(); err != nil {
			logger.Err(err).Msg("failed to close sql db")
		}
	}
}
