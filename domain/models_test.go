package domain

import "testing"

func TestDomainModels(t *testing.T) {
	if !SpeedHalf.Valid() || SpeedHalf.Label() != "0.5x" {
		t.Fatal("half speed should be valid")
	}
	window := CropWindow{StartMS: 100, EndMS: 2100}
	if window.DurationMS() != 2000 || !window.Contains(1000) {
		t.Fatal("crop calculations failed")
	}
	asset := VideoAsset{ID: "a", SourcePath: "a.mp4", DurationMS: 1000, FrameRate: 30, Width: 1920, Height: 1080}
	if !asset.Valid() || asset.AspectRatio() < 1.7 {
		t.Fatal("asset validation failed")
	}
}
