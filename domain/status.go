package domain

import "fmt"

const (
	StatusDraft     = "draft"
	StatusQueued    = "queued"
	StatusRendering = "rendering"
	StatusReady     = "ready"
	StatusFailed    = "failed"
	StatusArchived  = "archived"
)

func ValidStatus(status string) bool {
	switch status {
	case StatusDraft, StatusQueued, StatusRendering, StatusReady, StatusFailed, StatusArchived:
		return true
	default:
		return false
	}
}

func CanTransition(from, to string) bool {
	if !ValidStatus(from) || !ValidStatus(to) {
		return false
	}
	switch from {
	case StatusDraft:
		return to == StatusQueued || to == StatusArchived
	case StatusQueued:
		return to == StatusRendering || to == StatusFailed || to == StatusArchived
	case StatusRendering:
		return to == StatusReady || to == StatusFailed
	case StatusReady:
		return to == StatusArchived
	case StatusFailed:
		return to == StatusQueued || to == StatusArchived
	case StatusArchived:
		return false
	default:
		return false
	}
}

func TransitionMessage(from, to string) string {
	if CanTransition(from, to) {
		return fmt.Sprintf("status changed from %s to %s", from, to)
	}
	return fmt.Sprintf("status change from %s to %s is not allowed", from, to)
}

func StatusRank(status string) int {
	switch status {
	case StatusDraft:
		return 1
	case StatusQueued:
		return 2
	case StatusRendering:
		return 3
	case StatusReady:
		return 4
	case StatusFailed:
		return 5
	case StatusArchived:
		return 6
	default:
		return 0
	}
}
