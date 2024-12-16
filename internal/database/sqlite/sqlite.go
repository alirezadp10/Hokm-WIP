package sqlite

import (
    "github.com/alirezadp10/hokm/internal/database"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
    "os"
)

func GetNewConnection() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(os.Getenv("DB_NAME")), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }

    err = db.AutoMigrate(&database.Player{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    return db
}
