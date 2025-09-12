package repository

import (
	"backend/logs"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

func InitDB() error {
	if err := ConnectDB(); err != nil {
		return fmt.Errorf("ConnectDB failed: %w", err)
	}

	if err := ConnectRedis(); err != nil {
		return fmt.Errorf("ConnectRedis failed: %w", err)
	}

	if err := InitData(); err != nil {
		return fmt.Errorf("InitData failed: %w", err)
	}

	return nil
}

// 连接PostgreSQL
func ConnectDB() error {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	var err error
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大生命周期

	return nil
}

// 连接Redis
func ConnectRedis() error {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db := 0
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			return fmt.Errorf("invalid REDIS_DB value: %v", err)
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	logs.Sugar.Info("Redis connected successfully")
	return nil
}

func InitData() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	_, filename, _, ok := runtime.Caller(0) // 获取当前的文件名
	if !ok {
		logs.Sugar.Fatal("无法获取运行时调用者信息")
	}

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(filename)
	migrationsDir := filepath.Join(currentDir, "../../migrations")
	var err error
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	tx := DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
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
			fmt.Println("Execute: ", filePath)
			content, err := os.ReadFile(filePath)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("读取文件失败 %s: %v", filePath, err)
			}

			if err = tx.Exec(string(content)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("执行 SQL 失败 %s: %v", filePath, err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	fmt.Println("Migration completed successfully.")
	return nil
}
