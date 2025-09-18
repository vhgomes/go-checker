package config

import (
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

	err = db.AutoMigrate(&repository.Site{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&repository.SiteStatusHistory{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
