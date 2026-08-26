package report

import (
	"fmt"
	"strings"

	"slowpreview/domain"
	"slowpreview/engine"
	"slowpreview/timeline"
)

type PreviewSummary struct {
	Title      string
	Status     string
	Speed      string
	Resolution string
	Crop       string
	Duration   string
	FrameCount string
	Quality    string
	Warnings   []string
}

func BuildSummary(asset domain.VideoAsset, spec domain.PreviewSpec, plan engine.PreviewPlan) PreviewSummary {
	metrics := plan.Metrics
	return PreviewSummary{Title: asset.Title, Status: domain.StatusDraft, Speed: spec.Speed.Label(), Resolution: string(spec.Resolution), Crop: formatCrop(spec.Crop), Duration: formatDuration(metrics.OutputMS), FrameCount: fmt.Sprintf("%d input -> %d output frames", metrics.InputFrames, metrics.OutputFrames), Quality: qualityLabel(plan.QualityScore()), Warnings: append([]string(nil), plan.Warnings...)}
}

func formatCrop(crop domain.CropWindow) string {
	return fmt.Sprintf("%s to %s", formatDuration(crop.StartMS), formatDuration(crop.EndMS))
}

func formatDuration(milliseconds int64) string {
	seconds := milliseconds / 1000
	return fmt.Sprintf("%02d:%02d", seconds/60, seconds%60)
}

func qualityLabel(score int) string {
	switch {
	case score >= 80:
		return "excellent"
	case score >= 60:
		return "good"
	case score >= 40:
		return "balanced"
	default:
		return "rough"
	}
}

func RenderText(summary PreviewSummary) string {
	lines := []string{"Slow-motion preview", "Title: " + summary.Title, "Status: " + summary.Status, "Speed: " + summary.Speed, "Resolution: " + summary.Resolution, "Crop: " + summary.Crop, "Duration: " + summary.Duration, "Frames: " + summary.FrameCount, "Quality: " + summary.Quality}
	if len(summary.Warnings) > 0 {
		lines = append(lines, "Warnings: "+strings.Join(summary.Warnings, " | "))
	}
	return strings.Join(lines, "\n")
}

func TimelineLabel(asset domain.VideoAsset, spec domain.PreviewSpec) string {
	return fmt.Sprintf("%s at %s (%s)", asset.Title, timeline.FrameRateLabel(asset, spec.Interpolate), spec.Speed.Label())
}
