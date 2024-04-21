package domain

import "time"

type Message struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	ChatID    uint
	Number    uint
	Content   string
}
