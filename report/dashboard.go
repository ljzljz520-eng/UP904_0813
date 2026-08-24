package report

import (
	"fmt"
	"sort"
	"strings"

	"slowpreview/catalog"
	"slowpreview/domain"
)

type Dashboard struct {
	AssetCount    int
	PreviewCount  int
	StatusCounts  map[string]int
	Attention     int
	TopPreviews   []catalog.PreviewResult
	QueueHeadline string
}

func BuildDashboard(assets []domain.VideoAsset, records []domain.PreviewRecord, queue catalog.ReviewQueue, score func(domain.PreviewRecord) int) Dashboard {
	ranked := catalog.RankPreviews(records, score)
	top := ranked
	if len(top) > 5 {
		top = top[:5]
	}
	counts := catalog.SummarizeStatuses(records)
	attention := len(queue.Items)
	headline := "review queue is clear"
	if attention > 0 {
		headline = fmt.Sprintf("%d preview(s) need attention", attention)
	}
	return Dashboard{AssetCount: len(assets), PreviewCount: len(records), StatusCounts: counts, Attention: attention, TopPreviews: top, QueueHeadline: headline}
}

func (d Dashboard) StatusLine() string {
	statuses := make([]string, 0, len(d.StatusCounts))
	for status, count := range d.StatusCounts {
		statuses = append(statuses, fmt.Sprintf("%s=%d", status, count))
	}
	sort.Strings(statuses)
	return strings.Join(statuses, ", ")
}

func (d Dashboard) Render() string {
	lines := []string{
		"Slow-motion preview dashboard",
		fmt.Sprintf("Assets: %d", d.AssetCount),
		fmt.Sprintf("Previews: %d", d.PreviewCount),
		"Statuses: " + d.StatusLine(),
		"Attention: " + d.QueueHeadline,
	}
	for _, preview := range d.TopPreviews {
		lines = append(lines, fmt.Sprintf("%s %s %d", preview.Record.Spec.ID, StatusBadge(preview.Record.Status), preview.Score))
	}
	return strings.Join(lines, "\n")
}

func (d Dashboard) HasFailures() bool {
	return d.StatusCounts[domain.StatusFailed] > 0
}

func (d Dashboard) ReadyRatio() float64 {
	if d.PreviewCount == 0 {
		return 0
	}
	return float64(d.StatusCounts[domain.StatusReady]) / float64(d.PreviewCount)
}
