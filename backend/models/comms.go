package models

import (
	"time"
)

type ChatRoom struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:text"`
	Private     bool    `gorm:"not null;default:false"`
	Members     []Users `gorm:"many2many:chat_room_members;"`
	CreatedAt   time.Time
}

func (ChatRoom) TableName() string {
	return "chat_rooms"
}

type ChatMessage struct {
	ID         uint   `gorm:"primaryKey"`
	ChatRoomID uint   `gorm:"index;not null"`
	SenderID   uint   `gorm:"index;not null"`
	Message    string `gorm:"type:text;not null"`
	SentAt     time.Time
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
