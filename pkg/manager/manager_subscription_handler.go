package manager

import (
	"encoding/json"
	"github.com/stashapp/stash/pkg/logger"
)

type metadataUpdatePayload struct {
	Progress float64          `json:"progress"`
	Message  string           `json:"message"`
	Logs     []logger.LogItem `json:"logs"`
}

func (s *singleton) HandleMetadataUpdateSubscriptionTick(msg chan string) {
	var statusMessage string
	switch instance.Status {
	case Idle:
		statusMessage = "Idle"
	case Import:
		statusMessage = "Import"
	case Export:
		statusMessage = "Export"
	case Scan:
		statusMessage = "Scan"
	case Generate:
		statusMessage = "Generate"
	}
	payload := &metadataUpdatePayload{
		Progress: 0, // TODO
		Message:  statusMessage,
		Logs:     logger.LogCache,
	}
	payloadJSON, _ := json.Marshal(payload)

	msg <- string(payloadJSON)
}
