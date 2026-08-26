package service

import (
	"sort"

	"slowpreview/domain"
	"slowpreview/store"
)

type History struct {
	Events []domain.ActivityEvent
}

func (s *Service) History(previewID string) (History, error) {
	return LoadHistory(s.repo, previewID)
}

func LoadHistoryFromService(s *Service, previewID string) (History, error) {
	return s.History(previewID)
}

func LoadHistory(repo *store.Repository, previewID string) (History, error) {
	events, err := repo.ListEvents(previewID)
	if err != nil {
		return History{}, err
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].Seq < events[j].Seq })
	return History{Events: events}, nil
}

func (h History) Kinds() []string {
	kinds := make([]string, 0, len(h.Events))
	for _, event := range h.Events {
		kinds = append(kinds, event.Kind)
	}
	return kinds
}

func (h History) Contains(kind string) bool {
	for _, event := range h.Events {
		if event.Kind == kind {
			return true
		}
	}
	return false
}
