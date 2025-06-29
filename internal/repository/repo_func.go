package repository

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// RepoPageInfo repo分頁查詢資訊
type RepoPageInfo struct {
	Page     int    `json:"page" default:"1"`
	PageSize int    `json:"page_size" default:"10"`
	Sort     string `json:"sort" default:"id"`
	Order    string `json:"order" default:"desc"`
}

// RepoTestInit 初始化 repository 測試環境
// 使用 embedded-postgres 啟動 postgres 測試環境
// 測試完畢後停止 postgres
func RepoTestInit(m *testing.M) int {
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(5433))

	err := postgres.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start postgres")
	}

	// 註冊 signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info().Msg("Received interrupt, stopping postgres...")
		if err := postgres.Stop(); err != nil {
			log.Error().Err(err).Msg("failed to stop postgres on interrupt")
		}
		os.Exit(1)
	}()

	code := m.Run()

	err = postgres.Stop()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to stop postgres")
	}

	return code
}

func Paginate(pageInfo *RepoPageInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageInfo.Page <= 0 {
			pageInfo.Page = 1
		}

		offset := (pageInfo.Page - 1) * pageInfo.PageSize
		return db.Offset(offset).Limit(pageInfo.PageSize)
	}
}
