package catalog

import (
	"sort"
	"strings"

	"slowpreview/domain"
)

type PreviewQuery struct {
	Text       string
	Statuses   []string
	Speed      domain.PreviewSpeed
	Resolution domain.Resolution
	MinQuality int
}

type PreviewResult struct {
	Record domain.PreviewRecord
	Score  int
	Match  string
}

func MatchPreview(record domain.PreviewRecord, query PreviewQuery, score int) bool {
	if query.Text != "" {
		text := strings.ToLower(strings.TrimSpace(query.Text))
		joined := strings.ToLower(record.Spec.ID + " " + record.Spec.AssetID + " " + record.Spec.RequestedLabel)
		if !strings.Contains(joined, text) {
			return false
		}
	}
	if len(query.Statuses) > 0 && !containsStatus(query.Statuses, record.Status) {
		return false
	}
	if query.Speed != 0 && record.Spec.Speed != query.Speed {
		return false
	}
	if query.Resolution != "" && record.Spec.Resolution != query.Resolution {
		return false
	}
	return score >= query.MinQuality
}

func containsStatus(statuses []string, status string) bool {
	for _, candidate := range statuses {
		if candidate == status {
			return true
		}
	}
	return false
}

func RankPreviews(records []domain.PreviewRecord, score func(domain.PreviewRecord) int) []PreviewResult {
	results := make([]PreviewResult, 0, len(records))
	for _, record := range records {
		value := score(record)
		match := "standard"
		if value >= 80 {
			match = "high confidence"
		} else if value < 40 {
			match = "review needed"
		}
		results = append(results, PreviewResult{Record: record, Score: value, Match: match})
	}
	sort.SliceStable(results, func(a, b int) bool {
		if results[a].Score == results[b].Score {
			return results[a].Record.Spec.ID < results[b].Record.Spec.ID
		}
		return results[a].Score > results[b].Score
	})
	return results
}

func GroupByStatus(records []domain.PreviewRecord) map[string][]domain.PreviewRecord {
	groups := make(map[string][]domain.PreviewRecord)
	for _, record := range records {
		groups[record.Status] = append(groups[record.Status], record)
	}
	for status := range groups {
		sort.Slice(groups[status], func(i, j int) bool { return groups[status][i].Spec.ID < groups[status][j].Spec.ID })
	}
	return groups
}

func SummarizeStatuses(records []domain.PreviewRecord) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}
