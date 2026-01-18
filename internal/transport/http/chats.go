package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Polilo-User/test-task-hitalent/internal/chats/model"
	"github.com/Polilo-User/test-task-hitalent/internal/core/logging"
	"go.uber.org/zap"
)

type deletedChatResponse struct {
	Success bool `json:"success"`
}

func (s *Server) createChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var c model.Chat
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		logging.From(ctx).Error("failed to decode request body", zap.Error(err))
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	createdChat, err := s.chat.CreateChat(ctx, &c)
	if err != nil {
		logging.From(ctx).Error("failed to create chat", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": createdChat,
	})
}

func (s *Server) getChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"invalid chat id"}`, http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	var limit int64 = 20
	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit < 0 {
			http.Error(w, `{"error":"invalid limit parameter"}`, http.StatusBadRequest)
			return
		}
	}

	chat, err := s.chat.GetChat(ctx, id, 20)
	if err != nil {
		logging.From(ctx).Error("failed to get chat", zap.Error(err))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": chat,
	})
}

func (s *Server) deleteChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"invalid chat id"}`, http.StatusBadRequest)
		return
	}

	err = s.chat.DeleteChat(ctx, id)
	if err != nil {
		logging.From(ctx).Error("failed to delete chat", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"data": "deleted",
	})

	json.NewEncoder(w).Encode(deletedChatResponse{Success: true})
}

func extractID(path string) (string, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "chats" {
			return parts[i+1], nil
		}
	}

	return "", errors.New("chat ID not found in path")
}
func (s *Server) SetupRoutes() {
	http.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.createChat(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/chats/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.getChat(w, r)
		case http.MethodDelete:
			s.deleteChat(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
