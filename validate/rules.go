package validate

import (
	"path/filepath"
	"strings"

	"slowpreview/domain"
)

func NormalizeTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return "Untitled training clip"
	}
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) == 1 {
			words[i] = strings.ToUpper(word)
		} else {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func DefaultOutputPath(previewID string) string {
	clean := strings.ReplaceAll(strings.TrimSpace(previewID), " ", "-")
	if clean == "" {
		clean = "preview"
	}
	return filepath.Join("previews", clean+".mp4")
}

func ClampCrop(asset domain.VideoAsset, crop domain.CropWindow) domain.CropWindow {
	if crop.StartMS < 0 {
		crop.StartMS = 0
	}
	if crop.EndMS > asset.DurationMS {
		crop.EndMS = asset.DurationMS
	}
	if crop.EndMS <= crop.StartMS {
		fallback := int64(5000)
		if asset.DurationMS < fallback {
			fallback = asset.DurationMS
		}
		crop.EndMS = crop.StartMS + fallback
		if crop.EndMS > asset.DurationMS {
			crop.StartMS = asset.DurationMS - fallback
			crop.EndMS = asset.DurationMS
		}
	}
	return crop
}

func SuggestedResolution(asset domain.VideoAsset) domain.Resolution {
	if asset.Width >= 2560 && asset.Height >= 1440 {
		return domain.ResolutionHigh
	}
	if asset.Width >= 1920 && asset.Height >= 1080 {
		return domain.ResolutionMedium
	}
	return domain.ResolutionLow
}

func ParseSpeed(label string) (domain.PreviewSpeed, bool) {
	switch strings.ToLower(strings.TrimSpace(label)) {
	case "0.5", "0.5x", "half":
		return domain.SpeedHalf, true
	case "0.75", "0.75x":
		return domain.SpeedThreeQuarter, true
	case "1", "1x", "normal":
		return domain.SpeedNormal, true
	case "2", "2x", "double":
		return domain.SpeedDouble, true
	default:
		return 0, false
	}
}
