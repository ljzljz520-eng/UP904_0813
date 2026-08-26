package catalog

import (
	"fmt"
	"sort"

	"slowpreview/domain"
)

type ReviewQueue struct {
	Items []ReviewItem
}

type ReviewItem struct {
	PreviewID string
	AssetID   string
	Reason    string
	Priority  int
	Status    string
}

func BuildReviewQueue(records []domain.PreviewRecord) ReviewQueue {
	items := make([]ReviewItem, 0)
	for _, record := range records {
		reason, priority := reviewReason(record)
		if reason == "" {
			continue
		}
		items = append(items, ReviewItem{PreviewID: record.Spec.ID, AssetID: record.Spec.AssetID, Reason: reason, Priority: priority, Status: record.Status})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].PreviewID < items[j].PreviewID
		}
		return items[i].Priority > items[j].Priority
	})
	return ReviewQueue{Items: items}
}

func reviewReason(record domain.PreviewRecord) (string, int) {
	if record.Status == domain.StatusFailed {
		return "render failed", 100
	}
	if record.Status == domain.StatusDraft {
		return "waiting to be queued", 70
	}
	if record.Spec.Interpolate && record.Spec.Speed == domain.SpeedHalf {
		return "verify interpolated slow motion", 60
	}
	if record.Spec.Resolution == domain.ResolutionHigh && record.Status == domain.StatusReady {
		return "confirm high resolution output", 40
	}
	return "", 0
}

func FormatQueue(queue ReviewQueue) string {
	if len(queue.Items) == 0 {
		return "review queue is clear"
	}
	lines := make([]string, 0, len(queue.Items)+1)
	lines = append(lines, fmt.Sprintf("review queue: %d item(s)", len(queue.Items)))
	for _, item := range queue.Items {
		lines = append(lines, fmt.Sprintf("%s [%s] %s (priority %d)", item.PreviewID, item.Status, item.Reason, item.Priority))
	}
	return joinLines(lines)
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}
