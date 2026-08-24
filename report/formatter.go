package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"slowpreview/domain"
)

type Export struct {
	PreviewID string `json:"preview_id"`
	Status    string `json:"status"`
	Summary   string `json:"summary"`
	Command   string `json:"command"`
	Events    int    `json:"events"`
}

func BuildExport(record domain.PreviewRecord, eventCount int) Export {
	return Export{PreviewID: record.Spec.ID, Status: record.Status, Summary: record.Message, Command: record.Task.Command, Events: eventCount}
}

func MarshalExport(export Export) ([]byte, error) {
	return json.MarshalIndent(export, "", "  ")
}

func StatusBadge(status string) string {
	switch status {
	case domain.StatusReady:
		return "READY"
	case domain.StatusFailed:
		return "FAILED"
	case domain.StatusArchived:
		return "ARCHIVED"
	default:
		return strings.ToUpper(status)
	}
}

func ProgressText(status string) string {
	switch status {
	case domain.StatusDraft:
		return "configuration is being prepared"
	case domain.StatusQueued:
		return "waiting for render"
	case domain.StatusRendering:
		return "rendering preview"
	case domain.StatusReady:
		return "preview available"
	case domain.StatusFailed:
		return "needs correction"
	default:
		return fmt.Sprintf("status=%s", status)
	}
}
