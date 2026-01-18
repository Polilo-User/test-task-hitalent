package store

import (
	"context"

	"github.com/Polilo-User/test-task-hitalent/internal/chats/model"
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

func (s *Store) InsertChat(ctx context.Context, c *model.Chat) (*model.Chat, error) {
	err := s.db.Create(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Store) GetChat(ctx context.Context, id string) (*model.Chat, error) {
	var c model.Chat

	if err := s.db.Table("chats").Where("id = ?", id).Take(&c).Error; err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *Store) DeleteChat(ctx context.Context, id string) error {
	if err := s.db.Table("chats").Delete(&model.Chat{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil

}

func (s *Store) ChatExist(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.db.Model(new(model.Chat)).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).Error
	return exists, err
}
