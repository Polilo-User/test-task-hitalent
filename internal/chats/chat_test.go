package chats_test

import (
	"context"
	"testing"
	"time"

	"github.com/Polilo-User/test-task-hitalent/internal/chats"
	"github.com/Polilo-User/test-task-hitalent/internal/chats/mocks"
	"github.com/Polilo-User/test-task-hitalent/internal/chats/model"
	"github.com/Polilo-User/test-task-hitalent/internal/core/errors"
	modelMessage "github.com/Polilo-User/test-task-hitalent/internal/messages/model"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockStore(ctrl)
	m := mocks.NewMockMessage(ctrl)

	u := chats.New(s, m)
	assert.NotNil(t, u)
}

func TestChats_CreateChat_Success(t *testing.T) {
	type args struct {
		chat *model.Chat
	}
	tests := []struct {
		name     string
		args     args
		wantChat *model.Chat
	}{
		{
			name: "success",
			args: args{
				chat: &model.Chat{
					Title: pointer.ToString("testChat"),
				},
			},
			wantChat: &model.Chat{
				ID:        pointer.ToString("1"),
				Title:     pointer.ToString("testChat"),
				CreatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			m := mocks.NewMockMessage(ctrl)

			c := chats.New(s, m)
			require.NotNil(t, c)

			ctx := context.Background()

			s.EXPECT().InsertChat(gomock.Any(), tt.args.chat).Return(tt.wantChat, nil).Times(1)

			chat, err := c.CreateChat(ctx, tt.args.chat)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantChat, chat)
		})
	}
}

func TestChats_CreateChat_Error(t *testing.T) {
	type args struct {
		chat *model.Chat
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				chat: &model.Chat{
					Title: pointer.ToString("testChat"),
				},
			},
			wantErr: errors.New("test fail"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			m := mocks.NewMockMessage(ctrl)

			c := chats.New(s, m)

			require.NotNil(t, c)

			ctx := context.Background()

			s.EXPECT().InsertChat(gomock.Any(), tt.args.chat).Return(nil, tt.wantErr).Times(1)

			chat, err := c.CreateChat(ctx, tt.args.chat)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, chat)
		})
	}
}

func TestChats_GetChat_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantchat *model.Chat
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
			wantchat: &model.Chat{
				ID:        pointer.ToString("1"),
				Title:     pointer.ToString("testChat"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			m := mocks.NewMockMessage(ctrl)

			c := chats.New(s, m)
			require.NotNil(t, c)

			ctx := context.Background()

			s.EXPECT().
				GetChat(gomock.Any(), tt.args.id).
				Return(tt.wantchat, nil).
				Times(1)

			m.EXPECT().
				GetMessagesByChat(gomock.Any(), tt.args.id, int64(20)).
				Return([]modelMessage.Message{}, nil).
				Times(1)

			chat, err := c.GetChat(ctx, tt.args.id, 20)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantchat, chat)
		})
	}
}

func TestChats_GetChat_Error(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				id: "1",
			},
			wantErr: errors.New("test fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			m := mocks.NewMockMessage(ctrl)

			c := chats.New(s, m)
			require.NotNil(t, c)

			ctx := context.Background()

			s.EXPECT().
				GetChat(gomock.Any(), tt.args.id).
				Return(nil, tt.wantErr).
				Times(1)

			chat, err := c.GetChat(ctx, tt.args.id, int64(20))

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, chat)
		})
	}
}

func TestChats_DeleteChat_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			m := mocks.NewMockMessage(ctrl)

			c := chats.New(s, m)
			require.NotNil(t, c)

			ctx := context.Background()

			s.EXPECT().DeleteChat(gomock.Any(), tt.args.id).Return(nil).Times(1)

			err := c.DeleteChat(ctx, tt.args.id)
			assert.NoError(t, err)
		})
	}
}

func TestChats_DeleteChat_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				id: "1",
			},
			wantErr: errors.New("test fail"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			m := mocks.NewMockMessage(ctrl)

			c := chats.New(s, m)
			require.NotNil(t, c)

			ctx := context.Background()

			s.EXPECT().DeleteChat(gomock.Any(), tt.args.id).Return(tt.wantErr).Times(1)

			err := c.DeleteChat(ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
