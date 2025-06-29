package config

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	courseEntity "github.com/itmrchow/course-management-system/internal/domain/course/entity"
	teacherEntity "github.com/itmrchow/course-management-system/internal/domain/teacher/entity"
)

// NewPostgresDB 初始化 postgres db.
func NewPostgresDB(ctx context.Context, logger *zerolog.Logger) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC client_encoding=UTF8",
		viper.GetString("POSTGRES_HOST"),
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_DBNAME"),
		viper.GetString("POSTGRES_PORT"),
		viper.GetString("POSTGRES_SSLMODE"),
	)

	db, err := NewDB(ctx, postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Ctx(ctx).Msg("failed to init postgres db")
	}

	return db
}

func NewMemoryDB(ctx context.Context, logger *zerolog.Logger) *gorm.DB {
	viper.Set("POSTGRES_HOST", "localhost")
	viper.Set("POSTGRES_USER", "postgres")
	viper.Set("POSTGRES_PASSWORD", "postgres")
	viper.Set("POSTGRES_DBNAME", "postgres")
	viper.Set("POSTGRES_PORT", "5433")
	viper.Set("POSTGRES_SSLMODE", "disable")

	return NewPostgresDB(ctx, logger)
}

// NewDB 初始化 db.
func NewDB(ctx context.Context, dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {

	opts = append(opts, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	// teacher entities
	teacherEntities := []interface{}{
		&teacherEntity.Teacher{},
	}

	// course entities
	courseEntities := []interface{}{
		&courseEntity.Course{},
	}

	entities := append(teacherEntities, courseEntities...)

	// Auto migrate all entities
	err = db.AutoMigrate(
		entities...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return db, nil
}

// PingDB 呼叫 db.Ping() , 於初始化後呼叫.
func PingDB(ctx context.Context, logger *zerolog.Logger, db *gorm.DB) error {

	sqlDB, err := db.DB()
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("failed to get sql db")
		return err
	}
	logger.Info().Ctx(ctx).Msg("db initialized")
	err = sqlDB.Ping()
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("failed to ping db")
		return err
	}

	logger.Info().Ctx(ctx).Msg("db pinged")
	return nil
}
