package database

import (
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

    return db
}
