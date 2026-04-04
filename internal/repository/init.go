package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	// "strconv"
	"time"

	"github.com/chenzanhong/zlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func Init() (*Repository, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = InitPGData(db, context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to init pg data: %w", err)
	}

	return &Repository{
		DB: db,
		// Redis: rdb,
	}, nil
}

func InitPGData(db *gorm.DB, ctx context.Context) error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	migrationsDir := "./migrations" // 相对于 WORKDIR=/app
	var err error
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(migrationsDir, file.Name())
			zlog.Infow("Execute: " + filePath)
			content, err := os.ReadFile(filePath)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("读取文件失败 %s: %w", filePath, err)
			}

			if err = tx.WithContext(ctx).Exec(string(content)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("执行 SQL 失败 %s: %w", filePath, err)
			}
		}
	}

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	zlog.Info("Migration completed successfully.")

	return nil
}

// 连接PostgreSQL
func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大生命周期

	return db, nil
}
