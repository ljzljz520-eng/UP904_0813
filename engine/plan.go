package engine

import (
	"fmt"
	"strings"

	"slowpreview/domain"
	"slowpreview/timeline"
)

type PreviewPlan struct {
	Task       domain.RenderTask
	Metrics    timeline.QualityMetrics
	FilterList []string
	Warnings   []string
}

func BuildPlan(asset domain.VideoAsset, spec domain.PreviewSpec) PreviewPlan {
	metrics := timeline.CalculateMetrics(asset, spec)
	filters := buildFilters(spec)
	warnings := assessWarnings(asset, spec, metrics)
	task := domain.RenderTask{ID: "task-" + spec.ID, PreviewID: spec.ID, Speed: spec.Speed, RequestedLabel: spec.RequestedLabel, Resolution: spec.Resolution, Crop: spec.Crop, Interpolate: spec.Interpolate, Status: domain.StatusQueued}
	return PreviewPlan{Task: task, Metrics: metrics, FilterList: filters, Warnings: warnings}
}

func buildFilters(spec domain.PreviewSpec) []string {
	filters := []string{fmt.Sprintf("crop=start=%d:end=%d", spec.Crop.StartMS, spec.Crop.EndMS)}
	if spec.Speed != domain.SpeedNormal {
		filters = append(filters, fmt.Sprintf("setpts=%.2fx*PTS", 1/float64(spec.Speed)))
	}
	if spec.Interpolate {
		filters = append(filters, "minterpolate=fps=60")
	}
	filters = append(filters, fmt.Sprintf("scale=%d:%d", spec.Resolution.Width(), spec.Resolution.Height()))
	return filters
}

func assessWarnings(asset domain.VideoAsset, spec domain.PreviewSpec, metrics timeline.QualityMetrics) []string {
	warnings := make([]string, 0)
	if spec.Speed == domain.SpeedHalf && asset.FrameRate < 30 {
		warnings = append(warnings, "slow motion may look uneven below 30 fps")
	}
	if spec.Interpolate {
		warnings = append(warnings, "interpolation doubles estimated output frames")
	}
	if metrics.OutputMS > 60000 {
		warnings = append(warnings, "output duration is longer than one minute")
	}
	if spec.Resolution.Width() > asset.Width || spec.Resolution.Height() > asset.Height {
		warnings = append(warnings, "selected resolution is larger than source")
	}
	return warnings
}

func (p PreviewPlan) FilterSummary() string {
	return strings.Join(p.FilterList, ", ")
}

func (p PreviewPlan) QualityScore() int {
	return timeline.Score(p.Metrics)
}
