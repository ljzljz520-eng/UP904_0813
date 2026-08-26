package playback

import (
	"sort"
	"strings"

	"slowpreview/domain"
)

type Marker struct {
	ID       string
	Position int64
	Label    string
	Kind     string
}

type MarkerBoard struct {
	markers map[string]Marker
}

func NewMarkerBoard() *MarkerBoard {
	return &MarkerBoard{markers: make(map[string]Marker)}
}

func (b *MarkerBoard) Add(marker Marker, spec domain.PreviewSpec) bool {
	if b == nil || strings.TrimSpace(marker.ID) == "" {
		return false
	}
	if marker.Position < 0 || marker.Position > spec.DurationMS() {
		return false
	}
	if marker.Kind == "" {
		marker.Kind = "note"
	}
	b.markers[marker.ID] = marker
	return true
}

func (b *MarkerBoard) Remove(id string) bool {
	if b == nil {
		return false
	}
	if _, ok := b.markers[id]; !ok {
		return false
	}
	delete(b.markers, id)
	return true
}

func (b *MarkerBoard) List() []Marker {
	if b == nil {
		return nil
	}
	markers := make([]Marker, 0, len(b.markers))
	for _, marker := range b.markers {
		markers = append(markers, marker)
	}
	sort.SliceStable(markers, func(i, j int) bool {
		if markers[i].Position == markers[j].Position {
			return markers[i].ID < markers[j].ID
		}
		return markers[i].Position < markers[j].Position
	})
	return markers
}

func (b *MarkerBoard) At(position, tolerance int64) []Marker {
	if tolerance < 0 {
		tolerance = 0
	}
	result := make([]Marker, 0)
	for _, marker := range b.List() {
		distance := marker.Position - position
		if distance < 0 {
			distance = -distance
		}
		if distance <= tolerance {
			result = append(result, marker)
		}
	}
	return result
}

func (b *MarkerBoard) Labels() []string {
	labels := make([]string, 0)
	for _, marker := range b.List() {
		labels = append(labels, marker.Label)
	}
	return labels
}
