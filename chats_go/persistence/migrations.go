package persistence

import (
	"chats_go/domain"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.Chat{},
		&domain.Message{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database!")
	}
}
