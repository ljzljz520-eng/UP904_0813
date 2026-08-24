package validate

import (
	"fmt"
	"strings"

	"slowpreview/domain"
)

type Issue struct {
	Field   string
	Message string
}

func (i Issue) Error() string {
	return fmt.Sprintf("%s: %s", i.Field, i.Message)
}

type Result struct {
	Issues []Issue
}

func (r Result) OK() bool {
	return len(r.Issues) == 0
}

func (r Result) Error() string {
	if r.OK() {
		return ""
	}
	parts := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		parts = append(parts, issue.Error())
	}
	return strings.Join(parts, "; ")
}

func Asset(asset domain.VideoAsset) Result {
	result := Result{}
	if strings.TrimSpace(asset.ID) == "" {
		result.Issues = append(result.Issues, Issue{"id", "asset id is required"})
	}
	if strings.TrimSpace(asset.SourcePath) == "" {
		result.Issues = append(result.Issues, Issue{"source_path", "source path is required"})
	}
	if strings.TrimSpace(asset.Title) == "" {
		result.Issues = append(result.Issues, Issue{"title", "title is required"})
	}
	if asset.DurationMS <= 0 {
		result.Issues = append(result.Issues, Issue{"duration_ms", "duration must be positive"})
	}
	if asset.FrameRate < 1 || asset.FrameRate > 240 {
		result.Issues = append(result.Issues, Issue{"frame_rate", "frame rate must be between 1 and 240"})
	}
	if asset.Width < 320 || asset.Height < 180 {
		result.Issues = append(result.Issues, Issue{"dimensions", "video dimensions are too small"})
	}
	return result
}

func Preview(asset domain.VideoAsset, spec domain.PreviewSpec) Result {
	result := Result{}
	if !asset.Valid() {
		result.Issues = append(result.Issues, Issue{"asset", "asset is incomplete"})
	}
	if strings.TrimSpace(spec.ID) == "" {
		result.Issues = append(result.Issues, Issue{"id", "preview id is required"})
	}
	if spec.AssetID != asset.ID {
		result.Issues = append(result.Issues, Issue{"asset_id", "preview must reference selected asset"})
	}
	if !spec.Speed.Valid() {
		result.Issues = append(result.Issues, Issue{"speed", "speed must be 0.5x, 0.75x, 1x, or 2x"})
	}
	if !spec.Resolution.Valid() {
		result.Issues = append(result.Issues, Issue{"resolution", "unsupported preview resolution"})
	}
	if spec.Crop.StartMS < 0 {
		result.Issues = append(result.Issues, Issue{"crop.start_ms", "crop start cannot be negative"})
	}
	if spec.Crop.EndMS <= spec.Crop.StartMS {
		result.Issues = append(result.Issues, Issue{"crop.end_ms", "crop must have positive duration"})
	}
	if spec.Crop.EndMS > asset.DurationMS {
		result.Issues = append(result.Issues, Issue{"crop.end_ms", "crop cannot exceed asset duration"})
	}
	if spec.Crop.DurationMS() > 30000 {
		result.Issues = append(result.Issues, Issue{"crop", "preview crop cannot exceed 30 seconds"})
	}
	if strings.TrimSpace(spec.OutputPath) == "" {
		result.Issues = append(result.Issues, Issue{"output_path", "output path is required"})
	}
	return result
}
