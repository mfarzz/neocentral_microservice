package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/notification-service/internal/service"
	"neocentral-go/pkg/response"

	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct {
	svc *service.NotificationService
	sse *service.SSEManager
}

func NewNotificationHandler(svc *service.NotificationService, sse *service.SSEManager) *NotificationHandler {
	return &NotificationHandler{svc: svc, sse: sse}
}

func (h *NotificationHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	notifs, err := h.svc.GetUserNotifications(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": notifs,
	})
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.svc.MarkAsRead(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Notification marked as read",
	})
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.svc.MarkAllAsRead(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "All notifications marked as read",
	})
}

// SSE Endpoint
func (h *NotificationHandler) Stream(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Set required headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create channel for this client
	ch := make(chan []byte)
	h.sse.AddClient(userID, ch)
	defer h.sse.RemoveClient(userID, ch)

	// Send initial connection success message
	initialMsg, _ := json.Marshal(map[string]string{"message": "Connected to Notification Stream"})
	w.Write([]byte("event: connected\ndata: "))
	w.Write(initialMsg)
	w.Write([]byte("\n\n"))
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return // Client disconnected
		case payload := <-ch:
			// Send event to client
			w.Write([]byte("event: notification\ndata: "))
			w.Write(payload)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		}
	}
}
