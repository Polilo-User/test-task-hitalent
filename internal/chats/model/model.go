package model

import (
	"time"

	"github.com/Polilo-User/test-task-hitalent/internal/messages/model"
)

type Chat struct {
	ID        *string         `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Title     *string         `json:"title" db:"title"`
	CreatedAt *time.Time      `json:"created_at" db:"created_at"`
	Messages  []model.Message `json:"messages"`
}
