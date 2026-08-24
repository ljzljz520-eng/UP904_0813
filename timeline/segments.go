package timeline

import (
	"sort"

	"slowpreview/domain"
)

type Segment struct {
	Index    int
	StartMS  int64
	EndMS    int64
	Selected bool
	Label    string
}

func BuildSegments(asset domain.VideoAsset, segmentMS int64) []Segment {
	if segmentMS <= 0 {
		segmentMS = 1000
	}
	segments := make([]Segment, 0)
	for start := int64(0); start < asset.DurationMS; start += segmentMS {
		end := start + segmentMS
		if end > asset.DurationMS {
			end = asset.DurationMS
		}
		segments = append(segments, Segment{Index: len(segments), StartMS: start, EndMS: end, Label: segmentLabel(start, end)})
	}
	return segments
}

func segmentLabel(start, end int64) string {
	return formatMS(start) + " - " + formatMS(end)
}

func formatMS(value int64) string {
	seconds := value / 1000
	minutes := seconds / 60
	return formatUnit(minutes) + ":" + formatUnit(seconds%60)
}

func formatUnit(value int64) string {
	if value < 10 {
		return "0" + string(rune('0'+value))
	}
	return string(rune('0'+value/10)) + string(rune('0'+value%10))
}

func SelectSegments(segments []Segment, crop domain.CropWindow) []Segment {
	selected := make([]Segment, len(segments))
	copy(selected, segments)
	for i := range selected {
		selected[i].Selected = selected[i].StartMS < crop.EndMS && selected[i].EndMS > crop.StartMS
	}
	return selected
}

func MergeSelected(segments []Segment) domain.CropWindow {
	selected := make([]Segment, 0)
	for _, segment := range segments {
		if segment.Selected {
			selected = append(selected, segment)
		}
	}
	if len(selected) == 0 {
		return domain.CropWindow{}
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].StartMS < selected[j].StartMS })
	return domain.CropWindow{StartMS: selected[0].StartMS, EndMS: selected[len(selected)-1].EndMS}
}
