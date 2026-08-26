package ingest

import (
	"path/filepath"
	"strings"

	"slowpreview/domain"
	"slowpreview/validate"
)

type NormalizedAsset struct {
	Asset    domain.VideoAsset
	Warnings []string
}

func Normalize(entry ManifestEntry, sequence int) NormalizedAsset {
	asset := domain.VideoAsset{ID: strings.TrimSpace(entry.ID), SourcePath: filepath.Clean(strings.TrimSpace(entry.Path)), Title: validate.NormalizeTitle(entry.Title), DurationMS: entry.DurationMS, FrameRate: entry.FrameRate, Width: entry.Width, Height: entry.Height, CreatedStep: sequence}
	warnings := make([]string, 0)
	if filepath.IsAbs(asset.SourcePath) {
		warnings = append(warnings, "source path is absolute")
	}
	if asset.FrameRate < 24 {
		warnings = append(warnings, "source frame rate is below cinematic baseline")
	} else if asset.FrameRate > 120 {
		warnings = append(warnings, "source frame rate is high; preview may be expensive")
	}
	if asset.AspectRatio() > 2.2 {
		warnings = append(warnings, "wide aspect ratio may crop action")
	} else if asset.AspectRatio() < 1.2 {
		warnings = append(warnings, "tall aspect ratio may need portrait review")
	}
	return NormalizedAsset{Asset: asset, Warnings: warnings}
}

func NormalizeManifest(manifest Manifest) ([]NormalizedAsset, []string) {
	assets := make([]NormalizedAsset, 0, len(manifest.Entries))
	issues := make([]string, 0)
	for index, entry := range manifest.Entries {
		normalized := Normalize(entry, index+1)
		if !normalized.Asset.Valid() {
			issues = append(issues, "entry "+entry.ID+" is incomplete")
			continue
		}
		assets = append(assets, normalized)
	}
	return assets, issues
}

func CanonicalPath(path string) string {
	clean := filepath.Clean(strings.TrimSpace(path))
	if clean == "." {
		return ""
	}
	return clean
}

func MergeWarnings(groups ...[]string) []string {
	seen := make(map[string]bool)
	merged := make([]string, 0)
	for _, group := range groups {
		for _, warning := range group {
			if !seen[warning] {
				seen[warning] = true
				merged = append(merged, warning)
			}
		}
	}
	return merged
}
