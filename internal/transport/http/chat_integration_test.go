package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	chatModel "github.com/Polilo-User/test-task-hitalent/internal/chats/model"
	messagesModel "github.com/Polilo-User/test-task-hitalent/internal/messages/model"

	httptransport "github.com/Polilo-User/test-task-hitalent/internal/transport/http"
	"github.com/Polilo-User/test-task-hitalent/internal/transport/http/mocks"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseChatURL = "/v1/chats/"
	chatURL     = baseChatURL + "%s"
	messageURL  = chatURL + "/messages/"
)

func TestServer_CreateChat_Success(t *testing.T) {
	type args struct {
		chat chatModel.Chat
	}
	tests := []struct {
		name     string
		args     args
		wantChat chatModel.Chat
		wantCode int
	}{
		{
			name: "success",
			args: args{
				chat: chatModel.Chat{
					Title: pointer.ToString("testChat"),
				},
			},
			wantChat: chatModel.Chat{
				ID:        pointer.ToString("1"),
				Title:     pointer.ToString("testChat"),
				CreatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(u, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			u.EXPECT().
				CreateChat(gomock.Any(), gomock.AssignableToTypeOf(&chatModel.Chat{})).
				DoAndReturn(func(ctx context.Context, c *chatModel.Chat) (*chatModel.Chat, error) {
					return &tt.wantChat, nil
				}).Times(1)

			data, err := json.Marshal(tt.args.chat)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, baseChatURL, bytes.NewBuffer(data))
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data chatModel.Chat `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantChat, res.Data)
		})
	}
}

func TestServer_CreateChat_Error(t *testing.T) {
	type args struct {
		chat chatModel.Chat
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "fails",
			args: args{
				chat: chatModel.Chat{
					Title: pointer.ToString("testChat"),
				},
			},
			wantErr:  "test fail",
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(c, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			c.EXPECT().CreateChat(gomock.Any(), &tt.args.chat).Return(nil, errors.New(tt.wantErr)).Times(1)

			data, err := json.Marshal(tt.args.chat)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, baseChatURL, bytes.NewBuffer(data))
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Error string `json:"error"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantErr, res.Error)
		})
	}
}

func TestServer_GetChat_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantChat chatModel.Chat
		wantCode int
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
			wantChat: chatModel.Chat{
				ID:        pointer.ToString("1"),
				Title:     pointer.ToString("testChat"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				Messages:  []messagesModel.Message{},
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(c, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			c.EXPECT().
				GetChat(gomock.Any(), tt.args.id, int64(20)).
				Return(&tt.wantChat, nil).Times(1)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(chatURL, tt.args.id), nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data chatModel.Chat `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.EqualValues(t, tt.wantChat, res.Data)
		})
	}
}

func TestServer_GetChat_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "fails",
			args: args{
				id: "1",
			},
			wantErr:  "test fail",
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(c, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			c.EXPECT().GetChat(gomock.Any(), tt.args.id, int64(20)).Return(nil, errors.New(tt.wantErr)).Times(1)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(chatURL, tt.args.id), nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Error string `json:"error"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantErr, res.Error)
		})
	}
}

func TestServer_DeleteChat_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(c, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			c.EXPECT().DeleteChat(gomock.Any(), tt.args.id).Return(nil).Times(1)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(chatURL, tt.args.id), nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data string `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
		})
	}
}

func TestServer_DeleteChat_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantErr  string
	}{
		{
			name: "fails",
			args: args{
				id: "1",
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  "test fail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(c, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			c.EXPECT().DeleteChat(gomock.Any(), tt.args.id).Return(errors.New(tt.wantErr)).Times(1)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(chatURL, tt.args.id), nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Error string `json:"error"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantErr, res.Error)
		})
	}
}

func TestServer_CreateMessage_Success(t *testing.T) {
	type args struct {
		message messagesModel.Message
	}
	tests := []struct {
		name        string
		args        args
		wantMessage messagesModel.Message
		wantCode    int
	}{
		{
			name: "success",
			args: args{
				message: messagesModel.Message{
					Text:   pointer.ToString("testMessage"),
					ChatID: pointer.ToString("1"),
				},
			},
			wantMessage: messagesModel.Message{
				ID:        pointer.ToString("1"),
				Text:      pointer.ToString("testMessage"),
				CreatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(u, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			m.EXPECT().
				CreateMessage(gomock.Any(), gomock.AssignableToTypeOf(&messagesModel.Message{})).
				DoAndReturn(func(ctx context.Context, c *messagesModel.Message) (*messagesModel.Message, error) {
					return &tt.wantMessage, nil
				}).Times(1)

			data, err := json.Marshal(tt.args.message)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(messageURL, *tt.args.message.ChatID), bytes.NewBuffer(data))
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data messagesModel.Message `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantMessage, res.Data)
		})
	}
}

func TestServer_CreateMessage_Error(t *testing.T) {
	type args struct {
		message messagesModel.Message
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "fails",
			args: args{
				message: messagesModel.Message{
					Text:   pointer.ToString("testMessage"),
					ChatID: pointer.ToString("1"),
				},
			},
			wantErr:  "failed to create message",
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := mocks.NewMockChat(ctrl)
			m := mocks.NewMockMessage(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(c, m, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			m.EXPECT().CreateMessage(gomock.Any(), &tt.args.message).Return(nil, errors.New(tt.wantErr)).Times(1)

			data, err := json.Marshal(tt.args.message)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(messageURL, *tt.args.message.ChatID), bytes.NewBuffer(data))
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Error string `json:"error"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantErr, res.Error)
		})
	}
}
