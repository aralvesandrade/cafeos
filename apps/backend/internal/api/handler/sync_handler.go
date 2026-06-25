package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/infra/messaging"
)

type SyncHandler struct {
	publisher *messaging.Publisher
}

func NewSyncHandler(publisher *messaging.Publisher) *SyncHandler {
	return &SyncHandler{publisher: publisher}
}

type syncBatchRequest struct {
	Batch []syncItemRequest `json:"batch"`
}

type syncItemRequest struct {
	ClientID        string      `json:"client_id"`
	EventType       string      `json:"event_type"`
	Payload         interface{} `json:"payload"`
	ClientTimestamp string      `json:"client_timestamp"`
}

type syncResponse struct {
	Accepted int              `json:"accepted"`
	Errors   []syncErrorItem  `json:"errors"`
}

type syncErrorItem struct {
	ClientID string `json:"client_id"`
	Error    string `json:"error"`
}

func (h *SyncHandler) Sync(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req syncBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Batch) == 0 {
		writeError(w, "empty batch", http.StatusBadRequest)
		return
	}

	resp := syncResponse{}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	for _, item := range req.Batch {
		msg := messaging.SyncMessage{
			ClientID:        item.ClientID,
			EventType:       item.EventType,
			TenantID:        tenantID,
			Payload:         item.Payload,
			ClientTimestamp: item.ClientTimestamp,
			PublishedAt:     time.Now(),
		}

		queue := queueForEvent(item.EventType)
		if err := h.publisher.Publish(ctx, queue, msg); err != nil {
			resp.Errors = append(resp.Errors, syncErrorItem{
				ClientID: item.ClientID,
				Error:    err.Error(),
			})
			continue
		}
		resp.Accepted++
	}

	if len(resp.Errors) > 0 {
		writeJSON(w, resp, http.StatusAccepted)
		return
	}

	writeJSON(w, resp, http.StatusAccepted)
}

func queueForEvent(eventType string) string {
	switch eventType {
	case "operation.created", "operation.updated", "operation.deleted":
		return "sync.operations"
	case "stock.created", "stock.moved":
		return "sync.stock"
	case "harvest.production":
		return "sync.harvest"
	case "financial.created":
		return "sync.financial"
	case "labor.shift":
		return "sync.labor"
	default:
		return "sync.operations"
	}
}
