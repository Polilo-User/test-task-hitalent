package chats

import (
	"context"

	"github.com/Polilo-User/test-task-hitalent/internal/chats/model"
	"github.com/Polilo-User/test-task-hitalent/internal/core/errors"
	"github.com/Polilo-User/test-task-hitalent/internal/core/logging"
	messageModel "github.com/Polilo-User/test-task-hitalent/internal/messages/model"
	"go.uber.org/zap"
)

const (
	ErrChatNotFound = errors.Error("chat_not_found: chat not found")
)

type Store interface {
	InsertChat(ctx context.Context, user *model.Chat) (*model.Chat, error)
	GetChat(ctx context.Context, id string) (*model.Chat, error)
	DeleteChat(ctx context.Context, id string) error
	ChatExist(ctx context.Context, id string) (bool, error)
}

type Message interface {
	GetMessagesByChat(ctx context.Context, id string, limit int64) ([]messageModel.Message, error)
}

type ChatService struct {
	store    Store
	messages Message
}

func New(s Store, m Message) *ChatService {
	return &ChatService{
		store:    s,
		messages: m,
	}
}

func (c *ChatService) CreateChat(ctx context.Context, chat *model.Chat) (*model.Chat, error) {
	return c.store.InsertChat(ctx, chat)
}

func (c *ChatService) GetChat(ctx context.Context, id string, limit int64) (*model.Chat, error) {
	ch, err := c.store.GetChat(ctx, id)
	if err != nil {
		return nil, err
	}

	messages, err := c.messages.GetMessagesByChat(ctx, id, limit)
	if err != nil {
		return nil, err
	}

	ch.Messages = messages

	return ch, nil
}

func (c *ChatService) DeleteChat(ctx context.Context, id string) error {
	return c.store.DeleteChat(ctx, id)
}

func (c *ChatService) ChatExist(ctx context.Context, id string) error {
	ex, err := c.store.ChatExist(ctx, id)
	if err != nil {
		return err
	}
	if !ex {
		logging.From(ctx).Error("chat not found", zap.String("chat_id", id), zap.Error(err))
		return ErrChatNotFound
	}
	return nil
}
