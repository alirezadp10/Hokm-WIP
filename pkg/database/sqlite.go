package database

import (
    "github.com/alirezadp10/hokm/pkg/model"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
    "os"
)

func GetNewSqliteConnection() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(os.Getenv("DB_NAME")), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }

    err = db.AutoMigrate(&model.Player{}, &model.Game{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    return db
}
