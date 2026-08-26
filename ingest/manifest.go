package ingest

import (
	"fmt"
	"strconv"
	"strings"

	"slowpreview/domain"
)

type ManifestEntry struct {
	ID         string
	Path       string
	Title      string
	DurationMS int64
	FrameRate  int
	Width      int
	Height     int
}

type Manifest struct {
	Entries []ManifestEntry
}

func ParseManifest(lines []string) (Manifest, []string) {
	manifest := Manifest{Entries: make([]ManifestEntry, 0)}
	issues := make([]string, 0)
	for lineNumber, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entry, err := parseEntry(line)
		if err != nil {
			issues = append(issues, fmt.Sprintf("line %d: %v", lineNumber+1, err))
			continue
		}
		manifest.Entries = append(manifest.Entries, entry)
	}
	return manifest, issues
}

func parseEntry(line string) (ManifestEntry, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 7 {
		return ManifestEntry{}, fmt.Errorf("expected 7 fields, received %d", len(parts))
	}
	duration, err := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("duration is not an integer")
	}
	frameRate, err := strconv.Atoi(strings.TrimSpace(parts[4]))
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("frame rate is not an integer")
	}
	width, err := strconv.Atoi(strings.TrimSpace(parts[5]))
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("width is not an integer")
	}
	height, err := strconv.Atoi(strings.TrimSpace(parts[6]))
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("height is not an integer")
	}
	return ManifestEntry{ID: strings.TrimSpace(parts[0]), Path: strings.TrimSpace(parts[1]), Title: strings.TrimSpace(parts[2]), DurationMS: duration, FrameRate: frameRate, Width: width, Height: height}, nil
}

func (m Manifest) Assets() []domain.VideoAsset {
	assets := make([]domain.VideoAsset, 0, len(m.Entries))
	for _, entry := range m.Entries {
		assets = append(assets, domain.VideoAsset{ID: entry.ID, SourcePath: entry.Path, Title: entry.Title, DurationMS: entry.DurationMS, FrameRate: entry.FrameRate, Width: entry.Width, Height: entry.Height})
	}
	return assets
}

func (m Manifest) ValidEntries() []ManifestEntry {
	valid := make([]ManifestEntry, 0, len(m.Entries))
	for _, entry := range m.Entries {
		if entry.ID != "" && entry.Path != "" && entry.DurationMS > 0 && entry.FrameRate > 0 && entry.Width > 0 && entry.Height > 0 {
			valid = append(valid, entry)
		}
	}
	return valid
}

func (m Manifest) Find(id string) (ManifestEntry, bool) {
	for _, entry := range m.Entries {
		if entry.ID == id {
			return entry, true
		}
	}
	return ManifestEntry{}, false
}
