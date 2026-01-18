package model

import "time"

type Message struct {
	ID        *string    `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Text      *string    `json:"text" db:"text"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	ChatID    *string    `json:"chat_id" db:"chat_id"`
}
