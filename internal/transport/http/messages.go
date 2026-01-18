package http

import (
	"encoding/json"
	"net/http"

	"github.com/Polilo-User/test-task-hitalent/internal/chats"
	"github.com/Polilo-User/test-task-hitalent/internal/core/errors"
	"github.com/Polilo-User/test-task-hitalent/internal/core/logging"
	"github.com/Polilo-User/test-task-hitalent/internal/messages/model"

	"go.uber.org/zap"
)

func (s *Server) createMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var c model.Message
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		logging.From(ctx).Error("failed to decode request body", zap.Error(err))
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"invalid chat id"}`, http.StatusBadRequest)
		return
	}

	c.ChatID = &id

	createdMessage, err := s.message.CreateMessage(ctx, &c)
	if err != nil {
		logging.From(ctx).Error("failed to create message", zap.Error(err))

		if errors.Is(err, chats.ErrChatNotFound) {
			http.Error(w, `{"error":"chat not found"}`, http.StatusNotFound)
			return
		}

		http.Error(w, `{"error":"failed to create message"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": createdMessage,
	})
}
