package store

import (
	"context"

	"github.com/Polilo-User/test-task-hitalent/internal/messages/model"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) InsertMessage(ctx context.Context, c *model.Message) (*model.Message, error) {
	if err := s.db.Table("messages").Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Store) GetMessagesByChat(ctx context.Context, id string, limit int64) ([]model.Message, error) {
	var c []model.Message

	if err := s.db.Table("messages").Where("chat_id = ?", id).Limit(int(limit)).Find(&c).Error; err != nil {
		return nil, err
	}

	return c, nil
}
