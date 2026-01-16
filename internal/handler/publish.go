package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/k-omotani/aeron-sample/internal/aeron"
	"github.com/k-omotani/aeron-sample/internal/message"
)

// PublishHandler handles publishing messages via HTTP API
type PublishHandler struct {
	publisher *aeron.Publisher
	logger    *slog.Logger
}

// NewPublishHandler creates a new publish handler
func NewPublishHandler(publisher *aeron.Publisher, logger *slog.Logger) *PublishHandler {
	return &PublishHandler{
		publisher: publisher,
		logger:    logger.With("handler", "publish"),
	}
}

// PublishRequest is the request body for publishing
type PublishRequest struct {
	Amount int64 `json:"amount"`
}

// PublishResponse is the response for publish operations
type PublishResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

// Increment handles POST /api/counter/increment
func (h *PublishHandler) Increment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount == 0 {
		req.Amount = 1 // Default increment
	}

	requestID := uuid.New().String()

	msg, err := message.NewIncrementMessage(requestID, req.Amount, "http")
	if err != nil {
		h.logger.Error("failed to create message", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := h.publisher.Publish(ctx, msg); err != nil {
		h.logger.Error("failed to publish message", "error", err)
		http.Error(w, "failed to publish", http.StatusInternalServerError)
		return
	}

	h.logger.Info("increment message published",
		"requestID", requestID,
		"amount", req.Amount,
	)

	resp := PublishResponse{
		RequestID: requestID,
		Status:    "published",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
