package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/internal/llm"
	"github.com/stashapp/stash/pkg/logger"
)

// llmRoutes exposes the embedded LLM assistant. All routes sit behind the
// global authenticateHandler middleware (see server.go), so they share stash's
// session/API-key auth.
type llmRoutes struct {
	routes
	service *llm.Service
}

func (rs llmRoutes) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/status", rs.Status)
	r.Post("/chat", rs.Chat)
	r.Post("/confirm", rs.Confirm)
	return r
}

// Status reports whether the assistant is enabled/configured and which tools are
// available. Never exposes the API key.
func (rs llmRoutes) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rs.service.Status())
}

// Chat streams a single assistant turn as Server-Sent Events.
func (rs llmRoutes) Chat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConversationID string `json:"conversationId"`
		Message        string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // disable proxy buffering
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	emit := func(event string, data any) error {
		payload, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, payload); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	if err := rs.service.Chat(r.Context(), req.ConversationID, req.Message, emit); err != nil {
		logger.Errorf("[assistant] chat error: %v", err)
		_ = emit(llm.EventError, map[string]string{"error": err.Error()})
	}
}

// Confirm executes a single write tool that was proposed under the "ask" policy and
// approved by the user in the UI.
func (rs llmRoutes) Confirm(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	out, err := rs.service.ExecuteConfirmed(r.Context(), req.Name, req.Input)
	if err != nil {
		logger.Warnf("[assistant] confirm %q failed: %v", req.Name, err)
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": json.RawMessage(out)})
}
