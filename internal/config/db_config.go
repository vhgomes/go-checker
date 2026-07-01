package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-checker/internal/repository"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dbName := GetEnvOrDefault("DB_NAME", "test.db")
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("failed to connect database", zap.Error(err))
	}

	err = db.AutoMigrate(&repository.Site{}, &repository.User{}, &repository.SiteStatusHistory{})
	if err != nil {
		zap.L().Fatal("Erro ao conectar no banco", zap.Error(err))
	}

	return db
}

func InitRedis() *redis.Client {
	redisAddr := GetEnvOrDefault("REDIS_ADDR", "localhost:6379")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	fmt.Println("Connected to Redis!")
	return client
}
