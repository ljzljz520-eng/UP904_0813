package validate

import (
	"testing"

	"slowpreview/domain"
)

func TestValidationAndDefaults(t *testing.T) {
	asset := domain.VideoAsset{ID: "asset", SourcePath: "clip.mp4", Title: "clip", DurationMS: 10000, FrameRate: 60, Width: 1920, Height: 1080}
	if NormalizeTitle("  jump drill ") != "Jump Drill" {
		t.Fatal("title normalization failed")
	}
	if DefaultOutputPath("preview one") != "previews/preview-one.mp4" {
		t.Fatal("output path failed")
	}
	crop := ClampCrop(asset, domain.CropWindow{StartMS: -1, EndMS: 20000})
	if crop.StartMS != 0 || crop.EndMS != asset.DurationMS {
		t.Fatal("crop clamp failed")
	}
	if speed, ok := ParseSpeed("half"); !ok || speed != domain.SpeedHalf {
		t.Fatal("speed parser failed")
	}
}
