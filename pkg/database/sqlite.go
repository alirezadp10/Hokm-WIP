package database

import (
	"github.com/alirezadp10/hokm/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func GetNewSqliteConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Player{}, &model.Game{}); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	return db
}
