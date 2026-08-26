package catalog

import (
	"testing"

	"slowpreview/domain"
)

func TestAssetAndPreviewQueries(t *testing.T) {
	assets := []domain.VideoAsset{{ID: "a", Title: "Jump Drill", Width: 1920, Height: 1080}, {ID: "b", Title: "Serve Drill", Width: 1280, Height: 720}}
	index := NewAssetIndex(assets)
	if len(index.Find("jump")) != 1 || len(index.FilterByResolution(domain.ResolutionLow)) != 1 {
		t.Fatal("asset query failed")
	}
	records := []domain.PreviewRecord{{Spec: domain.PreviewSpec{ID: "p1", AssetID: "a", Speed: domain.SpeedHalf, RequestedLabel: "0.5x"}, Status: domain.StatusDraft}}
	if !MatchPreview(records[0], PreviewQuery{Statuses: []string{domain.StatusDraft}}, 70) {
		t.Fatal("preview match failed")
	}
}
