package engine

import (
	"strings"
	"testing"

	"slowpreview/domain"
)

func TestCommandBuilderIncludesFilters(t *testing.T) {
	asset := domain.VideoAsset{SourcePath: "clip.mp4", DurationMS: 10000, FrameRate: 60, Width: 1920, Height: 1080}
	spec := domain.PreviewSpec{ID: "cmd", Speed: domain.SpeedHalf, RequestedLabel: "0.5x", Resolution: domain.ResolutionMedium, Crop: domain.CropWindow{StartMS: 0, EndMS: 3000}, OutputPath: "out.mp4", Interpolate: true}
	command := NewCommandBuilder().Build(asset, BuildPlan(asset, spec))
	if !strings.Contains(command, "scale=1920:1080") || !strings.Contains(command, "minterpolate=fps=60") {
		t.Fatalf("command missing filters: %s", command)
	}
}
