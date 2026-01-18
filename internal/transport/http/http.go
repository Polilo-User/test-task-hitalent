//go:generate mockgen -destination=./mocks/http_mock.go -package mocks github.com/Polilo-User/test-task-hitalent/internal/transport/http Users,DB

package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	chmodel "github.com/Polilo-User/test-task-hitalent/internal/chats/model"
	msmodel "github.com/Polilo-User/test-task-hitalent/internal/messages/model"
	"github.com/gorilla/mux"
)

// go:generate mockgen -destination=./mocks/http_mock.go -package=mocks github.com/Polilo-User/test-task-hitalent/internal/transport/http Chat,Message,DB

type Chat interface {
	CreateChat(ctx context.Context, chat *chmodel.Chat) (*chmodel.Chat, error)
	GetChat(ctx context.Context, id string, limit int64) (*chmodel.Chat, error)
	DeleteChat(ctx context.Context, id string) error
}

type Message interface {
	CreateMessage(ctx context.Context, message *msmodel.Message) (*msmodel.Message, error)
	GetMessagesByChat(ctx context.Context, id string, limit int64) ([]msmodel.Message, error)
}

type DB interface {
	DB() (*sql.DB, error)
}

type Server struct {
	chat    Chat
	message Message
	db      DB
}

func New(c Chat, m Message, db DB) *Server {
	return &Server{
		chat:    c,
		message: m,
		db:      db,
	}
}

func (s *Server) AddRoutes(r *mux.Router) error {
	r.HandleFunc("/health", s.healthCheck).Methods(http.MethodGet)

	r = r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/chats/", s.createChat).Methods(http.MethodPost)                  // Done
	r.HandleFunc("/chats/{id}", s.getChat).Methods(http.MethodGet)                  // Done
	r.HandleFunc("/chats/{id}", s.deleteChat).Methods(http.MethodDelete)            // Done
	r.HandleFunc("/chats/{id}/messages/", s.createMessage).Methods(http.MethodPost) // Done

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	sql, err := s.db.DB()
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}
	if err := sql.PingContext(r.Context()); err != nil {
		handleError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	jsonRes := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}

	dataBytes, err := json.Marshal(jsonRes)
	if err != nil {
		handleError(ctx, w, err)
		return
	}

	if _, err := w.Write(dataBytes); err != nil {
		handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
