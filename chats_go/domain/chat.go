package domain

import (
	"time"
)

type Chat struct {
	ID               uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt        time.Time
	ApplicationToken string    `gorm:"size:36;index:idx_application_token_number,unique"`
	Number           uint      `gorm:"index:idx_application_token_number,unique"`
	Messages         []Message `gorm:"foreignkey:ChatID;association_foreignkey:ID"`
	MessagesCount    uint
}
