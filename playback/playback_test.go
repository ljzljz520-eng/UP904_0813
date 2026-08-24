package playback

import (
	"testing"

	"slowpreview/domain"
)

func TestPlaybackSessionAndMarkers(t *testing.T) {
	asset := domain.VideoAsset{ID: "asset", DurationMS: 6000, FrameRate: 60, Width: 1920, Height: 1080}
	spec := domain.PreviewSpec{ID: "preview", AssetID: asset.ID, Speed: domain.SpeedHalf, Crop: domain.CropWindow{StartMS: 0, EndMS: 4000}}
	session := NewSession(asset, spec)
	if session.Play() != "playing" || session.Seek(1000) != 1000 {
		t.Fatal("playback control failed")
	}
	board := NewMarkerBoard()
	if !board.Add(Marker{ID: "impact", Position: 1000, Label: "impact"}, spec) || len(board.At(1000, 0)) != 1 {
		t.Fatal("marker board failed")
	}
}
