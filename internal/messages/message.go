package messages

import (
	"context"

	"github.com/Polilo-User/test-task-hitalent/internal/messages/model"
)

type Store interface {
	GetMessagesByChat(ctx context.Context, id string, limit int64) ([]model.Message, error)
	InsertMessage(ctx context.Context, c *model.Message) (*model.Message, error)
}

type ChatService interface {
	ChatExist(ctx context.Context, id string) error
}

type MessageService struct {
	store Store
	c     ChatService
}

func New(s Store, c ChatService) *MessageService {
	return &MessageService{
		store: s,
		c:     c,
	}
}

func (c *MessageService) CreateMessage(ctx context.Context, m *model.Message) (*model.Message, error) {
	err := c.c.ChatExist(ctx, *m.ChatID)
	if err != nil {
		return nil, err
	}
	return c.store.InsertMessage(ctx, m)
}

func (c *MessageService) GetMessagesByChat(ctx context.Context, id string, limit int64) ([]model.Message, error) {
	return c.store.GetMessagesByChat(ctx, id, limit)
}
