package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

// Implement driver.Valuer for ChatRoom
func (cr ChatRoom) Value() (driver.Value, error) {
	return json.Marshal(cr)
}

// Implement driver.Scanner for ChatRoom
func (cr *ChatRoom) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ChatRoom scan")
	}

	return json.Unmarshal(data, &cr)
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

// Implement driver.Valuer for ChatMessage
func (cm ChatMessage) Value() (driver.Value, error) {
	return json.Marshal(cm)
}

// Implement driver.Scanner for ChatMessage
func (cm *ChatMessage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ChatMessage scan")
	}

	return json.Unmarshal(data, &cm)
}
