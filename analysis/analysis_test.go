package analysis

import (
	"testing"

	"slowpreview/domain"
)

func TestRecommendationAndChecklist(t *testing.T) {
	asset := domain.VideoAsset{ID: "analysis", DurationMS: 12000, FrameRate: 30, Width: 1280, Height: 720}
	crop := domain.CropWindow{StartMS: 0, EndMS: 4000}
	recommendations := Recommend(asset, crop)
	if len(recommendations) == 0 {
		t.Fatal("expected recommendations")
	}
	spec := domain.PreviewSpec{ID: "analysis-preview", AssetID: asset.ID, Speed: domain.SpeedThreeQuarter, Resolution: domain.ResolutionLow, Crop: crop, OutputPath: "out.mp4"}
	checklist := RunChecklist(asset, spec)
	if !checklist.Passed() {
		t.Fatalf("checklist failed: %s", checklist.Summary())
	}
}
