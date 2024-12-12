package database

import (
    "gopkg.in/telebot.v4"
    "gorm.io/gorm"
    "time"
)

func SavePlayer(db *gorm.DB, player *telebot.User, chatId int64) (*Player, error) {
    newPlayer := Player{
        Id:        player.ID,
        FirstName: player.FirstName,
        LastName:  player.LastName,
        Username:  player.Username,
        ChatId:    chatId,
        UpdatedAt: time.Now(),
        JoinedAt:  time.Now(),
    }
    err := db.Save(&newPlayer).Error
    return &newPlayer, err
}
