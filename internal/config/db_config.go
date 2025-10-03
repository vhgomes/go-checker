package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-checker/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&repository.Site{}, &repository.User{}, &repository.SiteStatusHistory{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	fmt.Println("Connected to Redis!")
	return client
}
