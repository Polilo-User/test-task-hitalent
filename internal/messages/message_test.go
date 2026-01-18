package messages_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/Polilo-User/test-task-hitalent/internal/messages"
	"github.com/Polilo-User/test-task-hitalent/internal/messages/mocks"
	"github.com/Polilo-User/test-task-hitalent/internal/messages/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessages_CreateMessage_Success(t *testing.T) {
	type args struct {
		message *model.Message
	}
	tests := []struct {
		name        string
		args        args
		wantMessage *model.Message
	}{
		{
			name: "success",
			args: args{
				message: &model.Message{
					ChatID: pointer.ToString("1"),
					Text:   pointer.ToString("testMessage"),
				},
			},
			wantMessage: &model.Message{
				ChatID:    pointer.ToString("1"),
				Text:      pointer.ToString("testMessage"),
				CreatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			c := mocks.NewMockChatService(ctrl)

			m := messages.New(s, c)
			require.NotNil(t, m)

			ctx := context.Background()

			s.EXPECT().InsertMessage(gomock.Any(), tt.args.message).Return(tt.wantMessage, nil).Times(1)

			c.EXPECT().ChatExist(gomock.Any(), *tt.args.message.ChatID).Return(nil).Times(1)

			message, err := m.CreateMessage(ctx, tt.args.message)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantMessage, message)
		})
	}
}

func TestMessages_CreateMessage_Error(t *testing.T) {
	type args struct {
		message *model.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				message: &model.Message{
					ChatID: pointer.ToString("1"),
					Text:   pointer.ToString("testMessage"),
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
			c := mocks.NewMockChatService(ctrl)

			m := messages.New(s, c)

			require.NotNil(t, m)

			ctx := context.Background()

			s.EXPECT().InsertMessage(gomock.Any(), tt.args.message).Return(nil, tt.wantErr).Times(1)

			c.EXPECT().ChatExist(gomock.Any(), *tt.args.message.ChatID).Return(nil).Times(1)

			message, err := m.CreateMessage(ctx, tt.args.message)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, message)
		})
	}
}
