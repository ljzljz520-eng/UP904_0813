package timeline

import (
	"testing"

	"slowpreview/domain"
)

func TestTimelineWorkflow(t *testing.T) {
	asset := domain.VideoAsset{DurationMS: 5000, FrameRate: 50, Width: 1920, Height: 1080}
	spec := domain.PreviewSpec{Speed: domain.SpeedHalf, Crop: domain.CropWindow{StartMS: 0, EndMS: 4000}}
	points := FramePoints(spec, asset)
	if len(points) == 0 {
		t.Fatal("expected frame points")
	}
	keys := KeyFrames(points, 10)
	if len(keys) < 2 {
		t.Fatal("expected key frames")
	}
	if EstimateFrames(asset, spec.Crop) != 200 {
		t.Fatalf("unexpected estimate: %d", EstimateFrames(asset, spec.Crop))
	}
}
